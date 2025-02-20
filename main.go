// Copyright 2024 The Prometheus Authors
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
	"flag"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("web.listen-address", ":9104", "Address to listen on for web interface and telemetry.")
	metricPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	certPath      = flag.String("tls.server-crt", "", "Path to PEM encoded file containing TLS server cert.")
	keyPath       = flag.String("tls.server-key", "", "Path to PEM encoded file containing TLS server key (unencyrpted).")
	silent        = flag.Bool("silent", false, "Disable logging of errors in handling stats lines")
)

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE|syslog.LOG_SYSLOG, "rsyslog_exporter")
	if e == nil {
		log.SetOutput(logwriter)
	}

	flag.Parse()
	exporter := newRsyslogExporter()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Print("interrupt received, exiting")
		os.Exit(0)
	}()

	go func() {
		exporter.run(*silent)
	}()

	prometheus.MustRegister(exporter)
	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// nolint:errcheck
		w.Write([]byte(`<html>
<head><title>Rsyslog exporter</title></head>
<body>
<h1>Rsyslog exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`))
	})

	if *certPath == "" && *keyPath == "" {
		log.Printf("Listening on %s", *listenAddress)
		log.Fatal(http.ListenAndServe(*listenAddress, nil))
	} else if *certPath == "" || *keyPath == "" {
		log.Fatal("Both tls.server-crt and tls.server-key must be specified")
	} else {
		log.Printf("Listening for TLS on %s", *listenAddress)
		log.Fatal(http.ListenAndServeTLS(*listenAddress, *certPath, *keyPath, nil))
	}
}
