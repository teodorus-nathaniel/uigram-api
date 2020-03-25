package posts

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID string             `json:"-" bson:"userId"`
	Title  string             `json:"title" bson:"title"`
	// Owner `json:"owner"` pake user id di atas
	Likes         []string `json:"likes" bson:"likes"`
	Dislikes      []string `json:"dislikes" bson:"dislikes"`
	Description   *string  `json:"description,omitempty" bson:"description"`
	Link          *string  `json:"link,omitempty" bson:"link,omitempty"`
	Images        []string `json:"images" bson:"images"`
	Timestamp     string   `json:"timestamp" bson:"timestamp"`
	LikeCount     int      `json:"likeCount"`
	DislikeCount  int      `json:"dislikeCount"`
	CommentsCount int      `json:"commentsCount"`
	Liked         bool     `json:"liked,omitempty"`
	Disliked      bool     `json:"disliked,omitempty"`
	Saved         bool     `json:"saved,omitempty"`
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

func (post *Post) addDerivedData() {
	post.LikeCount = len(post.Likes)
	post.DislikeCount = len(post.Dislikes)
	// post.Liked = cari itu ama saved jg ama disliked
}

func (post *Post) deriveToPost() {
	post.addDerivedData()

	post.Link = nil
	post.Description = nil
	post.Likes = nil
	post.Dislikes = nil
}

func (post *Post) deriveToPostDetail() {
	post.addDerivedData()
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
