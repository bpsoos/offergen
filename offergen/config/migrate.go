package config

import (
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type MigrateConfig struct {
	MigrationsDirPath string
	PostgresURL       string
}

func NewMigrateConfig() *MigrateConfig {
	var k = koanf.New(".")
	err := k.Load(env.Provider("", "__", func(s string) string { return strings.ToLower(s) }), nil)
	if err != nil {
		panic(err)
	}

	return &MigrateConfig{
		PostgresURL: k.MustString("postgres_url"),
	}
}
