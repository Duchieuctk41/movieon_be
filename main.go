package main

import (
	"context"
	"github.com/praslar/cloud0/logger"
	"github.com/sirupsen/logrus"
	"movieon_be/conf"
	"movieon_be/pkg/route"
	"movieon_be/pkg/utils"
	"os"
)

const (
	APPNAME = "movieon_be backend golang"
)

func main() {
	conf.SetEnv()
	logger.Init(APPNAME)
	utils.LoadMessageError()
	// Dev
	logger.DefaultLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		PadLevelText:     true,
		ForceQuote:       true,
		QuoteEmptyFields: true,
	})

	//if err := utils.InitAwsSession(); err != nil {
	//	logger.Tag("main").Error(err)
	//}
	app := route.NewService()
	ctx := context.Background()
	err := app.Start(ctx)
	if err != nil {
		logger.Tag("main").Error(err)
	}
	os.Clearenv()
}
