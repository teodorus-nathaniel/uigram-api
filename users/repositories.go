package users

import (
	"errors"
	"fmt"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(filter primitive.M) (*User, error) {
	var user User
	err := database.UsersCollection.FindOne(database.Context, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	return GetUser(primitive.M{
		"email": email,
	})
}

func GetUserById(id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return GetUser(primitive.M{
		"_id": oid,
	})
}

func InsertUser(document *User) (*User, error) {
	document.fillEmptyValues()
	duplicateEmail, err := GetUserByEmail(document.Email)
	if duplicateEmail != nil {
		return nil, errors.New("email is already used")
	}

	err = document.hashPassword()
	if err != nil {
		return nil, err
	}

	res, err := database.UsersCollection.InsertOne(database.Context, document)
	if err != nil {
		return nil, err
	}

	user, err := GetUserById(res.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func AppendArray(id, attribute, data string) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	oidData, err := primitive.ObjectIDFromHex(data)
	if err != nil {
		return nil, err
	}
	res, err := database.UsersCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
		"$push": primitive.M{
			attribute: primitive.M{
				"$each":     primitive.A{oidData.Hex()},
				"$position": 0,
			},
		},
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteArrayElement(id, attribute, data string) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	oidData, err := primitive.ObjectIDFromHex(data)
	if err != nil {
		return nil, err
	}
	res, err := database.UsersCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
		"$pull": primitive.M{
			attribute: oidData.Hex(),
		},
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func AppendFollower(id, follower string) (*mongo.UpdateResult, error) {
	return AppendArray(id, "followers", follower)
}

func AppendFollowing(id, following string) (*mongo.UpdateResult, error) {
	return AppendArray(id, "following", following)
}

func PullFollower(id, follower string) (*mongo.UpdateResult, error) {
	return DeleteArrayElement(id, "followers", follower)
}

func PullFollowing(id, following string) (*mongo.UpdateResult, error) {
	return DeleteArrayElement(id, "following", following)
}

func AddSavedPostDatabase(id, postId string) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return database.UsersCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
		"$push": primitive.M{
			"savedPosts": primitive.M{
				"$each":     primitive.A{postId},
				"$position": 0,
			},
		},
	})
}

func DeleteSavedPostDatabase(id, postId string) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return database.UsersCollection.UpdateOne(database.Context, primitive.M{"_id": oid}, primitive.M{
		"$pull": primitive.M{
			"savedPosts": postId,
		},
	})
}

func getUsersDataFromCursor(cursor *mongo.Cursor, self *User) []User {
	users := []User{}
	defer cursor.Close(database.Context)
	for cursor.Next(database.Context) {
		var user User
		cursor.Decode(&user)
		user.DeriveAttributesAndHideCredentials(self)

		users = append(users, user)
	}

	return users
}

func getUserFollowers(id string, self *User) ([]User, error) {
	user, err := GetUserById(id)
	if err != nil {
		return nil, err
	}

	ids := []primitive.ObjectID{}
	for _, id := range user.Followers {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			ids = append(ids, oid)
		}
	}

	cursor, err := database.UsersCollection.Find(database.Context, primitive.M{
		"_id": primitive.M{
			"$in": ids,
		},
	})

	if err != nil {
		return nil, err
	}

	return getUsersDataFromCursor(cursor, self), nil
}

func getUserFollowing(id string, self *User) ([]User, error) {
	user, err := GetUserById(id)
	if err != nil {
		return nil, err
	}

	ids := []primitive.ObjectID{}
	for _, id := range user.Following {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			ids = append(ids, oid)
		}
	}

	cursor, err := database.UsersCollection.Find(database.Context, primitive.M{
		"_id": primitive.M{
			"$in": ids,
		},
	})

	if err != nil {
		return nil, err
	}

	return getUsersDataFromCursor(cursor, self), nil
}

func updateUser(id *primitive.ObjectID, username, fullname, status, imagePath string) (*mongo.UpdateResult, error) {
	updateObj := primitive.D{}
	if username != "" {
		updateObj = append(updateObj, primitive.E{Key: "username", Value: username})
	}
	if fullname != "" {
		updateObj = append(updateObj, primitive.E{Key: "fullname", Value: fullname})
	}
	if status != "" {
		updateObj = append(updateObj, primitive.E{Key: "status", Value: status})
	}
	if imagePath != "" {
		updateObj = append(updateObj, primitive.E{Key: "profilePicture", Value: imagePath})
	}

	fmt.Println(updateObj)

	return database.UsersCollection.UpdateOne(database.Context, primitive.M{"_id": id}, primitive.M{"$set": updateObj})
}
