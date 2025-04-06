package config

import "github.com/joho/godotenv"

func initEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
