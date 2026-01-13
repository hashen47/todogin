package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr       string
	DBAddr           string
	DBUser           string
	DBPasswd         string 
	DBName           string
	JwtTokenLifetime string
	JwtSecretKey     string
}

func ConfigInit() (*Config, error) {
	c := new(Config)

	vals, err := godotenv.Read(".env")
	if err != nil {
		return nil, &ConfError{".env load fail", err}
	}

	// load keys
	var val string
	if val, err = getVal(&vals, "ServerAddr"); err != nil {
		return nil, err
	}
	c.ServerAddr = val

	if val, err = getVal(&vals, "DBAddr"); err != nil {
		return nil, err
	}
	c.DBAddr = val

	if val, err = getVal(&vals, "DBUser"); err != nil {
		return nil, err
	}
	c.DBUser = val

	if val, err = getVal(&vals, "DBPasswd"); err != nil {
		return nil, err
	}
	c.DBPasswd = val

	if val, err = getVal(&vals, "DBName"); err != nil {
		return nil, err
	}
	c.DBName = val

	if val, err = getVal(&vals, "JwtTokenLifetime"); err != nil {
		return nil, err
	}
	c.JwtTokenLifetime = val

	if val, err = getVal(&vals, "JwtSecretKey"); err != nil {
		return nil, err
	}
	c.JwtSecretKey = val

	return c, nil
}

func getVal(vals *map[string]string, key string) (string, error) {
	if val, ok := (*vals)[key]; ok {
		return val, nil
	}
	return "", &ConfError{fmt.Sprintf("%q key is not found", key), nil}
}
