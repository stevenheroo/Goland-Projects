package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var users = make(map[string]User)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Println("Error loading .env file")
	}

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", addAlbum)
	router.GET("/albums/:id", getAlbum)
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.GET("/login/verify", verifyLogin)
	router.GET("/users", getUsers)

	err := router.Run("localhost:8080")
	if err != nil {
		log.SetPrefix("Port failed :::: ")
		log.Fatal(err)
		return
	}

}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		log.SetPrefix("Post failed :::: ")
		log.Fatal(err)
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbum(c *gin.Context) {
	id := c.Param("id")
	for _, a := range albums {
		if id == a.ID {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "success",
				"data":    a,
			})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func signup(c *gin.Context) {

	var newUser User
	if err := c.BindJSON(&newUser); err != nil || len(newUser.Password) == 0 {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{
			"message": "Signup Failed",
			"error":   err,
		})
		return
	}

	hashPass, hashErr := hashPwd(newUser.Password)
	if hashErr != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{
			"message": "Failed",
			"error":   hashErr,
		})
		return
	}
	newUser.Password = hashPass
	users[newUser.Username] = newUser

	c.IndentedJSON(http.StatusExpectationFailed, gin.H{
		"message": "Success",
		"data":    newUser,
	})
}

func login(c *gin.Context) {

	var newUser User
	if err := c.BindJSON(&newUser); err != nil || len(newUser.Password) == 0 {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "Login Failed",
		})
		return
	}
	log.Print(newUser)
	if users[newUser.Username] == (User{}) {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "Login Failed",
			"token":   "",
		})
		return
	}

	token, err := auth(newUser)

	fmt.Println(err)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "Login Failed",
			"token":   "",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
}

func verifyLogin(c *gin.Context) {
	token, isPresent := getBearerToken(c)
	if !isPresent {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
		return
	}

	sub, err := validateToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Token",
			"subject": sub,
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "success",
		"subject": sub,
	})

}

func getUsers(c *gin.Context) {
	if len(users) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "success",
			"users":   map[string]User{},
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "success",
		"users":   users,
	})
}

//func userTable() map[string]User {
//	users["Steve"] = User{Username: "Steve", Password: "Test1234"}
//	users["Alex"] = User{Username: "Alex", Password: "Test1234"}
//
//	return users
//}
