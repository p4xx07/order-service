//go:build wireinject
// +build wireinject

package deps

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/meilisearch/meilisearch-go"
	"github.com/p4xx07/order-service/app"
	"github.com/p4xx07/order-service/app/domains/inventory"
	"github.com/p4xx07/order-service/app/domains/order"
	"github.com/p4xx07/order-service/app/domains/product"
	"github.com/p4xx07/order-service/app/domains/user"
	"github.com/p4xx07/order-service/configuration"
	"github.com/p4xx07/order-service/internal/db"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

func InjectApp(config *configuration.Configuration, logger *zap.SugaredLogger) (*app.App, error) {
	wire.Build(
		InitMeiliSearchClient,
		InitRedisClient,

		// handlers
		order.NewHandler,

		// services
		order.NewService,
		order.NewMeilisearchService,
		inventory.NewService,

		// stores
		ConnectDB,
		order.NewStore,
		inventory.NewStore,

		wire.Struct(new(app.App), "*"),
	)

	return nil, nil
}

func ConnectDB(configuration *configuration.Configuration) (*gorm.DB, error) {
	database, err := db.ConnectDB(configuration)
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(
		user.User{},
		product.Product{},
		inventory.Inventory{},
		order.Order{},
		order.OrderItem{},
	)

	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	}

	return database, nil
}

func InitRedisClient(configuration *configuration.Configuration) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     configuration.RedisHost + ":" + configuration.RedisPort,
		Password: configuration.RedisPassword,
		DB:       configuration.RedisDatabase,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	fmt.Println("Redis connected successfully")
	return client, nil
}

func InitMeiliSearchClient(configuration *configuration.Configuration) (meilisearch.ServiceManager, error) {
	host := fmt.Sprintf("%s:%d", configuration.MeiliSearchHost, configuration.MeiliSearchPort)
	client := meilisearch.New(
		host,
		meilisearch.WithAPIKey(configuration.MeiliSearchMasterKey),
	)
	index := client.Index("orders")
	sortable := []string{"created_at"}
	_, err := index.UpdateSortableAttributes(&sortable)
	if err != nil {
		return nil, err
	}
	return client, nil
}
