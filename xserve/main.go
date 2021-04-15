package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/markoczy/xtools/common/logger"
)

var (
	host    string
	port    int
	folder  string
	server  http.Handler
	mode    serverMode
	https   bool
	certKey string
	cert    string
	log     logger.Logger
)

type serverMode int

const (
	modeFile = serverMode(iota)
	modeApiTrace
)

func initFlags() {
	logFactory := logger.NewAutoFlagFactory()

	hostPtr := flag.String("host", "localhost", "the designated host")
	portPtr := flag.String("port", "<default>", "designated port (must be int or '<default>')")
	folderPtr := flag.String("folder", ".", "the path to serve")
	modePtr := flag.String("mode", "fileserver", "The server mode ('fileserver', 'apitrace'")
	httpsPtr := flag.Bool("tls", false, "Serve as HTTPS (i.e. TLS)")
	certPtr := flag.String("cert", ":exec/server.crt", "Path to TLS Certificate (use ':exec' to point to the executable path)")
	certKeyPtr := flag.String("cert-key", ":exec/server.key", "Path to TLS Certificate key (use ':exec' to point to the executable path)")

	logFactory.InitFlags()
	flag.Parse()

	log = logFactory.Create()

	switch strings.ToLower(*modePtr) {
	case "fileserver":
		mode = modeFile
	case "apitrace":
		mode = modeApiTrace
	default:
		flag.Usage()
		panic("Invalid Mode selected: " + *modePtr)
	}
	if *portPtr != "<default>" {
		var err error
		port, err = strconv.Atoi(*portPtr)
		if err != nil {
			flag.Usage()
			panic("Port could not be parsed, please provide int value or '<default>''")
		}
	} else {
		port = 80
		if *httpsPtr {
			port = 443
		}
	}
	host, folder, https, cert, certKey = *hostPtr, *folderPtr, *httpsPtr, *certPtr, *certKeyPtr
	cert = replacePath(cert)
	certKey = replacePath(certKey)
}

func main() {
	initFlags()

	// Handle HTTP: 1. Log, 2. Serve file
	if mode == modeFile {
		log.Info("Serving folder \"%s\" on \"%s:%d\"", folder, host, port)
		server = http.FileServer(http.Dir(folder))
	} else {
		server = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log.Info("Running API Debugger on \"%s:%d\"\n", host, port)
		})
	}
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("********************* BEGIN REQUEST *************************")
		log.Info("*** Request: %s %s from %s", r.Method, r.URL, r.RemoteAddr)
		log.Info("*** Headers:")
		for k, v := range r.Header {
			log.Info("***   %s : %s", k, strings.Join(v, " "))
		}
		log.Info("***   Referer : %s", r.Referer())

		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Body could not be read:", err.Error())
		}
		log.Info("*** Body:")
		body := "<empty>"
		if len(d) > 0 {
			body = string(d)
		}
		log.Info("***   %s", body)
		enableCors(&w)
		server.ServeHTTP(w, r)
		log.Info("********************* END REQUEST ***************************")
	}))

	var err error
	if https {
		err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", host, port), cert, certKey, nil)
	} else {
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)

	}
	if err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func replacePath(s string) string {
	var dirAbsPath string
	ex, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("Could not replace path: %w", err))
	}
	dirAbsPath = filepath.Dir(ex)
	return strings.ReplaceAll(s, ":exec", dirAbsPath)
}
