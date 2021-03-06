package repository

import (
	"fmt"
	"go-shorterer/model"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB model
type dbm struct {
	PsqlDB *gorm.DB
	ID     uint
}

type DB interface {
	SaveShortlink(*model.ShortLink, bool) error
	DeleteShortlink(string) error
	GetDestination(string) (string, error)
	SaveUser(*model.User) error
	CheckAPIKey(string) bool
}

// MakeDB create new DB object
func MakeDB() (DB, error) {
	DBMS := os.Getenv("DB_DRIVER")
	DBURL := generateDBURL()

	var (
		db  *gorm.DB
		err error
	)

	if db, err = gorm.Open(DBMS, DBURL); err != nil {
		return nil, err
	}

	// defer db.Close()

	if err = db.Debug().DropTableIfExists(&model.ShortLink{}, &model.User{}).Error; err != nil {
		return nil, err
	}
	if err = db.Debug().AutoMigrate(&model.ShortLink{}, &model.User{}).Error; err != nil {
		return nil, err
	}

	Db := &dbm{
		PsqlDB: db,
	}

	return Db, nil
}

func generateDBURL() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, name, password)

	// dbUrl := os.Getenv("DATABASE_URL")
	// return dbUrl
}
