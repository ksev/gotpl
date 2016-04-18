package lib

import (
	"database/sql"
	"log"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

// Config the 'main' configuration struct this maps directly to the config file
type Config struct {
	HTTPBind           string
	Production         bool
	DBConnectionString string
}

// CFG global Config instance
var CFG *Config

// DB refrence to the global DB pool
var DB *sql.DB

// LoadConfig load file as config file for the specific path, should be called early in the programs lifetime.
func LoadConfig(path string) {
	if _, err := toml.DecodeFile(path, &CFG); err != nil {
		// This cant be logged to file in production mode since i dont know what file to log to yet
		log.Fatal(err)
	}

	var err error
	DB, err = sql.Open("postgres", CFG.DBConnectionString)

	if err != nil {
		log.Fatal(err)
	}
}
