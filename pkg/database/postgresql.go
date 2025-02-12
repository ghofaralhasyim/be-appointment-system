package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/pkg/utils"
	_ "github.com/lib/pq"
)

func InitDbConnection() (*sql.DB, error) {
	connUrl, errConnUrl := utils.ConnURLBuilder("postgres")
	if errConnUrl != nil {
		return nil, errConnUrl
	}

	db, err := sql.Open("postgres", connUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	log.Println("check: db postgresql connected")
	return db, nil
}
