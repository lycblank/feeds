package conf

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	Server ServerConfig   `yaml:"server"`
	Swagger SwaggerConfig `yaml:"swagger"`
	Mysql      MysqlConfig `yaml:"mysql"`
	Feed FeedConfig `yaml:"feed"`
}

type MysqlConfig struct {
	UserName        string `yaml:"user_name"`
	Password        string `yaml:"password"`
	Addr            string `yaml:"addr"`
	Dbname          string `yaml:"dbname"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type ServerConfig struct {
	Host        string   `yaml:"host"`
	Port        string      `yaml:"port"`
}

type FeedConfig struct {
	FetchInterval int64 `yaml:"fetch_interval"`
}

type SwaggerConfig struct {
	Open bool          `yaml:"open"`
	DefaultAddr string `yaml:"default_addr"`
	Host string `yaml:"host"`
	BasePath string `yaml:"base_path"`
}

var config *Config
var configOnce sync.Once
func GetConfig() *Config {
	configOnce.Do(func(){
		if err := godotenv.Load(); err != nil {
			fmt.Println("not found .env file")
		}
		viper.AutomaticEnv()
		confPath := viper.GetString("CONF_PATH")
		if confPath == "" {
			confPath = "configs"
		}
		viper.SetConfigName("service")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(confPath)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		cfg := &Config{}
		if err := viper.Unmarshal(cfg, func(v *mapstructure.DecoderConfig){
			v.TagName = "yaml"
		}); err != nil {
			panic(err)
		}
		config = cfg
	})
	return config
}


