package users

import (
	"errors"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
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