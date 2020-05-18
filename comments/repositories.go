package comments

import (
	"errors"
	"fmt"
	"time"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCommentsDataFromCursor(cursor *mongo.Cursor, user *users.User) []Comment {
	comments := []Comment{}
	defer cursor.Close(database.Context)
	for cursor.Next(database.Context) {
		var comment Comment
		cursor.Decode(&comment)
		comment.deriveData(user)

		comments = append(comments, comment)
	}

	return comments
}

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

func getCommentsByPostId(postID string, sort string, limit, page int, user *users.User) ([]Comment, error) {
	opts := getFindQueryOptions(sort, limit, page)
	fmt.Println(opts)

	cursor, err := database.CommentsCollection.Find(database.Context, bson.M{"postId": postID, "parent": nil}, opts)
	if err != nil {
		return nil, err
	}

	comments := getCommentsDataFromCursor(cursor, user)

	return comments, nil
}

func getCommentsRepliesByParent(parentID string, sort string, limit, page int, user *users.User) ([]Comment, error) {
	opts := getFindQueryOptions(sort, limit, page)

	cursor, err := database.CommentsCollection.Find(database.Context, bson.M{
		"parent": parentID,
	}, opts)

	if err != nil {
		return nil, err
	}

	replies := getCommentsDataFromCursor(cursor, user)
	return replies, nil
}

func getCommentsRepliesCount(parentID string) int64 {
	repliesCount, err := database.CommentsCollection.CountDocuments(database.Context, bson.M{
		"parent": parentID,
	})
	if err != nil {
		return 0
	}

	return repliesCount
}

func GetCommentsCount(postID string) int64 {
	count, err := database.CommentsCollection.CountDocuments(database.Context, bson.M{"postId": postID})
	if err != nil {
		return 0
	}

	return count
}

func getComment(commentID string, user *users.User) (*Comment, error) {
	oid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, err
	}

	var comment Comment
	err = database.CommentsCollection.FindOne(database.Context, primitive.M{"_id": oid}).Decode(&comment)
	if err != nil {
		return nil, err
	}

	comment.deriveData(user)

	return &comment, nil
}

func updateCommentLikes(id string, like, dislike *bool, user *users.User) (*Comment, error) {
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

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID is not valid")
	}

	if action == true {
		database.CommentsCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
			"$pull": primitive.M{
				otherAttribute: user.ID.Hex(),
			},
		})
		database.CommentsCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
			"$push": primitive.M{
				attribute: primitive.M{
					"$each":     primitive.A{user.ID.Hex()},
					"$position": 0,
				},
			},
		})
	} else {
		database.CommentsCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
			"$pull": primitive.M{
				attribute: user.ID.Hex(),
			},
		})
	}

	return getComment(id, user)
}

func insertComment(postID, content string, user *users.User) (*Comment, error) {
	res, err := database.CommentsCollection.InsertOne(database.Context, primitive.M{
		"content":   content,
		"postId":    postID,
		"userId":    user.ID.Hex(),
		"likes":     []string{},
		"dislikes":  []string{},
		"timestamp": time.Now().Unix(),
	})

	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID)
	return getComment(id.Hex(), user)
}

func insertReply(postID, parentID, content string, user *users.User) (*Comment, error) {
	res, err := database.CommentsCollection.InsertOne(database.Context, primitive.M{
		"content":   content,
		"postId":    postID,
		"parent":    parentID,
		"userId":    user.ID.Hex(),
		"likes":     []string{},
		"dislikes":  []string{},
		"timestamp": time.Now().Unix(),
	})

	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID)
	return getComment(id.Hex(), user)
}
