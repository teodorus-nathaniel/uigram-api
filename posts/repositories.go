package posts

import (
	"time"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getFindQueryOptions(sort string, limit, page int) *options.FindOptions {
	if sort == "" {
		sort = "timestamp"
	}

	opts := options.Find()
	opts.SetSort(bson.D{
		primitive.E{Key: sort, Value: -1},
	})
	opts.SetSkip(int64(limit * (page - 1)))
	opts.SetLimit(int64(limit))

	return opts
}

func getPostsDataFromCursor(cursor *mongo.Cursor, user *users.User) []Post {
	posts := []Post{}
	defer cursor.Close(database.Context)
	for cursor.Next(database.Context) {
		var post Post
		cursor.Decode(&post)
		post.deriveToPost(user)

		posts = append(posts, post)
	}

	return posts
}

func getPosts(sort string, limit int, page int, user *users.User) ([]Post, error) {
	opts := getFindQueryOptions(sort, limit, page)

	cursor, err := database.PostsCollection.Find(database.Context, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	posts := getPostsDataFromCursor(cursor, user)

	return posts, nil
}

func getPostsByUserId(ids []string, limit, page int, user *users.User) ([]Post, error) {
	opts := getFindQueryOptions("", limit, page)

	cursor, err := database.PostsCollection.Find(database.Context, bson.M{
		"userId": bson.M{
			"$in": ids,
		},
	}, opts)

	if err != nil {
		return nil, err
	}

	posts := getPostsDataFromCursor(cursor, user)

	return posts, nil
}

func getPost(id string, user *users.User) (*Post, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var post Post
	err = database.Database.Collection("posts").FindOne(database.Context, primitive.M{"_id": oid}).Decode(&post)
	if err != nil {
		return nil, err
	}

	post.deriveToPostDetail(user)

	return &post, nil
}

func GetPostByOwner(id, sort string, limit, page int, user *users.User) ([]Post, error) {
	opts := getFindQueryOptions(sort, limit, page)

	cursor, err := database.PostsCollection.Find(database.Context, primitive.M{"userId": id}, opts)
	if err != nil {
		return nil, err
	}

	posts := getPostsDataFromCursor(cursor, user)

	return posts, nil
}

func insertPost(document Post, user *users.User) (*Post, error) {
	document.ID = primitive.NilObjectID
	document.Timestamp = time.Now().Unix()

	res, err := database.Database.Collection("posts").InsertOne(database.Context, document)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID)
	post, err := getPost(id.Hex(), user)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func getSavedPosts(sort string, limit, page int, user *users.User) ([]Post, error) {
	opts := getFindQueryOptions(sort, limit, page)

	var ids []primitive.ObjectID
	for _, id := range user.SavedPosts {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			ids = append(ids, oid)
		}
	}

	cursor, err := database.PostsCollection.Find(database.Context, bson.M{
		"_id": bson.M{
			"$in": ids,
		}}, opts)

	if err != nil {
		return nil, err
	}

	posts := getPostsDataFromCursor(cursor, user)

	return posts, nil
}
