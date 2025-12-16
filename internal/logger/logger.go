package logger

import "go.uber.org/zap"

var Log *zap.Logger

func InitLogger(isDev bool){

	var err error
	if isDev{
		Log,err=zap.NewDevelopment()
	}else{
		Log,err=zap.NewProduction()
	}

	if err!=nil{
		panic(err)
	}
} 

func Sugar() *zap.SugaredLogger{
	return Log.Sugar()
}