package database

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func getDatabaseConnection() string {
	databaseConn := os.Getenv("DATABASE")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")

	databaseConn = strings.Replace(databaseConn, "<username>", databaseUsername, 1)
	databaseConn = strings.Replace(databaseConn, "<password>", databasePassword, 1)

	return databaseConn
}

func InitializeClient() {
	databaseConn := getDatabaseConnection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		databaseConn,
	))

	if err != nil {
		log.Fatal("Error connecting to database...", err.Error())
	}
	Client = client
}
