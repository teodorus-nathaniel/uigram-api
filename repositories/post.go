package repositories

import (
	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPosts(filters bson.M) []models.Post{
	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "timestamp", Value: -1}})
	//TODO: options.SetLimit(10), INI DARI URL QUERY
	cursor, _ := database.Database.Collection("posts").Find(database.Context, filters, options)
	var posts []models.Post
	defer cursor.Close(database.Context)
	for cursor.Next(database.Context) {
		var post models.Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	return posts
}