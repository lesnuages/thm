package vmware

import (
	"github.com/spf13/viper"
)

var defaultConfig = `GOVC_USERNAME=""
GOVC_PASSWORD=""
GOVC_URL=""
GOVC_INSECURE=
GOVC_DATACENTER=""
GOVC_NETWORK=""
GOVC_DATASTORE=""
`

type Config struct {
	GovcURL        string `mapstructure:"GOVC_URL"`
	GovcUsername   string `mapstructure:"GOVC_USERNAME"`
	GovcPassword   string `mapstructure:"GOVC_PASSWORD"`
	GovcInsecure   bool   `mapstructure:"GOVC_INSECURE"`
	GovcDataCenter string `mapstructure:"GOVC_DATACENTER"`
	GovcDataStore  string `mapstructure:"GOVC_DATASTORE"`
	GovcNetwork    string `mapstructure:"GOVC_NETWORK"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	v.SetConfigType("env")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = v.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func GetDefaultConfig() string {
	return defaultConfig
}
