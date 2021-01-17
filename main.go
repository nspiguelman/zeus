package main

import (
	"github.com/nspiguelman/zeus/controllers"
	"github.com/nspiguelman/zeus/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	kahootController := controllers.NewKahootController()

	dsn := "host=full_db_postgres user=zeus_db password=password dbname=zeus port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if db == nil {
		panic(err.Error())
	}

	server := rest.NewServer(&kahootController, db)
	server.StartServer()
}
