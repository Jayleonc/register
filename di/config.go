package di

import (
	"github.com/Jayleonc/register/config_center"
)

func InitConfigClient() *config_center.Client {
	configCenter, err := config_center.NewClient()
	if err != nil {
		panic(err)
	}
	return configCenter
}
