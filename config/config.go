package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type UserConfig struct {
	UseProxy            bool     `yaml:"use_proxy"`
	RetryTimes          int      `yaml:"retry_times"`
	MaxGwei             int      `yaml:"max_gwei"`
	NeedDelayAct        bool     `yaml:"need_delay_act"`
	DelayActMin         int      `yaml:"delay_act_min"`
	DelayActMax         int      `yaml:"delay_act_max"`
	NeedDelayAcc        bool     `yaml:"need_delay_acc"`
	DelayAccMin         int      `yaml:"delay_acc_min"`
	DelayAccMax         int      `yaml:"delay_acc_max"`
	TelegramAlerts      bool     `yaml:"telegram_alerts"`
	BotToken            string   `yaml:"bot_token"`
	ChatID              int      `yaml:"chat_id"`
	NeedNonEth          bool     `yaml:"need_non_eth"`
	SideChain           string   `yaml:"side_chain"`
	OkxValueMin         float64  `yaml:"okx_value_min"`
	OkxValueMax         float64  `yaml:"okx_value_max"`
	OxkAPIKey           string   `yaml:"oxk_apiKey"`
	OxkSecret           string   `yaml:"oxk_secret"`
	OxkPassword         string   `yaml:"oxk_password"`
	RelayPercentMin     int      `yaml:"relay_percent_min"`
	RelayPercentMax     int      `yaml:"relay_percent_max"`
	BungeeTimes         int      `yaml:"bungee_times"`
	BungeeValueMin      float64  `yaml:"bungee_value_min"`
	BungeeValueMax      float64  `yaml:"bungee_value_max"`
	RefuelTo            []string `yaml:"refuel_to"`
	SelfTranTimes       int      `yaml:"self_tran_times"`
	SelfTransPercentMin int      `yaml:"self_trans_percent_min"`
	SelfTransPercentMax int      `yaml:"self_trans_percent_max"`
}

func ReadSettings(filepath string) UserConfig {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error reading config file: ", err)
		return UserConfig{}
	}

	var config UserConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal("Error decoding YAML: ", err)
		return UserConfig{}
	}

	return config
}
