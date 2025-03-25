package db

import (
	"fmt"
	"github.com/p4xx07/order-service/configuration"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDB *gorm.DB

func ConnectDB(configuration *configuration.Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configuration.DatabaseUsername,
		configuration.DatabasePassword,
		configuration.DatabaseHost,
		configuration.DatabasePort,
		configuration.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
