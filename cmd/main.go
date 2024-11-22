package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/vitalii-minchuk/alaska/cmd/api"
	"github.com/vitalii-minchuk/alaska/config"
	"github.com/vitalii-minchuk/alaska/db"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	initStorage(db)
	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	log.Println("DB successfully connected!")
}
