package config

import (
	"encoding/json"
	"github.com/team-zf/framework/utils"
)

func LoadConfig(filePath string) (conf *AppConfig, err error) {
	var bytes []byte
	bytes, err = utils.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &conf)
	return
}
