package databases

import (
	"fmt"
	"log"
	"sync"

	"github.com/jakkaphatminthana/isekai-shop-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	*gorm.DB
}

var (
	// Instance of postgresDatabase
	postgresDatabaseInstance *postgresDatabase
	// Singleton pattern (call once)
	once sync.Once
)

func NewPostgresDatabase(conf *config.Database) Database {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.DBName, conf.SSLMode, conf.Schema)

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			panic(err)
		}
		log.Printf("Connected to database %s", conf.DBName)

		postgresDatabaseInstance = &postgresDatabase{conn}
	})

	return postgresDatabaseInstance
}

// implement Database interface
func (db *postgresDatabase) Connect() *gorm.DB {
	return postgresDatabaseInstance.DB
}
