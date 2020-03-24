package posts

type Post struct {
	ID       string   `json:"id,omitempty" bson:"_id,omitempty"`
	UserID   string   `json:"-" bson:"userId"`
	Likes    []string `json:"-" bson:"likes"`
	Dislikes []string `json:"-" bson:"dislikes"`
	// Description string   `json:"-" bson:"description"` kalo dipisah ama post detail buang
	Title     string   `json:"title" bson:"title"`
	Images    []string `json:"images" bson:"images"`
	Link      string   `json:"link,omitempty" bson:"link,omitempty"`
	Timestamp string   `json:"timestamp" bson:"timestamp"`
	// Owner `json:"owner"` pake user id di atas
	LikeCount     int  `json:"likeCount"`
	DislikeCount  int  `json:"dislikeCount"`
	CommentsCount int  `json:"commentsCount"`
	Liked         bool `json:"liked,omitempty"`
	Disliked      bool `json:"disliked,omitempty"`
	Saved         bool `json:"saved,omitempty"`
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
