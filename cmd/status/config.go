package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/yushafro/03.12.2025/internal/config"
	"github.com/yushafro/03.12.2025/pkg/logger"
	"go.uber.org/zap"
)

func initConfig(ctx context.Context) *config.Config {
	log := logger.FromContext(ctx)

	viper.SetConfigName("local")

	viper.AddConfigPath("config/")
	viper.AddConfigPath(".")

	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("PORT", "8080")

	viper.SetDefault("CLIENT_TIMEOUT", "2s")

	viper.SetDefault("FILE_REPOSITORY_PATH", "status.json")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(ctx, "Fatal error config file", zap.Error(err))

		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	log.Info(ctx, "config file loaded")

	config := new(config.Config)
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(ctx, "Unable to decode config", zap.Error(err))

		panic(err)
	}

	return config
}
