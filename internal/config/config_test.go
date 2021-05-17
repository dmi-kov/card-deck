package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	cfgExpected := &Config{
		Database: struct {
			Database       string `hcl:"database"`
			Host           string `hcl:"host"`
			Port           int    `hcl:"port"`
			User           string `hcl:"user"`
			Password       string `hcl:"password"`
			MigrationsPath string `hcl:"migrations"`
			SSLMode        string `hcl:"sslmode"`
		}{
			Database:       "test",
			Host:           "localhost",
			Port:           5437,
			User:           "test",
			Password:       "test",
			MigrationsPath: "./db/migrations",
			SSLMode:        "disable",
		},
		App: struct {
			Listening         int  `hcl:"listening"`
			Prod              bool `hcl:"prod"`
			DisableStacktrace bool `hcl:"disableStacktrace"`
		}{
			Listening:         8083,
			Prod:              true,
			DisableStacktrace: true,
		},
	}

	cfgActual, err := NewConfig("config.example.hcl")
	require.NoError(t, err)

	require.Equal(t, cfgExpected, cfgActual)
}

func TestDatabase_connectionURL(t *testing.T) {
	db := Database{
		Database: "database",
		Host:     "localhost",
		Port:     1111,
		User:     "user",
		Password: "password",
		SSLMode:  "disable",
	}
	expectedURL := "postgres://user:password@localhost:1111/database?sslmode=disable"

	require.Equal(t, expectedURL, db.connectionURL())
}
