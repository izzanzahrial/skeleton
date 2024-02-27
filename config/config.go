package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Database struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
	Timezone string
}

func (d *Database) URL() string {
	// return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=%s", d.Username, d.Password, d.Host, d.Port, d.Name, d.Timezone)
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", d.Host, d.Username, d.Password, d.Name, d.Port, d.Timezone)
}

type Cache struct {
	Address  string
	Username string
	Password string
	Port     string
	DB       string
}

func (c *Cache) URL() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", c.Username, c.Password, c.Address, c.Port, c.DB)
}

func NewDatabase() (*Database, error) {
	var db Database
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, errors.New("environment DB_NAME must be set")
	}
	db.Name = dbName

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, errors.New("environment DB_HOST must be set")
	}
	db.Host = dbHost

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, errors.New("environment DB_PORT must be set")
	}
	db.Port = dbPort

	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		return nil, errors.New("environment DB_USERNAME must be set")
	}
	db.Username = dbUsername

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, errors.New("environment DB_PASSWORD must be set")
	}
	db.Password = dbPassword

	dbTimezone := os.Getenv("DB_TIMEZONE")
	if dbTimezone == "" {
		return nil, errors.New("environment DB_TIMEZONE must be set")
	}
	db.Timezone = dbTimezone

	return &db, nil
}
func NewCache() (*Cache, error) {
	var ch Cache
	cacheAddress := os.Getenv("CACHE_ADDRESS")
	if cacheAddress == "" {
		return nil, errors.New("environment CACHE_ADDRESS must be set")
	}
	ch.Address = cacheAddress

	cacheUsername := os.Getenv("CACHE_USERNAME")
	if cacheUsername == "" {
		return nil, errors.New("environment CACHE_USERNAME must be set")
	}
	ch.Username = cacheUsername

	cachePassword := os.Getenv("CACHE_PASSWORD")
	if cachePassword == "" {
		return nil, errors.New("environment CACHE_PASSWORD must be set")
	}
	ch.Password = cachePassword

	cachePort := os.Getenv("CACHE_PORT")
	if cachePort == "" {
		return nil, errors.New("environment CACHE_PORT must be set")
	}
	ch.Port = cachePort

	cacheDB := os.Getenv("CACHE_DB")
	if cacheDB == "" {
		return nil, errors.New("environment CACHE_DB must be set")
	}
	ch.DB = cacheDB

	return &ch, nil
}

type Producer struct {
	// Version           [4]uint
	// FlushBytes        int
	Addresses         []string
	Timeout           time.Duration
	ChannelBufferSize int
}

func NewProducer() (*Producer, error) {
	addrsString := os.Getenv("KAFKA_ADDRESSES")
	addresses := strings.Split(addrsString, ",")

	// versionString := os.Getenv("KAFKA_VERSION")
	// versionSlice := strings.Split(versionString, ".")
	// if len(versionSlice) != 4 {
	// 	return nil, errors.New("environment KAFKA_VERSION must be combination of 4 uint")
	// }

	// var version [4]uint
	// for i, v := range versionSlice {
	// 	ui64, err := strconv.ParseUint(v, 10, 64)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to parse version slice to uint: %w", err)
	// 	}

	// 	version[i] = uint(ui64)
	// }

	timeoutString := os.Getenv("KAFKA_TIMEOUT")
	if timeoutString == "" {
		return nil, errors.New("environment KAFKA_TIMEOUT must be set")
	}
	timeout, err := strconv.Atoi(timeoutString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeout string to int: %w", err)
	}

	// flushBytesString := os.Getenv("KAFKA_FLUSH_BYTES")
	// if flushBytesString == "" {
	// 	return nil, errors.New("environment KAFKA_FLUSH_BYTES must be set")
	// }
	// flushBytes, err := strconv.Atoi(flushBytesString)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse flush bytes string to int: %w", err)
	// }

	channelBufferSizeString := os.Getenv("KAFKA_CHANNEL_BUFFER_SIZE")
	channelBufferSize, err := strconv.Atoi(channelBufferSizeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse channel buffer size string to int: %w", err)
	}

	return &Producer{
		// Version:           version,
		// FlushBytes:        flushBytes,
		Addresses:         addresses,
		Timeout:           time.Duration(timeout) * time.Second,
		ChannelBufferSize: channelBufferSize,
	}, nil
}

type Consumer struct {
	// Version           [4]uint
	Addresses         []string
	Topics            []string
	GroupID           string
	AutoCommit        bool
	FetchBytes        int32
	MaxWait           time.Duration
	ChannelBufferSize int
}

func NewConsumer() (*Consumer, error) {
	addrsString := os.Getenv("KAFKA_ADDRESSES")
	addresses := strings.Split(addrsString, ",")

	topicsString := os.Getenv("KAFKA_TOPICS")
	if topicsString == "" {
		return nil, errors.New("environment KAFKA_TOPICS must be set")
	}
	topics := strings.Split(topicsString, ",")

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		return nil, errors.New("environment KAFKA_GROUP_ID must be set")
	}

	offsetsAutoCommitBoolStr := os.Getenv("KAFKA_OFFSETS_AUTOCOMMIT")
	if offsetsAutoCommitBoolStr == "" {
		return nil, errors.New("environment KAFKA_OFFSETS_AUTOCOMMIT must be set")
	}
	offsetsAutoCommitBool, err := strconv.ParseBool(offsetsAutoCommitBoolStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OffsetsAutoCommitBoolStr to bool: %w", err)
	}

	fetchBytesString := os.Getenv("KAFKA_FETCH_BYTES")
	if fetchBytesString == "" {
		return nil, errors.New("environment KAFKA_FETCH_BYTES must be set")
	}
	fetchBytesInt, err := strconv.Atoi(fetchBytesString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fetch bytes string to int: %w", err)
	}
	fetchBytes := int32(fetchBytesInt)

	maxWaitString := os.Getenv("KAFKA_MAX_WAIT")
	if maxWaitString == "" {
		return nil, errors.New("environment KAFKA_MAX_WAIT must be set")
	}
	maxWait, err := strconv.Atoi(maxWaitString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse max wait string to int: %w", err)
	}

	channelBufferSizeString := os.Getenv("KAFKA_CHANNEL_BUFFER_SIZE")
	channelBufferSize, err := strconv.Atoi(channelBufferSizeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse channel buffer size string to int: %w", err)
	}

	return &Consumer{
		// Version:           version,
		Addresses:         addresses,
		Topics:            topics,
		GroupID:           groupID,
		AutoCommit:        offsetsAutoCommitBool,
		FetchBytes:        fetchBytes,
		MaxWait:           time.Duration(maxWait) * time.Millisecond,
		ChannelBufferSize: channelBufferSize,
	}, nil
}
