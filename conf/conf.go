package conf

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	Listen         string
	KubeConfigPath string
}

// Conf ...
var Conf *Config

func defaultConf() (c *Config) {
	c = &Config{
		Listen:         "0.0.0.0:9090",
		KubeConfigPath: path.Join(filepath.Dir(os.Args[0]), "kubeconfig"),
	}
	return
}

// Init ...
func Init() {
	var (
		err        error
		configFile string
	)
	Conf = defaultConf()

	flag.StringVar(&configFile, "config", "", "config file")
	flag.Parse()

	log.Println("config file:", configFile)
	if _, err = toml.DecodeFile(configFile, Conf); err != nil {
		log.Println("config DecodeFile err:", err)
		return
	}
}
