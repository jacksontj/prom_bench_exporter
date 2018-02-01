package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

/*
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 3.2126e-05
*/

var opts struct {
	ListenPort int   `long:"port" default:"8888" description:"listen port"`
	Metrics    int64 `long:"metrics" default:"1" description:"how many metric names"`
	Labels     int64 `long:"labels" default:"1" description:"how many labels per metric"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	for metricCount := int64(0); metricCount < opts.Metrics; metricCount++ {
		metricName := fmt.Sprintf("test_metric_%d", metricCount)

		label := make([]string, opts.Labels)
		for labelCount := int64(0); labelCount < opts.Labels; labelCount++ {
			label[labelCount] = fmt.Sprintf(`label%d="%d"`, labelCount, labelCount)
		}
		labelSet := strings.Join(label, ", ")
		fmt.Fprintf(w, "# HELP %s test metric\n", metricName)
		fmt.Fprintf(w, "# TYPE %s gauge\n", metricName)
		fmt.Fprintf(w, "%s{%s} %d\n", metricName, labelSet, metricCount)
	}
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		logrus.Fatalf("Error parsing flags: %v", err)
	}

	http.HandleFunc("/", handler)
	fmt.Println(fmt.Sprintf(":%d", opts.ListenPort))
	http.ListenAndServe(fmt.Sprintf(":%d", opts.ListenPort), nil)
}
