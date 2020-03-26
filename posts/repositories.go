package posts

import (
	"context"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getPosts(sort string, limit int, page int) ([]Post, error) {
	opts := options.Find()
	opts.SetSort(bson.D{
		primitive.E{Key: "timestamp", Value: -1},
	})
	opts.SetSkip(int64(limit * (page - 1)))
	opts.SetLimit(int64(limit))

	cursor, err := database.PostsCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	defer cursor.Close(database.Context)
	for cursor.Next(database.Context) {
		var post Post
		cursor.Decode(&post)
		post.deriveToPost()

		posts = append(posts, post)
	}

	return posts, nil
}

func getPost(id string) (*Post, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var post Post
	err = database.Database.Collection("posts").FindOne(database.Context, primitive.M{"_id": oid}).Decode(&post)
	if err != nil {
		return nil, err
	}

	post.deriveToPostDetail()

	return &post, nil
}

func insertPost(document Post) (*Post, error) {
	res, err := database.Database.Collection("posts").InsertOne(database.Context, document)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID)
	post, err := getPost(id.String())
	if err != nil {
		return nil, err
	}

	return post, nil
}
