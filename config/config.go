package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Port     string
	Database database
	Cache    cache
}

type database struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
	Timezone string
}

func (d *database) URL() string {
	// return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=%s", d.Username, d.Password, d.Host, d.Port, d.Name, d.Timezone)
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", d.Host, d.Username, d.Password, d.Name, d.Port, d.Timezone)
}

type cache struct {
	Address  string
	Username string
	Password string
	Port     string
	DB       string
}

func (c *cache) URL() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", c.Username, c.Password, c.Address, c.Port, c.DB)
}

func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("environment PORT must be set")
	}

	db, err := newDatabase()
	if err != nil {
		return nil, err
	}

	cache, err := newCache()
	if err != nil {
		return nil, err
	}

	return &Config{Port: port, Database: db, Cache: cache}, nil
}

func newDatabase() (database, error) {
	var db database
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return database{}, errors.New("environment DB_NAME must be set")
	}
	db.Name = dbName

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return database{}, errors.New("environment DB_HOST must be set")
	}
	db.Host = dbHost

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return database{}, errors.New("environment DB_PORT must be set")
	}
	db.Port = dbPort

	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		return database{}, errors.New("environment DB_USERNAME must be set")
	}
	db.Username = dbUsername

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return database{}, errors.New("environment DB_PASSWORD must be set")
	}
	db.Password = dbPassword

	dbTimezone := os.Getenv("DB_TIMEZONE")
	if dbTimezone == "" {
		return database{}, errors.New("environment DB_TIMEZONE must be set")
	}
	db.Timezone = dbTimezone

	return db, nil
}
func newCache() (cache, error) {
	var ch cache
	cacheAddress := os.Getenv("CACHE_ADDRESS")
	if cacheAddress == "" {
		return cache{}, errors.New("environment CACHE_ADDRESS must be set")
	}
	ch.Address = cacheAddress

	cacheUsername := os.Getenv("CACHE_USERNAME")
	if cacheUsername == "" {
		return cache{}, errors.New("environment CACHE_USERNAME must be set")
	}
	ch.Username = cacheUsername

	cachePassword := os.Getenv("CACHE_PASSWORD")
	if cachePassword == "" {
		return cache{}, errors.New("environment CACHE_PASSWORD must be set")
	}
	ch.Password = cachePassword

	cachePort := os.Getenv("CACHE_PORT")
	if cachePort == "" {
		return cache{}, errors.New("environment CACHE_PORT must be set")
	}
	ch.Port = cachePort

	cacheDB := os.Getenv("CACHE_DB")
	if cacheDB == "" {
		return cache{}, errors.New("environment CACHE_DB must be set")
	}
	ch.DB = cacheDB

	return ch, nil
}
