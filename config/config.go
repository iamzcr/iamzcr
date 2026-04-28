package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	AdminPort    string
	FrontendPort string
	ProjectRoot  string
}

func Load() *Config {
	godotenv.Load()

	root, _ := findProjectRoot()

	Cfg = &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", "root"),
		DBName:       getEnv("DB_NAME", "stacklifes"),
		AdminPort:    getEnv("ADMIN_PORT", "8081"),
		FrontendPort: getEnv("FRONTEND_PORT", "8082"),
		ProjectRoot:  root,
	}
	return Cfg
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return dir, fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func (c *Config) AssetDir() string {
	return filepath.Join(c.ProjectRoot, "asset")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var Cfg *Config

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
