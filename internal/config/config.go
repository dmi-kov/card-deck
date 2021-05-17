package config

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	// required for database
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Config represents root configuration object
type Config struct {
	Database Database `hcl:"db,block"`
	App      App      `hcl:"app,block"`
}

// Database represents Postgres configuration
type Database struct {
	Database       string `hcl:"database"`
	Host           string `hcl:"host"`
	Port           int    `hcl:"port"`
	User           string `hcl:"user"`
	Password       string `hcl:"password"`
	MigrationsPath string `hcl:"migrations"`
	SSLMode        string `hcl:"sslmode"`
}

// App represents general application configuration
type App struct {
	Listening         int  `hcl:"listening"`
	Prod              bool `hcl:"prod"`
	DisableStacktrace bool `hcl:"disableStacktrace"`
}

// NewConfig reads configuration from given config path
func NewConfig(configPath string) (*Config, error) {
	var config Config

	err := hclsimple.DecodeFile(configPath, nil, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// NewDBConnection establish new connection to DB
func (c *Config) NewDBConnection(ctx context.Context) (*sqlx.DB, error) {
	if err := c.Database.migrateUP(); err != nil {
		return nil, errors.Wrap(err, "failed to execute migration")
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", c.Database.connectionURL())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to master database")
	}

	return db, nil
}

func (d *Database) migrateUP() error {
	m, err := migrate.New(fmt.Sprint("file://", d.MigrationsPath), d.connectionURL())
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (d *Database) connectionURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", d.User, d.Password, d.Host, d.Port, d.Database, d.SSLMode)
}

// InitLogger init Zap logger and calls ReplaceGlobals
func (a *App) InitLogger() (*zap.Logger, error) {
	var zapConfig zap.Config
	if a.Prod {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}
	if a.DisableStacktrace {
		zapConfig.DisableStacktrace = true
	}
	log, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed init Zap logger: %v", err)
	}
	zap.ReplaceGlobals(log)

	return log, nil
}
