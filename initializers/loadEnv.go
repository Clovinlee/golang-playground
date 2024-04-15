package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(f ...string) {
	var fName string

	if len(f) == 0 {
		fName = ".env"
	} else {
		fName = f[0]
	}

	err := godotenv.Load(fName)
	if err != nil {
		log.Fatal("Error loading env file: ", err)
	}
}
