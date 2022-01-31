package db

import "time"

type redisObject struct {
	Exp  int
	Data interface{}
}

var redisLite = map[string]redisObject{}

func redisSet(key string, value interface{}) {
	redisLite[key] = redisObject{}
	time.Now().Second()
}
