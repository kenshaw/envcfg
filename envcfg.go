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

const (
	// DefaultVarName is the default variable name to read the configuration
	// from.
	DefaultVarName = "APP_CONFIG"

	// DefaultConfigFile is the default config file path to read the
	// configuration from.
	DefaultConfigFile = "env/config"
)

// Filter is a filter that modifies a key returned from the envcfg.
type Filter func(*Envcfg, string) string

// Envcfg is config loaded from the environment.
type Envcfg struct {
	config *ini.File

	envVarName string
	configFile string

	filters map[string]Filter
}

// New creates an Envcfg.
func New(opts ...Option) (*Envcfg, error) {
	var err error

	// default values
	ec := &Envcfg{
		envVarName: DefaultVarName,
		configFile: DefaultConfigFile,
		filters:    make(map[string]Filter),
	}

	// apply options
	for _, o := range opts {
		o(ec)
	}

	// load environment data from $ENV{$envVarName} or from file $configFile
	if envdata := os.Getenv(ec.envVarName); envdata != "" {
		// if the data is supplied in $ENV{$envVarName}, then base64 decode the data
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

	// check if config data is like "$NAME||default"
	matches := envValRegexp.FindStringSubmatch(val)
	if len(matches) == 3 {
		// config data has $NAME, so read $ENV{$NAME}
		v := os.Getenv(matches[1])

		// if empty value, use the default
		if v == "" {
			val = matches[2]
		} else {
			val = v
		}
	}

	// apply filter
	if f, ok := ec.filters[key]; ok {
		return f(ec, val)
	}

	return val
}

// GetString retrieves a key from the environment or the supplied configuration
// data and returns it as a string.
//
// Alias for GetKey
func (ec *Envcfg) GetString(key string) string {
	return ec.GetKey(key)
}

// GetBool retrieves a key from the environment, or the supplied configuration
// data and returns it as a bool.
func (ec *Envcfg) GetBool(key string) bool {
	b, _ := strconv.ParseBool(ec.GetKey(key))
	return b
}

// GetFloat retrieves a key from the environment, or the supplied configuration
// data and returns it as a float64. Uses bitSize as the precision.
func (ec *Envcfg) GetFloat(key string, bitSize int) float64 {
	f, _ := strconv.ParseFloat(ec.GetKey(key), bitSize)
	return f
}

// GetInt64 retrieves a key from the environment, or the supplied configuration
// data and returns it as a int64. Uses base and bitSize to parse.
func (ec *Envcfg) GetInt64(key string, base, bitSize int) int64 {
	i, _ := strconv.ParseInt(ec.GetKey(key), base, bitSize)
	return i
}

// GetUint64 retrieves a key from the environment, or the supplied configuration
// data and returns it as a uint64. Uses base and bitSize to parse.
func (ec *Envcfg) GetUint64(key string, base, bitSize int) uint64 {
	u, _ := strconv.ParseUint(ec.GetKey(key), base, bitSize)
	return u
}

// GetInt retrieves a key from the environment, or the supplied configuration
// data and returns it as a int. Expects numbers to be base 10 and no larger
// than 32 bits.
func (ec *Envcfg) GetInt(key string) int {
	i, _ := strconv.Atoi(ec.GetKey(key))
	return i
}
