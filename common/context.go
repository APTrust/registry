package common

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

var ctx *APTContext

type APTContext struct {
	Config *Config
	DB     *pg.DB
}

// Context returns an APTContext object, which includes
// global config settings and a connection to the postgres database.
// It requires the environment variable APT_ENV to be set to something
// valid, such as "test", "dev", "integration", "demo", "staging" or
// "production". It loads the .env file that corresponds to that setting.
// If APT_ENV is not set to a valid setting, the app dies immediately.
//
// This will also exit if the app cannot connect to the database.
// If that happens, ensure the database is running and accepting
// connections at the specified location, and ensure that the db
// credentials are correct.
//
func Context() *APTContext {
	if ctx == nil {
		config := NewConfig()
		db := pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%d", config.DB.Host, config.DB.Port),
			User:     config.DB.User,
			Password: config.DB.Password,
			Database: config.DB.Name,
		})
		ctx = &APTContext{
			Config: config,
			DB:     db,
		}
	}
	return ctx
}
