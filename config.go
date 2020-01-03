package main

import (
	"io/ioutil"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

var globalConfig *config

// const value
const (
	DefaultConfigFile = "config.yaml"
	DefaultListenIP   = "0.0.0.0"
	DefaultListenPort = 8000
	DefaultRedisLoggerMaxSize = 0
)

type config struct {
	LogLevel string
	Service struct {
		ServiceName string
		DisplayName string
		Description string
	}
	Listener struct {
		IP   string
		Port int
	}
	RedisLogger struct {
		Enabled  bool
		Addr     string
		Port     int
		Password string
		Key      string
		DB       int
		MaxSize int
	}
}

func readConfig(fileName string) (*config, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	c := &config{}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return nil, fmt.Errorf(
			"Parse config file %s fail: %s", fileName, err.Error())
	}
	if c.RedisLogger.Enabled {
		if c.RedisLogger.Addr == "" {
			return nil, fmt.Errorf("Must provide the address of redis server")
		}
		if c.RedisLogger.Port < 1 || c.RedisLogger.Port > 65534 {
			return nil, fmt.Errorf("Invalid redis server port.")
		}
		if c.RedisLogger.Key == "" {
			return nil, fmt.Errorf("Must provide a list key in redis")
		}
		if c.RedisLogger.DB < 0 || c.RedisLogger.DB > 15 {
			return nil, fmt.Errorf("Invalid redis DB")
		}
		if c.RedisLogger.MaxSize < 0 {
			c.RedisLogger.MaxSize = 0
		}
	}
	return c, nil
}
