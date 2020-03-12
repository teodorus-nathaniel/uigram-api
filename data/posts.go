package data

import (
	"time"

	"github.com/teodorus-nathaniel/ui-share-api/models"
)

var DummyPosts = [5]models.Post{
	{
		ID:          "123",
		Images:      []string{"https://www.bigstockphoto.com/images/homepage/module-6.jpg"},
		Description: "DSIFUHFUIDasdoisdoifjsdf",
		Link:        "https://teodorus-nathaniel.github.io",
		Timestamp:   time.Now().String(),
		UserID:      "123",
	},
	{
		ID:          "1234",
		Images:      []string{"https://www.bigstockphoto.com/images/homepage/module-6.jpg"},
		Description: "DSIFUHFUIDasdoisdoifjsdf",
		Link:        "https://teodorus-nathaniel.github.io",
		Timestamp:   time.Now().String(),
		UserID:      "123",
	},
	{
		ID:          "12",
		Images:      []string{"https://www.bigstockphoto.com/images/homepage/module-6.jpg"},
		Description: "DSIFUHFUIDasdoisdoifjsdf",
		Link:        "https://teodorus-nathaniel.github.io",
		Timestamp:   time.Now().String(),
		UserID:      "123",
	},
	{
		ID:          "1223",
		Images:      []string{"https://www.bigstockphoto.com/images/homepage/module-6.jpg"},
		Description: "DSIFUHFUIDasdoisdoifjsdf",
		Link:        "https://teodorus-nathaniel.github.io",
		Timestamp:   time.Now().String(),
		UserID:      "123",
	},
	{
		ID:          "12113",
		Images:      []string{"https://www.bigstockphoto.com/images/homepage/module-6.jpg"},
		Description: "DSIFUHFUIDasdoisdoifjsdf",
		Link:        "https://teodorus-nathaniel.github.io",
		Timestamp:   time.Now().String(),
		UserID:      "123",
	},
}
