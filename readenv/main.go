package readenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ReadEnvKeys() (key string, email string) {
	err := godotenv.Load("./readenv/.env")

	if err != nil {
		log.Fatalln("Error loading .env file")
		return
	}

	return os.Getenv("KEY"), os.Getenv("EMAIL")
}
