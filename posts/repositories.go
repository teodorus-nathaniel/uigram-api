package posts

import (
	"errors"

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

	var cursor *mongo.Cursor
	var err error

	var notEqualUser *bson.M
	if user == nil {
		notEqualUser = &bson.M{}
	} else {
		notEqualUser = &bson.M{
			"userId": bson.M{
				"$ne": user.ID.Hex(),
			},
		}
	}

	if sort == "likesCount" {
		allowDiskUse := true
		aggregateOpts := options.Aggregate()
		aggregateOpts.AllowDiskUse = &allowDiskUse

		cursor, err = database.PostsCollection.Aggregate(database.Context, bson.A{
			bson.M{
				"$match": notEqualUser,
			},
			bson.M{
				"$addFields": bson.M{
					"likesCount": bson.M{
						"$size": bson.M{
							"$ifNull": bson.A{
								"$likes", bson.A{},
							},
						},
					},
				},
			},
			bson.M{
				"$sort": bson.M{
					"likesCount": -1,
					"timestamp":  -1,
				},
			},
			bson.M{
				"$skip": opts.Skip,
			},
			bson.M{
				"$limit": opts.Limit,
			},
		},
			aggregateOpts)
	} else {
		cursor, err = database.PostsCollection.Find(database.Context, notEqualUser, opts)
	}
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

func insertPost(images []string, title, description, link string, timestamp int64, user *users.User) (*Post, error) {
	res, err := database.PostsCollection.InsertOne(database.Context, bson.M{
		"userId":      user.ID.Hex(),
		"title":       title,
		"likes":       []string{},
		"dislikes":    []string{},
		"description": description,
		"link":        link,
		"images":      images,
		"timestamp":   timestamp,
	})
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

	ids := []primitive.ObjectID{}
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

func updatePostLikes(user *users.User, id string, like, dislike *bool) (*Post, error) {
	post, err := getPost(id, user)
	if err != nil {
		return nil, err
	}

	var attribute string
	var otherAttribute string
	var action bool
	if like == nil {
		attribute = "dislikes"
		otherAttribute = "likes"
		action = *dislike
	} else if dislike == nil {
		attribute = "likes"
		otherAttribute = "dislikes"
		action = *like
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID is not valid")
	}

	if action == true {
		database.PostsCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
			"$pull": primitive.M{
				otherAttribute: user.ID.Hex(),
			},
		})
		database.PostsCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
			"$push": primitive.M{
				attribute: primitive.M{
					"$each":     primitive.A{user.ID.Hex()},
					"$position": 0,
				},
			},
		})
	} else {
		database.PostsCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
			"$pull": primitive.M{
				attribute: user.ID.Hex(),
			},
		})
	}

	post, err = getPost(id, user)
	return post, err
}
