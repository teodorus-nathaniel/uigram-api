package posts

import (
	"context"
	"fmt"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func processDerivativeData(post *Post) {
	post.LikeCount = len(post.Likes)
	post.DislikeCount = len(post.Dislikes)
	// post.Liked = cari itu ama saved jg ama disliked
}

func GetPosts(sort string, limit int, page int) ([]Post, error) {
	options := options.Find()
	options.SetSort(bson.D{
		primitive.E{Key: "timestamp", Value: -1},
	})
	options.SetSkip(int64(limit * (page - 1)))
	options.SetLimit(int64(limit))

	cursor, err := database.PostsCollection.Find(context.TODO(), bson.M{}, options)
	if err != nil {
		return nil, err
	}

	fmt.Println(sort, limit, page)
	posts := []Post{}
	defer cursor.Close(database.Context)
	for cursor.Next(database.Context) {
		var post Post
		cursor.Decode(&post)

		processDerivativeData(&post)

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPost(filter bson.M) (*Post, error) {
	res := database.Database.Collection("posts").FindOne(database.Context, filter)
	var post Post
	err := res.Decode(&post)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func InsertPost(document bson.M) (*Post, error) {
	res, err := database.Database.Collection("posts").InsertOne(database.Context, document)

	post, err2 := GetPost(bson.M{
		"_id": res.InsertedID,
	})

	if err != nil {
		return nil, err
	} else if err2 != nil {
		return nil, err2
	}

	return post, nil
}
