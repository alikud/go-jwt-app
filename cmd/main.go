package main

import (
	"backend/pkg/handler"
	"backend/pkg/repository/postgres"
	"backend/pkg/service"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func newConfig() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	err := newConfig()
	if err != nil {
		log.Panic(err)
		panic("Error with config!")
	}

	pool, _, err := postgres.NewPostgresDB(postgres.PostgresConfig{
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSL_MODE"),
		MinConns: 10,
		MaxConns: 20,
		TimeOut:  5,
	},
	)
	if err != nil {
		panic("Error with creating postgres db")
	}

	repo := postgres.NewRepository(pool)
	services := service.NewService(repo)
	handler := handler.NewHandler(services)

	err = http.ListenAndServe(":8080", handler.InitRoutes())
	if err != nil {
		panic("Cant serve at this port")
	}
}
