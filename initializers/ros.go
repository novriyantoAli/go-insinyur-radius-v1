package initializers

import (
	"github.com/go-routeros/routeros/v3"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ROS []*routeros.Client

var ROS_CLIENT []domain.ROSClient

func ConnectROS() {

	addressess := viper.GetStringSlice(`routeros.address`)
	usernames := viper.GetStringSlice(`routeros.username`)
	passwords := viper.GetStringSlice(`routeros.password`)

	if (len(addressess) != len(usernames)) && (len(addressess) != len(passwords)) && (len(usernames) != len(passwords)) {
		logrus.Fatalln("config file required same size..")
	}

	// dial the file using foreach
	for i, value := range addressess {
		roscli, err := routeros.Dial(value, usernames[i], passwords[i])
		if err != nil {
			logrus.Fatalln("cant connect to routeros ", value)
		}

		ROS_CLIENT = append(ROS_CLIENT, domain.ROSClient{Name: value, Client: roscli})

		ROS = append(ROS, roscli)
	}

}
