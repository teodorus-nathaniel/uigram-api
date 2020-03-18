package posts

type Post struct {
	ID          string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID      string    `json:"userId" bson:"userId"`
	Likes       *[]string `json:"likes" bson:"likes"`
	Dislikes    *[]string `json:"dislikes" bson:"dislikes"`
	Images      []string  `json:"images" bson:"images"`
	Link        string    `json:"link,omitempty" bson:"link,omitempty"`
	Description string    `json:"description" bson:"description"`
	Timestamp   string    `json:"timestamp" bson:"timestamp"`
}
