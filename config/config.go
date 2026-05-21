package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Cwd                string
	EnvMode            string
	ServiceName        string
	Host               string
	Port               int
	BasePath           string
	Timezone           string
	DBURI              string
	DBPassword         string
	EncryptionKey      string
	JWTPrivateKey      string
	JWTPublicKey       string
	UploadLimitSize    int
	ImageCompressLevel int
	RedisURI           string
	MQTTHost           string
	MQTTPort           int
	MQTTProtocol       string
	MQTTUser           string
	MQTTPassword       string
	MQTTPath           string
	MQTTTopic          string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s: %+v", err.Error(), err)
	}
	return &Config{
		Cwd:                cwd,
		EnvMode:            getEnv("ENV_MODE", "development"),
		ServiceName:        getEnv("SERVICE_NAME", "demo-rest-api"),
		Host:               getEnv("HOST", "0.0.0.0"),
		Port:               getEnvAsInt("PORT", 8000),
		BasePath:           getEnv("BASE_PATH", "api"),
		Timezone:           getEnv("TZ", "Asia/Bangkok"),
		DBURI:              getEnv("DB_URI", ""),
		DBPassword:         getEnv("DB_PASS", ""),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", ""),
		JWTPrivateKey:      getEnv("JWT_PRIVATE_KEY", ""),
		JWTPublicKey:       getEnv("JWT_PUBLIC_KEY", ""),
		UploadLimitSize:    getEnvAsInt("UPLOAD_LIMIT_SIZE", 10),
		ImageCompressLevel: getEnvAsInt("IMAGE_COMPRESS_LEVEL", 70),
		RedisURI:           getEnv("REDIS_URI", ""),
		MQTTHost:           getEnv("MQTT_HOST", ""),
		MQTTPort:           getEnvAsInt("MQTT_PORT", 8084),
		MQTTProtocol:       getEnv("MQTT_PROTOCOL", "wss"),
		MQTTUser:           getEnv("MQTT_USER", ""),
		MQTTPassword:       getEnv("MQTT_PASSWORD", ""),
		MQTTPath:           getEnv("MQTT_PATH", "mqtt"),
		MQTTTopic:          getEnv("MQTT_TOPIC", "demo"),
	}, nil
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
