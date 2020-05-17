package comments

import (
	"github.com/teodorus-nathaniel/uigram-api/users"
	"github.com/teodorus-nathaniel/uigram-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Owner struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	ProfilePic string `json:"profilePic,omitempty"`
	Followed   bool   `json:"followed"`
}

type Comment struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content      string             `json:"content" bson:"content"`
	PostID       string             `json:"postId" bson:"postId"`
	UserID       string             `json:"-" bson:"userId"`
	Likes        []string           `json:"likes" bson:"likes"`
	Dislikes     []string           `json:"dislikes" bson:"dislikes"`
	Parent       *string            `json:"-" bson:"parent,omitempty"`
	Timestamp    int64              `json:"timestamp" bson:"timestamp"`
	Owner        Owner              `json:"owner" bson:"-"`
	Liked        bool               `json:"liked" bson:"-"`
	Disliked     bool               `json:"disliked" bson:"-"`
	LikeCount    int                `json:"likeCount" bson:"-"`
	DislikeCount int                `json:"dislikeCount" bson:"-"`
	RepliesCount int64              `json:"repliesCount" bson:"-"`
}

func (comment *Comment) deriveData(user *users.User) {
	owner, err := users.GetUserById(comment.UserID)
	if err != nil {
		return
	}

	comment.Owner.ID = owner.ID.Hex()
	comment.Owner.ProfilePic = owner.ProfilePic
	comment.Owner.Username = owner.Username

	if comment.Likes == nil {
		comment.LikeCount = 0
	} else {
		comment.LikeCount = len(comment.Likes)
	}
	if comment.Dislikes == nil {
		comment.DislikeCount = 0
	} else {
		comment.DislikeCount = len(comment.Dislikes)
	}

	if user != nil {
		comment.Liked = utils.SearchArray(comment.Likes, user.ID.Hex())
		comment.Disliked = utils.SearchArray(comment.Dislikes, user.ID.Hex())
	}

	if comment.Parent == nil {
		comment.RepliesCount = getCommentsRepliesCount(comment.ID.Hex())
	}
}
