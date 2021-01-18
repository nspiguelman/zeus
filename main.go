package main

import (
	"github.com/nspiguelman/zeus/controllers"
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/rest"
	"log"
)

func main() {
	kahootController := controllers.NewKahootController()

	// dsn := "host=full_db_postgres user=zeus_db password=password dbname=zeus port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	d := data.New()

	sqlDB, _ := d.DB.DB()

	if err := sqlDB.Ping(); err != nil {
		log.Panic(err.Error())
	}

	server := rest.NewServer(&kahootController)
	server.StartServer()
	// TODO: free resources if panic
}
