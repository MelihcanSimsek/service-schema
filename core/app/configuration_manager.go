package app

import "Service-schema/core/postgresql"

type ConfigurationManager struct {
	PostgresqlConfig postgresql.Config
}

func NewConfigurationManager() *ConfigurationManager {
	postgreSqlConfig := getPostgreSqlConfig()
	return &ConfigurationManager{
		PostgresqlConfig: postgreSqlConfig,
	}
}

func getPostgreSqlConfig() postgresql.Config {
	return postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		UserName:              "postgres",
		Password:              "password",
		DbName:                "product_service",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	}
}
