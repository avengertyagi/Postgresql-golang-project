package models

import "github.com/akshit_tyagi/postgresql_project/internal/config"

var registry []interface{}

func Register(model interface{}) {
	registry = append(registry, model)
}

func Registered() []interface{} {
	return registry
}

func AutoMigrate() error {
	return config.DB.AutoMigrate(registry...)
}
