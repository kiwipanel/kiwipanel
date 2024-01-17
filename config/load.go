package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func NewConfigServer() *Config {
	once.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		host := os.Getenv("SERVER_HOST")
		port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))

		if err != nil {
			panic("wrong server port (check your .env)")
		}
		readTimeout, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
		if err != nil {
			panic("wrong server read timeout (check your .env)")
		}

		writeTimeout, err := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
		if err != nil {
			panic("wrong server write timeout (check your .env)")
		}
		idleTimeout, err := strconv.Atoi(os.Getenv("SERVER_IDLE_TIMEOUT"))
		if err != nil {
			panic("wrong server idle timeout (check your .env)")
		}
		instance = &Config{
			Server: &serverConfig{
				Addr:         fmt.Sprintf("%s:%d", host, port),
				ReadTimeout:  time.Duration(readTimeout) * time.Second,
				WriteTimeout: time.Duration(writeTimeout) * time.Second,
				IdleTimeout:  time.Duration(idleTimeout) * time.Second,
			},
		}
	})

	return instance
}
