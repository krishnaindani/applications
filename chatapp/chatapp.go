package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog/log"
	"os"
)

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
