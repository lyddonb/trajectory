/*
Live swap config file class taken from:
http://openmymind.net/Golang-Hot-Configuration-Reload/
*/

package trajectory

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	HostPort  int
	RedisHost string
}

var (
	config      *Config
	configLock  = new(sync.RWMutex)
	defaultPort = 1200
	host        = "127.0.0.1"
)

func loadConfig(fail bool) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		handleError("open", err)
		return
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		handleError("parse", err)
		return
	}

	setConfig(temp)
}

func handleError(step string, err error) {
	log.Println(step, " config: ", err)
	setDefaultConfig()
}

func setDefaultConfig() {
	c := new(Config)
	c.HostPort = defaultPort
	c.RedisHost = host
	setConfig(c)
}

func setConfig(c *Config) {
	configLock.Lock()
	config = c
	configLock.Unlock()
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// go calls init on start
func init() {
	loadConfig(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(false)
			log.Println("Reloaded")
		}
	}()
}
