// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	_ "net/http/pprof"
	"log"
	"os"
        "fmt"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "sampleexporter"
        subsystem = "subsystem"
)
var hostnameLabelValues []string
var hostnameEnabled = prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "enabled"), "Is this server enabled?.", []string{"hostname"}, nil)

type Exporter struct {
        hostname	string
        enabled	        bool
}

func NewExporter() (*Exporter, error) {
        var hostname string
        var enabled bool
        hostname = "muj-server"
        enabled = true

        return &Exporter{
                hostname: hostname,
		enabled: enabled,
        }, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
        fmt.Println("This is describe func")
        fmt.Println(e.hostname)
        fmt.Println(e.enabled)

        ch <- hostnameEnabled

}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
        fmt.Println("This is collect func")
        fmt.Println(e.hostname)
        fmt.Println(e.enabled)

        var hostname string
        var enabled int

        hostname = e.hostname
        if e.enabled == true {
		enabled = 1
	}

        hostnameLabelValues = []string{hostname}
        ch <- prometheus.MustNewConstMetric(hostnameEnabled, prometheus.GaugeValue, float64(enabled), hostnameLabelValues...)

}

func main() {
        var (
                app             = kingpin.New("libvirt_exporter", "Prometheus metrics exporter for libvirt")
                listenAddress   = app.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9177").String()
                metricsPath     = app.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
        )

        kingpin.MustParse(app.Parse(os.Args[1:]))

        exporter, _ := NewExporter()

        prometheus.MustRegister(exporter)
        prometheus.MustRegister(version.NewCollector("sample_exporter"))

        http.Handle(*metricsPath, promhttp.Handler())
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                w.Write([]byte(`
                        <html>
                        <head><title>Sample Exporter</title></head>
                        <body>
                        <h1>Sample Exporter</h1>
                        <p><a href='` + *metricsPath + `'>Metrics</a></p>
                        </body>
                        </html>`))
        })
        log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
