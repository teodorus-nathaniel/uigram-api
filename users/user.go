package users

import (
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

type User struct {
	ID             *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Password       *string             `json:"password,omitempty" bson:"password"`
	Email          string              `json:"email" bson:"email"`
	Username       string              `json:"username" bson:"username"`
	Fullname       string              `json:"fullname" bson:"fullname"`
	ProfilePic     string              `json:"profilePicture,omitempty" bson:"profilePicture"`
	Followers      []string            `json:"-" bson:"followers"`
	Following      []string            `json:"-" bson:"following"`
	FollowingCount int                 `json:"followingCount" bson:"-"`
	FollowersCount int                 `json:"followersCount" bson:"-"`
	Status         string              `json:"status" bson:"status"`
}

func (user *User) fillEmptyValues() {
	if user.Followers == nil {
		user.Followers = []string{}
	}
	if user.Following == nil {
		user.Following = []string{}
	}
}

func (user *User) hashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	pass := string(hash)
	user.Password = &pass
	return nil
}

func (user *User) IsRightPassword(password string) bool {
	bytePassToCheck := []byte(password)
	bytePass := []byte(*user.Password)
	err := bcrypt.CompareHashAndPassword(bytePass, bytePassToCheck)
	if err != nil {
		return false
	}

	return true
}

func (user *User) HideCredentials() {
	user.Password = nil
}

func (user *User) DeriveAttributesAndHideCredentials() {
	user.HideCredentials()
	user.FollowersCount = len(user.Followers)
	user.FollowingCount = len(user.Following)
}

func ValidateEmailPassword(email, password string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		return errors.New("invalid email format")
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
	err := ValidateEmailPassword(data.Email, *data.Password)
	if err != nil {
		return err
	}

	if len(data.Username) < 6 || len(data.Username) > 20 {
		return errors.New("username must be between 6 and 20 characters")
	}

	return nil
}
