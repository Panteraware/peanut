package peanut

import (
	"github.com/google/go-github/v67/github"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	URLs     []string
	Interval int64
	Port     int
	Repos    []Repo
}

type Repo struct {
	Owner string
	Name  string
	Token string
	Url   string
	Cache Cache
}

type Cache struct {
	LastUpdate       time.Time
	IsOutdated       bool
	RefreshCache     bool
	LoadCache        bool
	Latest           github.RepositoryRelease
	CacheReleaseList []github.RepositoryRelease
}

var Config config

func ConfigInit() {
	Config = config{
		Interval: 0,
		URLs:     strings.Split(os.Getenv("URLs"), ","),
		Port:     getEnvAsInt("PORT", 3000),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
