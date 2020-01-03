package hooks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

// HookConfig stores configuration needed to setup the hook
type HookConfig struct {
	Key      string
	Format   string
	App      string
	Host     string
	Password string
	Hostname string
	Port     int
	DB       int
	MaxSize  int
}

// RedisHook to sends logs to Redis server
type RedisHook struct {
	RedisClient    *redis.Client
	RedisHost      string
	RedisKey       string
	LogstashFormat string
	AppName        string
	Hostname       string
	RedisPort      int
	MaxSize        int
}

// NewRedisHook creates a hook to be added to an instance of logger
func NewRedisHook(config HookConfig) (*RedisHook, error) {

	hostPort := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     hostPort,
		Password: config.Password,
		DB:       config.DB,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}

	if config.Format != "v0" && config.Format != "v1" && config.Format != "access" {
		return nil, fmt.Errorf("unknown message format")
	}

	return &RedisHook{
		RedisHost:      hostPort,
		RedisClient:    client,
		RedisKey:       config.Key,
		LogstashFormat: config.Format,
		AppName:        config.App,
		Hostname:       config.Hostname,
		MaxSize:        config.MaxSize,
	}, nil

}

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	var msg interface{}

	switch hook.LogstashFormat {
	case "v0":
		msg = createV0Message(entry, hook.AppName, hook.Hostname)
	case "v1":
		msg = createV1Message(entry, hook.AppName, hook.Hostname)
	case "access":
		msg = createAccessLogMessage(entry, hook.AppName, hook.Hostname)
	default:
		fmt.Println("Invalid LogstashFormat")
	}

	// Marshal into json message
	js, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error creating message for REDIS: %s", err)
	}

	// send message
	_, err = hook.RedisClient.RPush(hook.RedisKey, js).Result()
	if err != nil {
		return fmt.Errorf("error sending message to REDIS: %s", err)
	}
	if hook.MaxSize > 0 {
		_, err = hook.RedisClient.LTrim(
			hook.RedisKey, (int64)(-hook.MaxSize), -1).Result()
		if err != nil {
			return fmt.Errorf("error trimming list: %s", err)
		}
	}

	return nil
}

// Levels returns the available logging levels.
func (hook *RedisHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func createV0Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["@source_host"] = hostname
	m["@message"] = entry.Message

	fields := make(map[string]interface{})
	fields["level"] = entry.Level.String()
	fields["application"] = appName

	for k, v := range entry.Data {
		fields[k] = v
	}
	m["@fields"] = fields

	return m
}

func createV1Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["host"] = hostname
	m["message"] = entry.Message
	m["level"] = entry.Level.String()
	m["application"] = appName
	for k, v := range entry.Data {
		m[k] = v
	}

	return m
}

func createAccessLogMessage(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["message"] = entry.Message
	m["@source_host"] = hostname

	fields := make(map[string]interface{})
	fields["application"] = appName

	for k, v := range entry.Data {
		fields[k] = v
	}
	m["@fields"] = fields

	return m
}
