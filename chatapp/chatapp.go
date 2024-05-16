package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Channel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID        int    `json:"id"`
	ChannelID int    `json:"channel_id"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
}

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("error getting working directory")
	}

	db, err := sql.Open("sqlite", wd+"database.db")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("error opening database")
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error closing database")
		}
	}(db)

	router := gin.Default()

	//healthcheck
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ready",
		})
	})

	// create endpoints
	router.POST("/users", func(c *gin.Context) {})
	router.POST("/channels", func(c *gin.Context) {})
	router.POST("/messages", func(c *gin.Context) {})

	// list endpoints

	router.GET("/channels", func(c *gin.Context) {})
	router.GET("/messages", func(c *gin.Context) {})

	// login endpoints
	router.POST("/login", func(c *gin.Context) {})

	port := "8080"

	log.Info().
		Msgf("starting webserver on %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal().
			Err(err).
			Msg("error starting server")
	}
}

func createUser(c *gin.Context, db *sql.DB) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func login(c *gin.Context, db *sql.DB) {

	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE username = ?, user.Username AND password = ?", user.Username, user.Password)

	var id int

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
