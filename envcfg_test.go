package envcfg

import (
	"fmt"
	"log"
	"testing"
)

func TestEmpty(t *testing.T) {

}

func TestKeyExist(t *testing.T) {
	const (
		serverHostKey        = "server.host"
		serverNonexistingKey = "server.nonexisting"
		nonExistingSection   = "nonexisting.section"
	)

	config, err := New(ConfigFile("_example/sample.config"))
	if err != nil {
		log.Fatalf("error initializing config: %v", err)
	}

	if !config.KeyExist(serverHostKey) {
		t.Fatalf("%q should exist", serverHostKey)
	}
	if config.KeyExist(serverNonexistingKey) {
		t.Fatalf("%q shouldn't exists", serverNonexistingKey)
	}
	if config.KeyExist(nonExistingSection) {
		t.Fatalf("%q shouldn't exists", nonExistingSection)
	}
}
