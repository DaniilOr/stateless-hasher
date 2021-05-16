package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string) error {
	size := 16
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return err
	}
	key := base64.RawStdEncoding.EncodeToString(buf)

	mux := http.NewServeMux()
	mux.HandleFunc("/hash", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("start handle %s by: %s", request.URL.Path, key)
		h := sha256.New()
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		h.Write(b)
		_, err = writer.Write([]byte(hex.EncodeToString(h.Sum(nil))))
		if err != nil {
			log.Print(err)
		}
		log.Printf("finish handle %s by: %s", request.URL.Path, key)
	})
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("start handle %s by: %s", request.URL.Path, key)
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Printf("finish handle %s by: %s", request.URL.Path, key)
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("server %s started at %s", key, addr)
	return server.ListenAndServe()
}

