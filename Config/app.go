package Config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"sync"
)

var (
	AppConfig       TomlConfig
	initConfigMutex sync.Mutex
)

func InitConfig() {
	initConfigMutex.Lock()
	defer initConfigMutex.Unlock()
	content, readErr := ioutil.ReadFile("./config.toml")
	if readErr != nil {
		log.WithField("err", readErr).WithField("op", "read config.toml").Fatal()
	}
	if _, decodeConfigErr := toml.Decode(string(content), &AppConfig); decodeConfigErr != nil {
		log.WithField("err", decodeConfigErr).WithField("op", "decode config.toml").Fatal()
	}
}
