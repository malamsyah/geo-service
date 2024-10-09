package main

import (
	"fmt"

	"github.com/malamsyah/geo-service/internal/db"
	"github.com/malamsyah/geo-service/internal/handler"
	"github.com/malamsyah/geo-service/pkg/config"
	"gorm.io/gorm"
)

func main() {
	var dbConn *gorm.DB
	var err error

	fmt.Println("Starting server...")

	conf := config.Instance()
	dbConn, err = db.ConnectPostgres(conf)
	if err != nil {
		panic(err)
	}

	err = db.Migrate(dbConn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Setting up router...")
	r := handler.SetupRouter(conf, dbConn)

	err = r.Run(":" + config.Instance().AppPort)
	if err != nil {
		panic(err)
	}
}
