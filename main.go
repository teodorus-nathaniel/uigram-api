package main

import (
	"context"
	"fmt"
	_ "fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/teodorus-nathaniel/ui-share-api/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initializeRoutes(router *gin.RouterGroup) {
	routes.UsePostsRoutes(router)
}

func getDatabaseConnection() string {
	databaseConn := os.Getenv("DATABASE")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")

	databaseConn = strings.Replace(databaseConn, "<username>", databaseUsername, 1)
	databaseConn = strings.Replace(databaseConn, "<password>", databasePassword, 1)

	return databaseConn
}

func connectToDatabase() (*mongo.Client, error) {
	databaseConn := getDatabaseConnection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(
		databaseConn,
	))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env", err.Error())
	}

	client, err := connectToDatabase()
	if err != nil {
		log.Fatal("Error connecting to Database", err.Error())
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	routerGroup := router.Group("/api/v1")
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")

	router.Run()
}
