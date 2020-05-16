package posts

import (
	"errors"

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

type Post struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID        string             `json:"-" bson:"userId"`
	Title         string             `json:"title" bson:"title"`
	Owner         *Owner             `json:"owner" bson:"-"`
	Likes         []string           `json:"-" bson:"likes"`
	Dislikes      []string           `json:"-" bson:"dislikes"`
	Description   *string            `json:"description,omitempty" bson:"description"`
	Link          *string            `json:"link,omitempty" bson:"link,omitempty"`
	Images        []string           `json:"images" bson:"images"`
	Timestamp     int64              `json:"timestamp" bson:"timestamp"`
	LikeCount     int                `json:"likeCount" bson:"-"`
	DislikeCount  int                `json:"dislikeCount" bson:"-"`
	CommentsCount int                `json:"commentsCount" bson:"-"`
	Liked         bool               `json:"liked,omitempty" bson:"-"`
	Disliked      bool               `json:"disliked,omitempty" bson:"-"`
	Saved         bool               `json:"saved,omitempty" bson:"-"`
}

func (post *Post) fillEmptyValues() {
	if post.Likes == nil {
		post.Likes = []string{}
	}
	if post.Dislikes == nil {
		post.Dislikes = []string{}
	}
}

func (post *Post) validateData() error {
	if post.Title == "" {
		return errors.New("post title can't be empty")
	}
	if len(post.Images) == 0 {
		return errors.New("post images must be at least one")
	}

	return nil
}

func (post *Post) processData(user *users.User) {
	post.LikeCount = len(post.Likes)
	post.DislikeCount = len(post.Dislikes)

	res, _ := users.GetUserById(post.UserID)
	post.Owner = &Owner{ID: res.ID.Hex(), ProfilePic: res.ProfilePic, Username: res.Username, Followed: false}
	if user != nil {
		post.Owner.Followed = utils.SearchArray(user.Following, post.Owner.ID)
		post.Saved = utils.SearchArray(user.SavedPosts, post.ID.Hex())
		post.Liked = utils.SearchArray(post.Likes, user.ID.Hex())
		post.Disliked = utils.SearchArray(post.Dislikes, user.ID.Hex())
	}

}

func (post *Post) deriveToPost(user *users.User) {
	post.processData(user)

	post.Link = nil
	post.Description = nil
	post.Likes = nil
	post.Dislikes = nil
}

func (post *Post) deriveToPostDetail(user *users.User) {
	post.processData(user)
}

// id: string;
// title: string;
// owner: UserBasicInfo;
// images: string[];
// likeCount: number;
// dislikeCount: number;
// commentsCount: number;
// liked?: boolean;
// disliked?: boolean;
// timestamp: Date;
// saved: boolean;

// likes: string[];
// dislikes: string[];
// description: string;
// link: string;
