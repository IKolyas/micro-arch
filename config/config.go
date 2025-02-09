package config

import (
	"fmt"
	"os"
	"strconv"
)

type (
	Config struct {
		AppConfig struct {
			AppPort   string
			JwtSecret string
		}
		DBConfig struct {
			Host     string
			Port     string
			User     string
			Password string
			Dbname   string
			Replicas int
		}
	}
)

var (
	Cnf = Config{}
)

func (c *Config) Load() error {

	q, err := strconv.Atoi(os.Getenv("PGSQL_SLAVES"))
	if err != nil {
		return err
	}

	c.AppConfig.AppPort = os.Getenv("APP_PORT")
	c.AppConfig.JwtSecret = os.Getenv("JWT_SECRET")

	c.DBConfig.Host = os.Getenv("PGSQL_HOST")
	c.DBConfig.Port = os.Getenv("PGSQL_PORT")
	c.DBConfig.User = os.Getenv("PGSQL_USER")
	c.DBConfig.Password = os.Getenv("PGSQL_PASSWORD")
	c.DBConfig.Dbname = os.Getenv("PGSQL_DB")
	c.DBConfig.Replicas = q

	if c.AppConfig.AppPort == "" {
		return fmt.Errorf("APP_PORT is not set")
	}
	if c.AppConfig.JwtSecret == "" {
		return fmt.Errorf("JWT_SECRET is not set")
	}
	if c.DBConfig.Host == "" {
		return fmt.Errorf("PGSQL_HOST is not set")
	}
	if c.DBConfig.Port == "" {
		return fmt.Errorf("PGSQL_PORT is not set")
	}
	if c.DBConfig.User == "" {
		return fmt.Errorf("PGSQL_USER is not set")
	}
	if c.DBConfig.Password == "" {
		return fmt.Errorf("PGSQL_PASSWORD is not set")
	}
	if c.DBConfig.Dbname == "" {
		return fmt.Errorf("PGSQL_DB is not set")
	}

	fmt.Printf("Config loaded successfully: %+v\n", c)
	return nil
}
