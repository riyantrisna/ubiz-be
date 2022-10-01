package configs

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Config is a struct that will receive configuration options via environment variables.
type Config struct {
	Mode     string `mapstructure:"MODE"`
	Address  string `mapstructure:"ADDRESS"`
	Database struct {
		Mysql struct {
			DbHost          string `mapstructure:"DBHOST"`
			DbPort          int    `mapstructure:"DBPORT"`
			DbUser          string `mapstructure:"DBUSER"`
			DbPass          string `mapstructure:"DBPASS"`
			MaxIdleConns    int    `mapstructure:"MAXIDLECONNS"`
			MaxOpenConns    int    `mapstructure:"MAXOPENCONNS"`
			ConnMaxLifeTime int    `mapstructure:"CONNMAXLIFETIME"`
			ConnMaxIdleTime int    `mapstructure:"CONNMAXIDLETIME"`
		} `mapstructure:"MYSQL"`
	} `mapstructure:"DATABASE"`
	JWT struct {
		Key            string        `mapstructure:"KEY"`
		Expired        time.Duration `mapstructure:"EXPIRED"`
		ExpiredRefresh time.Duration `mapstructure:"EXPIRED_REFRESH"`
	} `mapstructure:"JWT"`
	DefaultPassword string `mapstructure:"DEFAULT_PASSWORD"`
	DefaultLang     string `mapstructure:"DEFAULT_LANG"`
	Files           struct {
		Photo string `mapstructure:"PHOTO"`
	} `mapstructure:"FILES"`
}

var (
	conf Config
	once sync.Once
)

// Get are responsible to load env and get data an return the struct
func Get() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalln("Failed reading config file.", err)
	}

	once.Do(func() {
		log.Println("Service configuration initialized.")
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatalln("Failed unmarshalling config.", err)
		}
	})

	return &conf
}
