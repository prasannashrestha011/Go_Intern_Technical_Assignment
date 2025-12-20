package utils

import (
	"main/internal/logger"
	"main/internal/models"
	"sync"
	"time"

	"go.uber.org/zap"
)

var OrderQueue = make(chan models.Order,100)

var wg sync.WaitGroup
func GenerateOrderInvoice(order models.Order){
	wg.Add(1)
	defer wg.Done()

	logger.Log.Info("Generating Invoice",zap.String("orderID: ",order.ID.String()))
	time.Sleep(2*time.Second)
	logger.Log.Info("Generating Invoice completed",zap.String("orderID: ",order.ID.String()),zap.Time("time",time.Now()))
	
}

func InitOrderWorker(){
	go func ()  {
		for order:= range OrderQueue{
			go GenerateOrderInvoice(order)
		}
	}()
}

