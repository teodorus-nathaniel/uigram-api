package posts

import (
	"github.com/teodorus-nathaniel/uigram-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPosts(filters bson.M) []Post {
	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "timestamp", Value: -1}})
	//TODO: options.SetLimit(10), INI DARI URL QUERY
	cursor, _ := database.Database.Collection("posts").Find(database.Context, filters, options)
	var posts []Post
	defer cursor.Close(database.Context)
	for cursor.Next(database.Context) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	return posts
}
