// Copyright 2024 eve.  All rights reserved.

//nolint:errcheck
package options

import (
	"errors"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"

	"github.com/snail-plus/gopkg/log"
)

var _ IOptions = (*HealthOptions)(nil)

// HealthOptions defines options for redis cluster.
type HealthOptions struct {
	// Enable debugging by exposing profiling information.
	HTTPProfile        bool   `json:"enable-http-profiler" mapstructure:"enable-http-profiler"`
	HealthCheckPath    string `json:"check-path" mapstructure:"check-path"`
	HealthCheckAddress string `json:"check-address" mapstructure:"check-address"`
}

// NewHealthOptions create a `zero` value instance.
func NewHealthOptions() *HealthOptions {
	return &HealthOptions{
		HTTPProfile:        false,
		HealthCheckPath:    "/healthz",
		HealthCheckAddress: "0.0.0.0:20250",
	}
}

// Validate verifies flags passed to HealthOptions.
func (o *HealthOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to redis storage for a specific APIServer to the specified FlagSet.
func (o *HealthOptions) AddFlags(fs *pflag.FlagSet, prefixs ...string) {
	fs.BoolVar(&o.HTTPProfile, "health.enable-http-profiler", o.HTTPProfile, "Expose runtime profiling data via HTTP.")
	fs.StringVar(&o.HealthCheckPath, "health.check-path", o.HealthCheckPath, "Specifies liveness health check request path.")
	fs.StringVar(&o.HealthCheckAddress, "health.check-address", o.HealthCheckAddress, "Specifies liveness health check bind address.")
}

// ServeHealthCheck 启动健康检查服务。返回 error 而非调用 Fatalf，
// 避免作为库使用时直接终止宿主进程。
func (o *HealthOptions) ServeHealthCheck() error {
	r := mux.NewRouter()

	r.HandleFunc(o.HealthCheckPath, handler).Methods(http.MethodGet)
	if o.HTTPProfile {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		r.HandleFunc("/debug/pprof/{_:.*}", pprof.Index)
		log.Warnw("pprof is enabled; ensure the health check address is not publicly reachable",
			"addr", o.HealthCheckAddress)
	}

	server := &http.Server{
		Addr:              o.HealthCheckAddress,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Infow("Starting health check server", "path", o.HealthCheckPath, "addr", o.HealthCheckAddress)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("Error serving health check endpoint: %v", err)
		return err
	}
	return nil
}

func handler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status": "ok"}`))
}
