package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger Logger
	Server ServerConfig
	Kafka  KafkaConfig
	Email  EmailConfig
}

type Logger struct {
	Level string `env:"LOG_LEVEL,required"`
}

type ServerConfig struct {
	Addr string `env:"SERVER_ADDR,required"`
}

type KafkaConfig struct {
	Brokers   []string `env:"KAFKA_BROKERS" envSeparator:"," required:"true"` // Список брокеров Kafka
	Topic     string   `env:"KAFKA_TOPIC,required"`                           // Название топика
	GroupID   string   `env:"KAFKA_GROUP_ID,required"`                        // Идентификатор группы потребителей
	Username  string   `env:"KAFKA_USERNAME"`                                 // Имя пользователя (если нужна авторизация)
	Password  string   `env:"KAFKA_PASSWORD"`                                 // Пароль (если нужна авторизация)
	UseTLS    bool     `env:"KAFKA_USE_TLS" envDefault:"false"`               // Использовать ли TLS
	TLSCACert string   `env:"KAFKA_TLS_CA_CERT"`                              // Путь к CA-сертификату (если нужен TLS)
	TLSCert   string   `env:"KAFKA_TLS_CERT"`                                 // Путь к клиентскому сертификату (если нужен TLS)
	TLSKey    string   `env:"KAFKA_TLS_KEY"`                                  // Путь к клиентскому ключу (если нужен TLS)
}

type EmailConfig struct {
	Username string `env:"EMAIL_USERNAME,required"` // Email-адрес отправителя
	Password string `env:"EMAIL_PASSWORD,required"`
	Host     string `env:"EMAIL_SMTP_HOST,required"`
	Port     string `env:"EMAIL_SMTP_PORT,required"`
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		_ = godotenv.Load()
		if err := env.Parse(&config); err != nil {
			log.Fatal(err)
		}
	})
	return &config
}
