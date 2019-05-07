package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type postView struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	PreBody string `json:"pre_body"`
	Body    string `json:"body"`
}

type postModel struct {
	gorm.Model
	Title   string `json:"title"`
	PreBody string `json:"pre_body"`
	Body    string `json:"body"`
}

type TestView struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Class     string `json:"class"`
	Points    uint   `json:"points"`
	Result    uint   `json:"result"`
}

type TestModel struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Class     string `json:"class"`
	Points    uint   `json:"points"`
	Result    uint   `json:"result"`
}

var db *gorm.DB

var dbUrl = "host=ec2-54-221-236-144.compute-1.amazonaws.com port=5432 user=jdriqytivsymsx dbname=d21u0n9iqblmf4 password=bdaedd4403c337f27fe794073ddc7c1650f3f841a67570a12cca5f0e1d72fbbe"
var dbUrlDev = "host=localhost port=5432 user=admin dbname=test password=admin sslmode=disable"

func initMigration() {
	var err error
	db, err = gorm.Open("postgres", dbUrl) //sslmode=disable

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	db.AutoMigrate(&postModel{})
	db.AutoMigrate(&TestModel{})
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	initMigration()

	r := gin.Default()

	r.LoadHTMLFiles("./info.html")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET, PUT, POST, DELETE"},
		AllowHeaders:     []string{"Origin, Authorization, Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: true,
		MaxAge:           100 * time.Hour,
	}))

	api := r.Group("/api")
	{
		api.GET("/", homeApi)
		api.POST("/post", createPost)
		api.DELETE("/delete", deletePost)
		api.GET("/delete/posts", deletePosts)
		api.PUT("/edit/post", editPost)
		api.GET("/post", getPost)
		api.GET("/posts", allPosts)
		api.GET("/count", countPosts)
	}

	r.GET("api/test/", homeApi)
	r.POST("api/test/", postTests)
	r.GET("api/test/users", getTestUsers)

	_ = r.Run(":" + port)
}

func homeApi(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"status":  http.StatusOK,
	// 		"message": "This is main page",
	// 	})

	c.HTML(http.StatusOK, "info.html", gin.H{
		"title": "test",
	})
}

func createPost(c *gin.Context) {
	post := postModel{
		Title:   c.PostForm("title"),
		PreBody: c.PostForm("pre_body"),
		Body:    c.PostForm("body"),
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
	preBody := c.PostForm("pre_body")
	body := c.PostForm("body")

	postId := c.PostForm("id")

	db.First(&post, postId)

	db.Model(&post).Update(postModel{
		Title:   title,
		PreBody: preBody,
		Body:    body,
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

	_post := postView{
		Id:      post.ID,
		Title:   post.Title,
		PreBody: post.PreBody,
		Body:    post.Body,
	}

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
			Id:      item.ID,
			Title:   item.Title,
			PreBody: item.PreBody,
			Body:    item.Body,
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

func postTests(c *gin.Context) {
	var points uint
	var result uint

	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	class := c.PostForm("class")

	ans1 := c.PostForm("ans1")
	ans2 := c.PostForm("ans2")
	ans3 := c.PostForm("ans3")

	if ans1 == "2" {
		points += 1
	}
	if ans2 == "3" {
		points += 1
	}
	if ans3 == "2" {
		points += 1
	}

	if points < 1 {
		result = 2
	}
	if points == 1 {
		result = 3
	}
	if points == 2 {
		result = 4
	}
	if points == 3 {
		result = 5
	}

	test := TestModel{
		FirstName: firstName,
		LastName:  lastName,
		Class:     class,
		Points:    points,
		Result:    result,
	}

	db.Create(&test)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"points": points,
		"result": result,
	})
}

func getTestUsers(c *gin.Context) {
	var tests []TestModel
	var _tests []TestView

	db.Find(&tests)

	if len(tests) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No posts found!",
		})

		return
	}

	for _, item := range tests {
		_tests = append(_tests, TestView{
			Id:        item.ID,
			FirstName: item.FirstName,
			LastName:  item.LastName,
			Class:     item.Class,
			Points:    item.Points,
			Result:    item.Result,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   _tests,
	})
}
