package posts

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/screenshot"
	"github.com/teodorus-nathaniel/uigram-api/users"
	"github.com/teodorus-nathaniel/uigram-api/utils"
)

func getUserFromMiddleware(c *gin.Context) *users.User {
	user, _ := c.Get("user")
	if user == nil {
		return nil
	}
	return user.(*users.User)
}

func getPostsHandler(c *gin.Context) {
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	posts, err := getPosts(sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getPostsFeedsHandler(c *gin.Context) {
	_, limit, page := utils.GetQueryStringsForPagination(c)

	data, _ := c.Get("user")
	user := data.(*users.User)

	followingIds := user.Following
	posts, err := getPostsByUserId(followingIds, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getPostHandler(c *gin.Context) {
	id := c.Param("id")
	post, err := getPost(id, getUserFromMiddleware(c))

	if err != nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail("post with id: "+id+" was not found"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"post": post}))
}

func postPostHandler(c *gin.Context) {
	user := getUserFromMiddleware(c)
	description, _ := c.GetPostForm("description")
	title, _ := c.GetPostForm("title")
	link, _ := c.GetPostForm("link")
	files := c.Request.MultipartForm.File["files"]
	images, _ := c.GetPostFormArray("images")

	var dataImages []string
	var filePaths []string
	timestamp := time.Now().Unix()
	filesLen := len(files)

	for _, image := range images {
		tokens := strings.Split(image, "--")
		if len(tokens) < 2 {
			dataImages = append(dataImages, image)
			continue
		}
		idx, err := strconv.Atoi(tokens[1])
		if err != nil || idx >= filesLen {
			dataImages = append(dataImages, image)
			continue
		}

		path := "img/post-" +
			strconv.Itoa(idx) +
			strconv.FormatInt(timestamp, 10) +
			"-" + user.ID.Hex() +
			filepath.Ext(files[idx].Filename)
		filePaths = append(filePaths, path)

		path = "http://localhost:8080/" + path
		dataImages = append(dataImages, path)
	}

	if len(filePaths) != len(files) {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail("Files not match"))
		return
	}

	for idx, file := range files {
		err := c.SaveUploadedFile(file, filePaths[idx])

		if err != nil {
			c.JSON(http.StatusInternalServerError, jsend.GetJSendFail("File upload fail"))
			return
		}
	}

	post, err := insertPost(dataImages, title, description, link, timestamp, user)
	post.deriveToPost(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"post": post}))
}

func getUserPostHandler(c *gin.Context) {
	id := c.Param("id")
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	posts, err := GetPostByOwner(id, sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail("Fail fetching user posts"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getUserSavedPostHandler(c *gin.Context) {
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	posts, err := getSavedPosts(sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail("Fail fetching user posts"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

type screenshotReqBody struct {
	URL string `json:"url"`
}

func postScreenshot(c *gin.Context) {
	var url screenshotReqBody
	json.NewDecoder(c.Request.Body).Decode(&url)

	matched, err := regexp.MatchString(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`, url.URL)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail("Invalid URL"))
		return
	}

	if !strings.HasPrefix(url.URL, "http://") && !strings.HasPrefix(url.URL, "https://") {
		url.URL = "http://" + url.URL
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, screenshot.FullScreenshot(url.URL, 50, &buf)); err != nil {
		log.Fatal(err)
	}
	cleanedURL := strings.ReplaceAll(url.URL, "/", "")
	cleanedURL = strings.ReplaceAll(cleanedURL, ".", "")
	cleanedURL = strings.ReplaceAll(cleanedURL, ":", "")
	filepath := "img/" + cleanedURL + strconv.Itoa(rand.Intn(1000)) + ".jpeg"
	fmt.Println(strings.ReplaceAll(url.URL, "/", ""))
	fmt.Println(filepath)
	if err := ioutil.WriteFile(filepath, buf, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	res := "http://localhost:8080/" + filepath

	c.JSON(http.StatusCreated, jsend.GetJSendSuccess(gin.H{"url": res}))
}

type patchLikesReqBody struct {
	Like    *bool `json:"like,omitempty"`
	Dislike *bool `json:"dislike,omitempty"`
}

func patchLikes(c *gin.Context) {
	user := getUserFromMiddleware(c)
	id, _ := c.Params.Get("id")
	var body patchLikesReqBody
	json.NewDecoder(c.Request.Body).Decode(&body)

	post, err := updatePostLikes(user, id, body.Like, body.Dislike)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"post": post}))
}
