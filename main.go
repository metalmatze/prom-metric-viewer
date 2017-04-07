package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/urfave/cli"
)

type Metric struct {
	Name        string
	Type        string
	Help        string
	Cardinality int
}

func main() {
	app := cli.NewApp()

	app.Action = ViewAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "sort",
			Usage: "sort by name, type or help",
		},
		cli.BoolFlag{},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func ViewAction(c *cli.Context) error {
	url := c.Args().First()
	sortFlag := c.String("sort")

	if url == "" {
		return errors.New("Please provide a URL, like http://localhost:8080/metrics")
	}

	metrics, err := getMetrics(url)
	if err != nil {
		return err
	}

	switch sortFlag {
	case "type":
		sort.Slice(metrics, func(i, j int) bool {
			return metrics[i].Type < metrics[j].Type
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Name\tType\tCardinality\tHelp")
	for _, metric := range metrics {
		if metric.Name != "" {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", metric.Name, metric.Type, metric.Cardinality, metric.Help)
		}
	}
	return w.Flush()

}

func getMetrics(url string) ([]Metric, error) {
	metricsMap := make(map[string]Metric, 64)

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

	scanner := bufio.NewScanner(resp.Body)
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
			name := parseMetric(line)
			metric := metricsMap[name]
			metric.Cardinality = metric.Cardinality + 1
			metricsMap[name] = metric
		}
	}

	metrics := make([]Metric, len(metricsMap))
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

func parseMetric(line string) string {
	spaces := strings.Split(line, " ")
	return strings.Split(spaces[0], "{")[0]
}
