package models

type User struct {
	ID             string   `json:"_id" bson:"_id"`
	Username       string   `json:"username" bson:"username"`
	Password       string   `json:"password" bson:"password"`
	Email          string   `json:"email" bson:"email"`
	ProfilePicture string   `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"`
	Followers      []string `json:"followers" bson:"followers"`
	Following      []string `json:"following" bson:"following"`
	Saved          bool     `json:"saved"`
	Liked          bool     `json:"liked,omitempty"`
	Disliked       bool     `json:"disliked,omitempty"`
}
