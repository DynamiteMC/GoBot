package config

import (
	"encoding/json"
	"os"
)

type ConfigS struct {
	ID                 int64  `json:"id"`
	Token              string `json:"token"`
	Prefix             string `json:"prefix"`
	InfoPrefix         string `json:"info_prefix"`
	MuteRole           int64  `json:"mute_role"`
	DisgraceRole       int64  `json:"disgrace_role"`
	OpenAIKey          string `json:"open_ai_key"`
	DynamiteLogChannel int64  `json:"dynamite_log_channel"`
	DynamitePassword   string `json:"dynamite_password"`
	DynamiteServer     string `json:"dynamite_server"`
}

var Config ConfigS

func (config *ConfigS) Load() bool {
	d, err := os.ReadFile("config.json")
	if err != nil {
		d, _ := json.Marshal(config)
		os.WriteFile("config.json", d, 0755)
		return false
	}
	json.Unmarshal(d, config)
	return true
}

func (config *ConfigS) Save() {
	d, _ := json.Marshal(config)
	os.WriteFile("config.json", d, 0755)
}
