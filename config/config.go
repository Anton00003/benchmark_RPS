package config

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DSN          string
	TextRequest  string
	Level        string
	TestTime     time.Duration
	ParallelProc int
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	testTime, err := getEnvDur("TEST_TIME")
	if err != nil {
		return nil, errors.WithMessage(err, "error getting TEST_TIME")
	}

	return &Config{
		DSN:          getEnv("DSN"),
		TextRequest:  getEnv("TEXT_REQUEST"),
		TestTime:     testTime,
		ParallelProc: getEnvInt("PARALLEL_PROC"),
		Level:        getEnv("LEVEL"),
	}, nil
}

func getEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return ""
}

func getEnvInt(key string) int {
	valS := getEnv(key)
	valI, err := strconv.Atoi(valS)
	if err != nil {
		return 0
	}
	return valI
}

func getEnvDur(key string) (time.Duration, error) {
	return time.ParseDuration(getEnv(key))

}
