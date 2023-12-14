package data

import (
	"context"
	"log"
)

const (
	mongoDB  = "mongodb"
	postgres = "postgres"
	mySQL    = "mySQL"
)

type UserDatabase interface {
	InitDB(context.Context, DBConfig) error
	GetUserByID(id string) (User, error)
	CreateUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
}

type User struct {
	Name       string
	Email      string
	ProfileURL string
}

type DBConfig struct {
	DBName     string
	ConnString string
}

func GetDBClient(l *log.Logger, conf DBConfig) UserDatabase {
	switch conf.DBName {
	case mongoDB:
		return NewMongoDBClient(l)
	default:
		return nil
	}
}
