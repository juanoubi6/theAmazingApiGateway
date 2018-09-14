package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	ENV        string
	PORT       string
	JWT_SECRET string
	CORS       string

	DB_TYPE     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string

	USER_ADMINISTRATION_SERVICE string
	POST_MANAGEMENT_SERVICE     string
	NOTIFICATIONS_SERVICE       string
}

var instance *Config

func GetConfig() *Config {
	if instance == nil {
		err := readEnv()
		if err != nil {
			panic(err)
		}
		config := newConfig()
		instance = &config
	}
	return instance
}

func newConfig() Config {
	return Config{
		ENV:        GetEnv("ENV", "develop"),
		PORT:       GetEnv("PORT", "8999"),
		JWT_SECRET: GetEnv("JWT_SECRET", "j8Ah4kO3"),
		CORS:       GetEnv("CORS", ""),

		USER_ADMINISTRATION_SERVICE: GetEnv("USER_ADMINISTRATION_SERVICE", "http://localhost:5000"),
		POST_MANAGEMENT_SERVICE:     GetEnv("POST_MANAGEMENT_SERVICE", "http://localhost:5003"),
		NOTIFICATIONS_SERVICE:       GetEnv("NOTIFICATIONS_SERVICE", "http://localhost:5004"),

		DB_TYPE:     GetEnv("DB_TYPE", "mysql"),
		DB_USERNAME: GetEnv("DB_USERNAME", "root"),
		DB_PASSWORD: GetEnv("DB_PASSWORD", "root"),
		DB_HOST:     GetEnv("DB_HOST", "127.0.0.1"),
		DB_PORT:     GetEnv("DB_PORT", "3306"),
		DB_NAME:     GetEnv("DB_NAME", "amazing-code-database"),
	}
}

func GetEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}

func readEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), "=")
		if len(values) == 2 {
			err = os.Setenv(values[0], values[1])
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
