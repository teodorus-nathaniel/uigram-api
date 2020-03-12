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

var Database *mongo.Database
var Context context.Context

func getDatabaseConnection() string {
	databaseConn := os.Getenv("DATABASE_CONNECTION")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")

	databaseConn = strings.Replace(databaseConn, "<username>", databaseUsername, 1)
	databaseConn = strings.Replace(databaseConn, "<password>", databasePassword, 1)

	return databaseConn
}

func InitializeClient() {
	databaseConn := getDatabaseConnection()
	databaseName := os.Getenv("DATABASE")

	Context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(Context, options.Client().ApplyURI(
		databaseConn,
	))

	if err != nil {
		log.Fatal("Error connecting to database...", err.Error())
	}
	Database = client.Database(databaseName)
}
