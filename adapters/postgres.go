package adapters

import (
	"context"
	"errors"
	"os"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/lithammer/shortuuid/v3"
)

type Database struct {
	Conn *pgx.Conn
}

type ImageData struct {
	Id   string
	Name string
	Date string
}

// function initializes database connection
func InitDB() (*Database, error) {
	// get database url
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, errors.New("$DATABASE_URL is not set")
	}

	// establish connection to database
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &Database{conn}, nil
}

// function closes database connection
func (db *Database) CloseDB() {
	db.Conn.Close(context.Background())
}

// function gets data for image specified by $id
func (db *Database) GetImageData(c *gin.Context, id string) (*ImageData, error) {
	var name, date string

	// fetch data by id
	err := db.Conn.QueryRow(c, "select name, date from images where id=$1", id).Scan(&name, &date)
	if err != nil {
		return nil, err
	}

	return &ImageData{id, name, date}, nil
}

// function gets data for all images
func (db *Database) GetAllImageData(c *gin.Context) ([]ImageData, error) {
	var images []ImageData
	
	rows, err := db.Conn.Query(c, "select * from images")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, name, date string

		err := rows.Scan(&id, &name, &date)
		if err != nil {
			return nil, err
		}

		images = append(images, ImageData{id, name, date})
	}

	return images, nil
}

// function inserts image data specified by $name into database
// returns created unique id for image
func (db *Database) InsertNewImageData(c *gin.Context, name string) (string, error) {
	// get id and date
	id := shortuuid.New()
	date := time.Now().Format(time.UnixDate)

	// insert data into database
	_, err := db.Conn.Exec(c, "insert into images (id, name, date) values ($1, $2, $3)", id, name, date)
	if err != nil {
		return "", err
	}

	return id, nil
}

// function deletes image data specified by $id
func (db *Database) DeleteImageData(c *gin.Context, id string) error {
	// delete image from database
	_, err := db.Conn.Exec(c, "delete from images where id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
