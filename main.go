package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/urfave/cli"
)

// RawMetric is the actual metrics exported for prometheus
type RawMetric struct {
	Element string  `json:"element"`
	Value   float64 `json:"value"`
}

// Metric has a name, type, help and cardinality.
type Metric struct {
	Name       string
	Type       string
	Help       string
	RawMetrics []RawMetric
}

// Cardinality is the number of different labels in one metric.
func (m Metric) Cardinality() int {
	return len(m.RawMetrics)
}

func main() {
	app := cli.NewApp()

	app.Action = ViewAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Read metrics from file",
		},
		cli.StringFlag{
			Name:  "sort",
			Usage: "Sort by name, type or help",
		},
		cli.BoolFlag{
			Name:  "web",
			Usage: "Show the overview via a website",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// ViewAction runs the default cli app.
func ViewAction(c *cli.Context) error {
	url := c.Args().First()
	sortFlag := c.String("sort")
	fileFlag := c.String("file")
	webFlag := c.Bool("web")

	var metrics []Metric
	if fileFlag != "" {
		fileMetrics, err := FileMetrics(fileFlag)
		if err != nil {
			return err
		}
		metrics = fileMetrics
	} else {
		if url == "" {
			return errors.New("Please provide a URL, like http://localhost:8080/metrics")
		}
		httpMetrics, err := HTTPMetrics(url)
		if err != nil {
			return err
		}
		metrics = httpMetrics
	}

	switch sortFlag {
	case "type":
		sort.Slice(metrics, func(i, j int) bool {
			return metrics[i].Type < metrics[j].Type
		})
	case "cardinality":
		sort.Slice(metrics, func(i, j int) bool {
			return metrics[i].Cardinality() < metrics[j].Cardinality()
		})
	case "help":
		sort.Slice(metrics, func(i, j int) bool {
			return metrics[i].Help < metrics[j].Help
		})
	default:
		sort.Slice(metrics, func(i, j int) bool {
			return metrics[i].Name < metrics[j].Name
		})
	}

	if webFlag {
		return printWeb(metrics)
	}

	return printCli(metrics)
}

type jsonMetric struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Cardinality int    `json:"cardinality"`
	Help        string `json:"help"`
}

func printWeb(metrics []Metric) error {
	box := packr.NewBox("./public")
	tem := template.New("index")
	tem, err := tem.Parse(box.String("index.html"))
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		tem.Execute(w, metrics)
	})

	http.HandleFunc("/metrics.json", func(w http.ResponseWriter, r *http.Request) {
		queryName := r.URL.Query().Get("name")
		if queryName != "" {
			for _, metric := range metrics {
				if metric.Name == queryName {
					payload, err := json.Marshal(metric.RawMetrics)
					if err != nil {
						log.Println(err)
						http.Error(w, "failed to marshal raw metrics", http.StatusInternalServerError)
						return
					}

					w.Write(payload)
					w.Header().Set("Content-Type", "application/json")
					return
				}
			}

			fmt.Fprintln(w, "metric name was not found")
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		var jsonMetrics []jsonMetric
		for _, metric := range metrics {
			jsonMetrics = append(jsonMetrics, jsonMetric{
				Name:        metric.Name,
				Type:        metric.Type,
				Cardinality: metric.Cardinality(),
				Help:        metric.Help,
			})
		}

		data, err := json.Marshal(jsonMetrics)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to marshal json", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data)
	})

	http.HandleFunc("/metrics.md", func(w http.ResponseWriter, r *http.Request) {
		var filtered []Metric

		contains := r.URL.Query().Get("contains")
		if contains != "" {
			for _, metric := range metrics {
				if strings.Contains(metric.Name, contains) {
					filtered = append(filtered, metric)
				}
			}
		} else {
			filtered = metrics
		}

		tab := tabwriter.NewWriter(w, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(tab, "|Name\t|Type\t|Cardinality\t|Help")
		fmt.Fprintln(tab, "|----\t|----\t|-----------\t|----")
		for _, metric := range filtered {
			if metric.Name != "" {
				fmt.Fprintf(tab, "| %s\t| %s\t| %d\t| %s\n", metric.Name, metric.Type, metric.Cardinality(), metric.Help)
			}
		}
		tab.Flush()
	})

	http.Handle("/build.js", http.FileServer(box))

	return http.ListenAndServe(":8080", nil)
}

func printCli(metrics []Metric) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Name\tType\tCardinality\tHelp")
	for _, metric := range metrics {
		if metric.Name != "" {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", metric.Name, metric.Type, metric.Cardinality(), metric.Help)
		}
	}
	return w.Flush()
}

// FileMetrics reads a file and returns its metrics.
func FileMetrics(filename string) ([]Metric, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseMetrics(file)
}

// HTTPMetrics makes an HTTP GET request again url and returns its metrics.
func HTTPMetrics(url string) ([]Metric, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMetrics(resp.Body)
}

func parseMetrics(r io.Reader) ([]Metric, error) {
	metricsMap := make(map[string]Metric, 64)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# HELP") {
			name, help := parseHelp(line)
			metricsMap[name] = Metric{
				Name: name,
				Help: help,
			}
		} else if strings.HasPrefix(line, "# TYPE") {
			name, types := parseType(line)
			metric := metricsMap[name]
			metric.Type = types
			metricsMap[name] = metric
		} else {
			name, raw := parseRawMetric(line)
			metric := metricsMap[name]
			metric.RawMetrics = append(metric.RawMetrics, raw)
			metricsMap[name] = metric
		}
	}

	var metrics []Metric
	for _, metric := range metricsMap {
		if metric.Name != "" {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func parseHelp(line string) (string, string) {
	splits := strings.SplitN(line, " ", 4)
	return splits[2], splits[3]
}

func parseType(line string) (string, string) {
	splits := strings.SplitN(line, " ", 4)
	return splits[2], splits[3]
}

func parseRawMetric(line string) (string, RawMetric) {
	spaces := strings.Split(line, " ")
	name := parseName(line)

	// value is the last element separated by a space
	rawValue := spaces[len(spaces)-1]
	// element is all elements except the last, which is the value
	rawElement := strings.Join(spaces[:len(spaces)-1], " ")

	if rawValue == "NaN" {
		return name, RawMetric{Element: rawElement, Value: 0}
	}

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		log.Println("failed to parse value from metric to float64", rawValue)
		return "", RawMetric{Element: rawElement, Value: 0}
	}

	return name, RawMetric{Element: rawElement, Value: value}
}

func parseName(line string) string {
	spaces := strings.Split(line, " ")
	return strings.Split(spaces[0], "{")[0]
}
