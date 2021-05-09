package handlers

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jonathanlucki/shopify-challenge/adapters"
)

type UploadForm struct {
	Name string                `form:"name" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type GetImageResponse struct {
	Url  string
	Id   string
	Name string
	Date string
}

type GetImagesResponse struct {
	Images []GetImageResponse
}

type UploadImageResponse struct {
	Url string
	Id string
}

type DeleteImageResponse struct {
	Id string
}

func GetImages() gin.HandlerFunc {
	return func(c *gin.Context) {
		var imagesResponse GetImagesResponse = GetImagesResponse{}

		// get database connection
		db, ok := c.Keys["Database"].(*adapters.Database)
		if !ok {
			log.Printf("Error getting images: No database connection")
        	c.String(http.StatusBadRequest, "Error deleting image: No database connection")
			return
		}

		// get image data
		images, err := db.GetAllImageData(c)
		if err != nil {
			log.Printf("Error getting images: %s", err.Error())
			c.String(http.StatusBadRequest, "Error getting images: %s", err.Error())
			return
		}

		// iterate through images and construct images Response
		for _, image := range images {
			imageUrl, err := adapters.GetImageUrl(image.Id)
			if err != nil {
				log.Printf("Error getting image url: %s", err.Error())
				c.String(http.StatusBadRequest, "Error getting image url: %s", err.Error())
				return
			}

			imageResponse := GetImageResponse{imageUrl, image.Id, image.Name, image.Date}
			imagesResponse.Images = append(imagesResponse.Images, imageResponse)
		}

		// send get images response
		log.Printf("Images data retrieved")
		c.JSON(http.StatusOK, imagesResponse)
	}
}

// helper function to verify file is png
// this is not a very secure check but will do for our simple purposes
func checkFile(file *multipart.FileHeader) error {
	if file.Header.Get("content-type") != "image/png" {
		return errors.New("file is not a png file")
	}

	return nil
}

// handler for image uploading (via multipart form)
func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uploadForm UploadForm

		// get database connection
		db, ok := c.Keys["Database"].(*adapters.Database)
		if !ok {
			log.Printf("Error uploading image: No database connection")
        	c.String(http.StatusBadRequest, "Error deleting image: No database connection")
			return
		}

		// get storage connection
		storage, ok := c.Keys["Storage"].(*adapters.Storage)
		if !ok {
			log.Printf("Error uploading image: No storage connection")
        	c.String(http.StatusBadRequest, "Error deleting image: No storage connection")
			return
		}

		// bind upload form
		if err := c.ShouldBind(&uploadForm); err != nil {
			log.Printf("Error uploading image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error uploading image: %s", err.Error())
			return
		}

		// check file type
		err := checkFile(uploadForm.File)
		if err != nil {
			log.Printf("Error uploading image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error uploading image: %s", err.Error())
			return
		}

		// insert image data
		id, err := db.InsertNewImageData(c, uploadForm.Name)
		if err != nil {
			log.Printf("Error uploading image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error uploading image: %s", err.Error())
			return
		}

		// upload image to storage
		err = storage.UploadImage(c, id, uploadForm.File)
		if err != nil {
			log.Printf("Error uploading image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error uploading image: %s", err.Error())
			return
		}

		// get image url
		url, err := adapters.GetImageUrl(id)
		if err != nil {
			log.Printf("Error getting image url: %s", err.Error())
			c.String(http.StatusBadRequest, "Error getting image url: %s", err.Error())
			return
		}

		// send image upload response
		response := UploadImageResponse{url, id}
		log.Printf("Image uploaded - ID:%s", id)
		c.JSON(http.StatusOK, response)
	}
}

// handler for get image
func GetImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get id paramater
		id := c.Param("id")

		// get database connection
		db, ok := c.Keys["Database"].(*adapters.Database)
		if !ok {
			log.Printf("Error getting image: No database connection")
        	c.String(http.StatusBadRequest, "Error deleting image: No database connection")
			return
		}

		// get image data
		imageData, err := db.GetImageData(c, id)
		if err != nil {
			log.Printf("Error getting image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error getting image: %s", err.Error())
			return
		}

		// get image url
		url, err := adapters.GetImageUrl(imageData.Id)
		if err != nil {
			log.Printf("Error getting image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error getting image: %s", err.Error())
			return
		}

		// send get image response
		response :=  GetImageResponse{url, id, imageData.Name, imageData.Date}
		log.Printf("Image data retrieved - ID: %s", id)
		c.JSON(http.StatusOK, response)
	}
}

// handler for image deletion
func DeleteImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get id parameter
		id := c.Param("id")

		// get database connection
		db, ok := c.Keys["Database"].(*adapters.Database)
		if !ok {
			log.Printf("Error deleting image: No database connection")
        	c.String(http.StatusBadRequest, "Error deleting image: No database connection")
			return
		}

		// get storage connection
		storage, ok := c.Keys["Storage"].(*adapters.Storage)
		if !ok {
			log.Printf("Error deleting image: No storage connection")
        	c.String(http.StatusBadRequest, "Error deleting image: No storage connection")
			return
		}

		// delete image from storage
		err := storage.DeleteImage(c, id)
		if err != nil {
			log.Printf("Error deleting image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error deleting image: %s", err.Error())
			return
		}

		// delete image data
		err = db.DeleteImageData(c, id)
		if err != nil {
			log.Printf("Error deleting image: %s", err.Error())
			c.String(http.StatusBadRequest, "Error deleting image: %s", err.Error())
			return
		}

		// send deletion response
		response :=  DeleteImageResponse{id}
		log.Printf("Image deleted - ID:%s", id)
		c.JSON(http.StatusOK, response)
	}
}
