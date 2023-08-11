package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/cors"
  _"fmt"
	_"net/http"
	"net/url"
	"os"
	"path/filepath"
	_"strings"
)
func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Bad request"})
		return
	}

	filename := file.Filename
	ext := filepath.Ext(filename)
	allowedExts := []string{".pdf", ".docx"}

	// Check if the file extension is allowed
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
		
	if !allowed {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "File type not allowed"})
		return
	}
		// Get the selected folder from the request
		selectedFolder := c.PostForm("folder")
		if selectedFolder == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No folder selected"})
			return
	}

	// Save the file in the project directory
	dst := filepath.Join("uploads", filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, FileUploadResponse{Message: "File uploaded successfully"})
}

func downloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join("uploads", filename)

	// Check if the file exists
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "File not found"})
		return
	}

	c.File(filepath)
}

type FileUploadResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
r := gin.Default()
r.Use(cors.Default())
r.ForwardedByClientIP = true
r.SetTrustedProxies([]string{"127.0.0.1"})

// Define the proxy URL
	proxyURL, _ := url.Parse("http://localhost:8080")

	// Create a custom transport with the proxy settings
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
	}

	// Create a route that makes a request using the custom HTTP client
	r.GET("/", func(c *gin.Context) {
		resp, err := client.Get("http://localhost:8080")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		// Process the response or return it as needed
		// ...

		c.String(http.StatusOK, "Request through proxy successful")
	})

	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	// Upload endpoint
	r.POST("/upload", uploadHandler)

	// Download endpoint
	r.GET("/download/:filename", downloadHandler)

	r.Run(":8080")
}