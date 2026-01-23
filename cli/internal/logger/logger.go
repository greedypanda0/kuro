package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init(debug bool) {
	var err error

	if debug {
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}
}
