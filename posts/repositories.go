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

func GetPost(filter bson.M) (*Post, error) {
	res := database.Database.Collection("posts").FindOne(database.Context, filter)
	var post Post
	err := res.Decode(&post)

	if err != nil {
		return nil , err
	}

	return &post, nil
}

func InsertPost(document bson.M)(*Post, error) {
	res, err := database.Database.Collection("posts").InsertOne(database.Context, document)

	post, err2 := GetPost(bson.M{
		"_id" : res.InsertedID,
	})

	if err != nil {
		return nil, err
	}else if err2 != nil{
		return nil, err2
	}

	return post, nil
}