package utils

import (
	"main/internal/logger"

	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

var Mailer *resend.Client
func InitEmailClient(apiKey string) {

	Mailer= resend.NewClient(apiKey)
	if Mailer==nil{
		logger.Log.Error("Failed to initialized the mailer system")
	}
	logger.Log.Info("Resend mailer client initialized successfully")
}

func SendEmail(params *resend.SendEmailRequest)(success bool,message string,err error){
	sent,err:=Mailer.Emails.Send(params)
	if err!=nil{

		return false,"Failed to send email request",err
	}
	logger.Log.Info("SUCCESS: Mailer operation successful",zap.String("RequestID",sent.Id))
	return true,"",nil
}