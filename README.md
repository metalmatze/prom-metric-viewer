# prom-metric-viewer

A simple CLI tool to view prometheus metrics of exporters with style.
It also includes a fancy Web UI to be even quicker viewing, sorting and searching metrics.

![screenrecording.gif](screenrecording.gif)

## Install

Dependencies: Go (with `$GOPATH` setup), Node.js (best with yarn).

```bash
go get -u -v github.com/metalmatze/prom-metric-viewer
cd $GOPATH/src/github.com/metalmatze/prom-metric-viewer
make deps ui packr install
```

You are now able to run prom-metric-viewer from anywhere on your system.

## Usage

Examples:
```bash
# Gets the metrics and prints them on your terminal
prom-metric-viewer https://localhost:9100/metrics 
# Gets the metrics and starts a web UI on localhost:8080
prom-metric-viewer --web https://localhost:9100/metrics 

# Reads the metrics from a file and prints them on your terminal
prom-metric-viewer -f examples/node-exporter 
# Reads the metrics from a file and starts a web UI on localhost:8080
prom-metric-viewer --web -f examples/node-exporter 

# Sort by name, type or help
prom-metric-viewer --sort name -f examples/node-exporter 
```
