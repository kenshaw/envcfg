// examples/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brankas/envcfg"
)

func main() {
	// load config from the default APP_CONFIG environment variable
	config, err := envcfg.New()
	if err != nil {
		log.Fatal(err)
	}

	// load additional config from SOME_OTHER_VAR environment variable
	config2, err := envcfg.New(
		envcfg.VarName("SOME_OTHER_VAR"),
	)
	if err != nil {
		log.Fatal(err)
	}
	config2 = config2

	// read a config key
	val := config.GetKey("mysection.mykeyname")
	log.Printf("> val: %s", val)

	// create a http.Server with a host, port, and TLS based on config pulled
	// from environment
	s := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", config.Host(), config.Port()),
		TLSConfig: config.TLS(nil),
	}
	log.Fatal(s.ListenAndServe())
}
