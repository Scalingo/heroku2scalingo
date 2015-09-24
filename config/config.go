package config

import (
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Scalingo/envconfig"
	"github.com/Scalingo/go-scalingo"
)

type Config struct {
	ApiUrl     string
	apiHost    string
	ApiVersion string
	ConfigDir  string
	AuthFile   string
}

var (
	env = map[string]string{
		"API_URL":     "https://api.scalingo.com",
		"API_VERSION": "1",
		"CONFIG_DIR":  ".config/scalingo",
		"AUTH_FILE":   "auth",
	}
	C Config
)

func init() {
	home := HomeDir()
	if home == "" {
		panic("The HOME environment variable must be defined")
	}

	env["CONFIG_DIR"] = filepath.Join(home, env["CONFIG_DIR"])
	env["AUTH_FILE"] = filepath.Join(env["CONFIG_DIR"], env["AUTH_FILE"])

	for k := range env {
		vEnv := os.Getenv(k)
		if vEnv == "" {
			os.Setenv(k, env[k])
		}
	}

	envconfig.Process("", &C)

	u, err := url.Parse(C.ApiUrl)
	if err != nil {
		panic("API_URL is not a valid URL " + err.Error())
	}

	C.apiHost = strings.Split(u.Host, ":")[0]

	scalingo.ApiAuthenticator = Authenticator
	scalingo.ApiUrl = C.ApiUrl
	scalingo.ApiVersion = C.ApiVersion
}

func HomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
