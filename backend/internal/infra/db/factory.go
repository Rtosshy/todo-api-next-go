package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	InstanceSQLite int = iota
	InstancePostgres
)

var (
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

func NewDatabaseSQLFactory(instance int) (db *gorm.DB, err error) {
	switch instance {
	case InstancePostgres:
		configs := NewConfigPostgres()
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", configs.User, configs.Password, configs.Host, configs.Port, configs.Database, configs.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	case InstanceSQLite:
		configs := NewConfigSQLite()
		db, err = gorm.Open(sqlite.Open(configs.Database), &gorm.Config{})
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
	return db, err
}
