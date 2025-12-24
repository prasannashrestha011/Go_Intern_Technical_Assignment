package utils

import (
	"context"
	"encoding/json"
	"main/internal/database"
	"main/internal/logger"

	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

type EmailJob struct {
	To   string
	From string
	Subj string
	Body string
}

func EnqueueEmail(ctx context.Context,job *EmailJob) error{
	data,err:=json.Marshal(job)
	if err!=nil{
		return err
	}
	return database.RDB.LPush(ctx,"email_queue",data).Err()
}

func EmailWorker(ctx context.Context){
	for{
		result,err:=database.RDB.BRPop(ctx,0,"email_queue").Result()
		if err!=nil{
			logger.Log.Error("Failed to pop email job",zap.Error(err))
			return
		}
	logger.Log.Info("Email worker triggered")
		var job EmailJob
		err=json.Unmarshal([]byte(result[1]),&job)
		if err!=nil{
			logger.Log.Error("Invalid job data",zap.Error(err))
			return
		}

	email_req_details:=&resend.SendEmailRequest{
		From: job.From,
		To: []string{job.To},
		Subject:job.Subj,
		Html:job.Body,
	}
	_,_,err=SendEmail(email_req_details)
	if err!=nil{
     logger.Log.Error("Email send failed", zap.Error(err))
		}
	}
}