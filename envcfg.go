// Package envcfg provides a common way to pull configuration variables from
// the environment or files on disk.
package envcfg

import (
	"bytes"
	"encoding/base64"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/knq/ini"
)

type Filter func(*Envcfg, string) string

// Envcfg
type Envcfg struct {
	envVarName string
	configFile string
	config     *ini.File

	filters map[string]Filter
}

// New creates an Envcfg.
func New(opts ...Option) (*Envcfg, error) {
	var err error

	// default values
	ec := &Envcfg{
		envVarName: "APP_CONFIG",
		configFile: "env/config",
		filters:    make(map[string]Filter),
	}

	// apply options
	for _, o := range opts {
		o(ec)
	}

	// load environment data from $ENV{$envVarName} or from file $PWD/$configFile
	if envdata := os.Getenv(ec.envVarName); envdata != "" {
		// if the data is supplied in the $ENV, then base64 decode the data
		var data []byte
		data, err = base64.StdEncoding.DecodeString(envdata)
		if err == nil {
			r := bytes.NewReader(data)
			ec.config, err = ini.Load(r)
		}
	} else {
		ec.config, err = ini.LoadFile(ec.configFile)
	}

	// ensure no err
	if err != nil {
		return nil, err
	}

	// set git style config
	ec.config.SectionNameFunc = ini.GitSectionNameFunc
	ec.config.SectionManipFunc = ini.GitSectionManipFunc
	ec.config.ValueManipFunc = func(val string) string {
		val = strings.TrimSpace(val)

		if str, err := strconv.Unquote(val); err == nil {
			val = str
		}

		return val
	}

	return ec, nil
}

// envValRegexp matches the definition of "$NAME||VAL"
var envValRegexp = regexp.MustCompile(`(?i)^\$([a-z][a-z0-9_]*)\|\|(.+)$`)

// GetKey retrieves a key from the environment, or the supplied configuration
// data.
func (ec *Envcfg) GetKey(key string) string {
	val := ec.config.GetKey(key)

	// check if its $NAME is present
	matches := envValRegexp.FindStringSubmatch(val)
	if len(matches) == 3 {
		// has $NAME, so look at ENV{$NAME}
		v := os.Getenv(matches[1])

		// if empty value, use the default value
		if v == "" {
			val = matches[2]
		} else {
			val = v
		}
	}

	// apply key value filter
	if f, ok := ec.filters[key]; ok {
		return f(ec, val)
	}

	return val
}
