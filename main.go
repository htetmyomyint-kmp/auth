package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/htetmyomyint-kmp/leave-tracker/auth/data"
	"github.com/htetmyomyint-kmp/leave-tracker/auth/handlers"
)

type Config struct {
	DB data.DBConfig
}

func main() {
	l := log.Default()

	conf := getConfig()

	dbClient := data.GetDBClient(l, conf.DB)
	log.Println(dbClient)
	if err := dbClient.InitDB(context.TODO(), conf.DB); err != nil {
		panic("cannot init database")
	}

	setUpServer(l, dbClient)
}

func setUpServer(l *log.Logger, dbClient data.UserDatabase) {
	r := gin.Default()

	apiGp := r.Group("/api")
	authHandler := handlers.NewGoogleAuthHandler(l, dbClient)

	apiGp.GET("/signup", authHandler.Signup)

	apiGp.GET("/login", authHandler.Login)

	apiGp.GET("/auth", authHandler.Auth)

	r.Run(":8080")
}

func getConfig() Config {
	return Config{
		DB: data.DBConfig{
			DBName:     "mongodb",
			ConnString: "mongodb://localhost:27017",
		},
	}

}
