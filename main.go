package main

// Based on https://github.com/kubernetes/kubernetes/blob/release-1.24/test/images/agnhost/webhook/main.go

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"k8s.io/klog/v2"
)

func (c *Config) serveMutateNamespace(w http.ResponseWriter, r *http.Request) {
	serve(w, r, newDelegateToV1AdmitHandler(c.mutateNamespace))

}

func main() {
	if len(os.Args) != 2 {
		klog.Fatal("Usage: %s <configuration-file>")
	}
	configFile := os.Args[1]
	config, err := NewConfigFromFile(configFile)
	if err != nil {
		klog.Fatal("Loading configuration: %e", err)
	}

	http.HandleFunc("/mutating-namespace", config.serveMutateNamespace)

	sCert, err := tls.LoadX509KeyPair(config.Server.TLS.CertFile, config.Server.TLS.KeyFile)
	if err != nil {
		klog.Fatal("Loading: CertFile: %s, CertKey: %s: Error: ",
			config.Server.TLS.CertFile, config.Server.TLS.KeyFile, err)
	}
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Server.Port),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{sCert},
		},
	}
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		klog.Fatal("Starting HTTP server: %e", err)
	}
}
