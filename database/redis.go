package database

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // "localhost:6379" gibi
		Password: "",                      // Redis şifre yoksa boş
		DB:       0,                       // default db
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis bağlantısı başarısız:", err)
	}

	log.Println("Redis'e bağlanıldı")
}
