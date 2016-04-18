package main

import (
	"flag"
	"fmt"
	"gotpl/lib"
	"gotpl/web/handlers"
	"net/http"
	"net/http/pprof"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	path := flag.String("config", "../config/dev.toml", "Config file path")
	flag.Parse()

	lib.LoadConfig(*path)

	// Router
	mux := mux.NewRouter()

	// Handlers
	handlers.RegisterHome(mux)

	// Profiling support
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Instrumentation
	mux.Handle("/metrics", prometheus.Handler())

	// HTTP middleware
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	if !lib.CFG.Production {
		n.Use(negroni.NewLogger())
	}

	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseHandler(mux)

	fmt.Printf("Listening on %s\n", lib.CFG.HTTPBind)
	http.ListenAndServe(lib.CFG.HTTPBind, n)
}
