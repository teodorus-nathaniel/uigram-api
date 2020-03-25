package users

import (
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credentials struct {
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

type User struct {
	ID         *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Password   string              `json:"password" bson:"password"`
	Email      string              `json:"email" bson:"email"`
	Username   string              `json:"username" bson:"username"`
	ProfilePic string              `json:"profilePicture,omitempty" bson:"profilePicture"`
	Followers  []string            `json:"followers" bson:"followers"`
	Following  []string            `json:"following" bson:"following"`
}

func (user *User) fillEmptyValues() {
	if user.Followers == nil {
		user.Followers = []string{}
	}
	if user.Following == nil {
		user.Following = []string{}
	}
}

func ValidateEmailPassword(email, password string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		return errors.New("invalid email")
	}
	if len(password) < 5 {
		return errors.New("password must be more than 5 characters")
	}

	return nil
}

func (data *Credentials) ValidateEmailPassword() error {
	return ValidateEmailPassword(data.Email, data.Password)
}

func (data *User) ValidateData() error {
	err := ValidateEmailPassword(data.Email, data.Password)
	if err != nil {
		return err
	}

	if len(data.Username) < 6 || len(data.Username) > 20 {
		return errors.New("username must be between 6 and 20 characters")
	}

	return nil
}
