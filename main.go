package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "b7rs1dkcgh9pgztf2hiw-postgresql.services.clever-cloud.com"
	port     = "50013"
	user     = "unoidcfwvryb8dxgh9o5"
	password = "Gc2dxLAx3zuurM3ayPYGFcBIHaS7bp"
	dbname   = "b7rs1dkcgh9pgztf2hiw"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func main() {
	// Configurar conexão com o banco de dados
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result []album
	for rows.Next() {
		var a album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			log.Fatal(err)
		}
		result = append(result, a)
	}

	c.IndentedJSON(http.StatusOK, result)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	_, err := db.Exec("INSERT INTO albums (id, title, artist, price) VALUES ($1, $2, $3, $4)", newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a album
	err := db.QueryRow("SELECT id, title, artist, price FROM albums WHERE id = $1", id).Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
			return
		}
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, a)
}
