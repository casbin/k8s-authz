package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

const (
	port = "8080"
)

var (
	tlscert, tlskey string
)

func main() {
	flag.StringVar(&tlscert, "tlsCertFile", "certs/ca.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "certs/ca.key", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		glog.Errorf("Filed to load key pair: %v", err)
	}
	server := &http.Server{
		Addr:      fmt.Sprintf(":%v", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{certs}},
	}
	cs := CasbinServerHandler{}
	router := mux.NewRouter()
	router.HandleFunc("/validate", cs.serve)

	if err := server.ListenAndServeTLS("", ""); err != nil {
		glog.Error("Server error", err)
	}
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Failed to listen and serve webhook server: %v", err)
		}
	}()

	glog.Info("Server running listening in port: ", port)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	glog.Info("Shutting down webhook server...")
	if err := server.Shutdown(context.Background()); err != nil {
		glog.Error("Unable to shutdown the server", err)
	}

}
