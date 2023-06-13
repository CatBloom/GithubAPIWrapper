package utils

import (
	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load("/go/src/work/.env")
	if err != nil {
		panic(err)
	}
}
