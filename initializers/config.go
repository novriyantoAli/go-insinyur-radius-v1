package initializers

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() {

	if os.Getenv("ENV") == "production" {
		viper.SetConfigName("config_prod")
	} else {
		viper.SetConfigName("config_dev")
	}
	// name of config file (without extension)
	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/ir/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.ir") // call multiple times to add many search paths
	viper.AddConfigPath(".")         // optionally look for config in the working directory
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if viper.GetBool(`debug`) {
		logrus.Infoln("SERVICE RUN IN DEBUG MODE")
	}
}
