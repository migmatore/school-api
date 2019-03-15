package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type postView struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
}

type postModel struct {
	gorm.Model
	Title string `json:"title"`
}

var db *gorm.DB

func initMigration() {
	var err error
	db, err = gorm.Open("postgres",
		"host=ec2-184-72-238-22.compute-1.amazonaws.com port=5432 user=tvblcbvgwsxrhs dbname=da0mec5vqlq9bk password=9ed13e8cf2b3b5a78f43b5a9925744747e8680c200c25dbb6a42dc2aba286f58") //sslmode=disable

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	db.AutoMigrate(&postModel{})
}

func main() {
	initMigration()

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET, PUT, POST, DELETE"},
		AllowHeaders:     []string{"Origin, Authorization, Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: true,
		MaxAge:           50 * time.Hour,
	}))

	api := r.Group("/api")
	{
		api.GET("/", home)
		api.POST("/post", createPost)
		api.GET("/delete", deletePost)
		api.GET("/delete/posts", deletePosts)
		api.POST("/edit/post", editPost)
		api.GET("/post", getPost)
		api.GET("/posts", allPosts)
		api.GET("/count", countPosts)
	}

	_ = r.Run("5000")
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "This is main page",
	})
}

func createPost(c *gin.Context) {
	post := postModel{
		Title: c.PostForm("title"),
	}

	db.Create(&post)

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "post created successfully!",
	})
}

func deletePost(c *gin.Context) {
	var post postModel

	postId := c.Query("id")

	db.Unscoped().Delete(&post, postId)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "post delete successfully!",
	})
}

func deletePosts(c *gin.Context) {
	var posts []postModel

	db.Unscoped().Delete(&posts)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "deleted all posts!",
	})
}

func editPost(c *gin.Context) {
	var post postModel

	title := c.PostForm("title")
	postId := c.PostForm("id")

	db.First(&post, postId)

	db.Model(&post).Update(postModel{
		Title: title,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "post " + string(post.ID) + "has been changed!",
	})
}

func getPost(c *gin.Context) {
	var post postModel

	postId := c.Query("id")

	db.First(&post, postId)

	if post.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No post found!",
		})

		return
	}

	_post := postView{Id: post.ID, Title: post.Title}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   _post,
	})
}

func allPosts(c *gin.Context) {
	var posts []postModel
	var _posts []postView

	db.Find(&posts)

	if len(posts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No posts found!",
		})

		return
	}

	for _, item := range posts {
		_posts = append(_posts, postView{
			Id:    item.ID,
			Title: item.Title,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   _posts,
	})
}

func countPosts(c *gin.Context) {
	var posts []postModel

	db.Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   len(posts),
	})
}
