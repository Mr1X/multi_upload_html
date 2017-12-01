package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})

	router.POST("/upload", func(c *gin.Context) {
		// name := c.PostForm("name")
		// fmt.Println(name)
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		filename := header.Filename

		fmt.Println(file, err, filename)

		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusCreated, "upload successful")
	})

	router.POST("/multi/upload", func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(200000)
		if err != nil {
			log.Fatal(err)
		}

		formdata := c.Request.MultipartForm

		files := formdata.File["upload"]
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			out, err := os.Create(files[i].Filename)

			defer out.Close()

			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(out, file)

			if err != nil {
				log.Fatal(err)
			}

			c.String(http.StatusCreated, "upload successful \n")
		}
	})

	router.Run(":8000")
}
