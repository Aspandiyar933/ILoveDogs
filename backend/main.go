package main

import (
	"log"

	"github.com/Aspandiyar933/Ilovedogs/api"
	"github.com/Aspandiyar933/Ilovedogs/config"
	"github.com/Aspandiyar933/Ilovedogs/database"
	"github.com/Aspandiyar933/Ilovedogs/store"
	"github.com/go-sql-driver/mysql"
)

func main() {

	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sqlStorage := database.NewMySQLStorage(cfg)
	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(db)

	api := api.NewAPIServer(":3000", store)
	api.Serve()
}
