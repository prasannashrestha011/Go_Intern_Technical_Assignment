package utils

import (
	"context"
	"crypto/rand"
	"main/internal/database"
	"main/internal/logger"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type VerificationCode struct {
	Code  int
	ExpiresAt time.Time
}



func StoreVerificationCode(ctx context.Context,email string,code int )error{

	duration := 3 * time.Minute
	email=strings.ToLower(strings.TrimSpace(email))
	key:="verify:"+email
	err := database.RDB.Set(ctx, key, code, duration).Err()
	if err != nil {
		return err
	}
	logger.Log.Info("Verification code",zap.Int("Code",code))
	return nil

}
func GenerateVerificationCode()int {
max := big.NewInt(90000000)          // 0–89,999,999
	n, _ := rand.Int(rand.Reader, max)   // random number in [0, 89,999,999]
	n.Add(n, big.NewInt(10000000))       // shift to 10,000,000–99,999,999
	return int(n.Int64())
}

func VerifyCode(ctx context.Context,email string,code int)bool{
	
	key:="verify:"+email
	value,err:=database.RDB.Get(ctx,key).Result()
	if err!=nil{
		if err == redis.Nil {
            logger.Log.Info("Verification code not found or expired", zap.String("email", email))
            return false
        }
		logger.Log.Error("Redis retrieval error",zap.Error(err))
		return false
	}
	storedCode,err:=strconv.Atoi(value)
    if err != nil {
        logger.Log.Error("Failed to convert Redis value to int", zap.String("value", value), zap.Error(err))
        return false
    }

	return storedCode==code
}


func DeleteVerificationCode(ctx context.Context,email string)error{
	key:="verify:"+email
	err:=database.RDB.Del(ctx,key).Err()
	if err!=nil && err!=redis.Nil{
		return err
	} 
	return nil
}

