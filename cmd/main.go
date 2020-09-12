/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	v1 "k8s.io/api/admission/v1"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	// TODO: try this library to see if it generates correct json patch
	// https://github.com/mattbaird/jsonpatch
)

func main() {
	// webhooks must have server certs
	CertFile := "/cert/webhook.crt"
	KeyFile := "/cert/webhook-key.crt"
	sCert, err := tls.LoadX509KeyPair(CertFile, KeyFile)
	if err != nil {
		klog.Fatal(err)
	}
	tlsconfig := &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}


        mux := http.NewServeMux()
        mux.HandleFunc("/mutate", whsvr.serve)
        whsvr.server.Handler = mux

        // start webhook server in new rountine
        go func() {
                if err := whsvr.server.ListenAndServeTLS("", ""); err != nil {
                        glog.Errorf("Failed to listen and serve webhook server: %v", err)
                }
        }()

        // listening OS shutdown singal
        signalChan := make(chan os.Signal, 1)
        signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
        <-signalChan

        glog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...")
        whsvr.server.Shutdown(context.Background())

}
