package main

import (
	"fmt"
	"net/http"
	_ "net/url"
	"os"
	"path/filepath"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Define a custom response struct
type FileUploadResponse struct {
    Message     string `json:"message"`
    FileName    string `json:"filename"`
    FileSize    int64  `json:"filesize"`
    ContentType string `json:"content_type"`
    Folder      string `json:"folder"`
}

func uploadHandler(c *gin.Context) {
    fmt.Println("uploadHandler called")

    // Retrieve the string parameter
    selectedFolder := c.PostForm("folder")
    fmt.Println("Selected folder:", selectedFolder)

    if selectedFolder == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No folder selected"})
        return
    }

    // Define the absolute path to the "uploads" directory
    uploadsBasePath := "./"

	    // Define the folder path
		folderPath := filepath.Join(uploadsBasePath, selectedFolder)

    // Create the "uploads" folder if it doesn't exist
    if err := os.MkdirAll(uploadsBasePath, 0755); err != nil {
        fmt.Println("Error creating subdirectory:", err)
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create subdirectory"})
        return
    }

    // Retrieve the multipart form
    form, err := c.MultipartForm()

    if err != nil {
        c.JSON(http.StatusBadRequest, nil)
        return
    }

    files := form.File["file"]
    selectedFolder = form.Value["folder"][0] // Retrieve the selected folder
    fmt.Println("formFfile", files)
    fmt.Println(c.Request.MultipartForm.File)
    var uploadedFiles []string // To store the uploaded file names
    fmt.Println("Entering loop")

    // Check if no files were uploaded, and create an empty PDF if necessary
    if len(files) == 0 {
        emptyPDFPath := filepath.Join(folderPath, selectedFolder+".pdf")
        emptyPDF, err := os.Create(emptyPDFPath)
        if err != nil {
            fmt.Println("Error creating empty PDF:", err)
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create empty PDF"})
            return
        }
        emptyPDF.Close()
        uploadedFiles = append(uploadedFiles, selectedFolder+".pdf")
    }

    for _, file := range files {
        // Save the uploaded file to the specified path
        dst := filepath.Join(folderPath, file.Filename)
        if err := c.SaveUploadedFile(file, dst); err != nil {
            fmt.Println("Failed to save file:", err)
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save file"})
            return
        }
        // fmt.Println("File uploaded:", file)
        fmt.Println("uploads aloop", uploadedFiles)

        uploadedFiles = append(uploadedFiles, file.Filename) // Store the file name
    }
	submissionTime := time.Now()
    c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully", 
		"files": uploadedFiles,
		"submissionTime": submissionTime.Format(time.RFC3339),
	})
}

func downloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join("uploads", filename)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "File not found"})
		return
	}

	c.File(filepath)
}

// func fileListHandler(c *gin.Context) {
//     uploadsBasePath := "./"
    
//     files, err := os.ReadDir(uploadsBasePath)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve file list"})
//         return
//     }
    
//     var fileNames []string
//     for _, file := range files {
//         if !file.IsDir() {
//             fileNames = append(fileNames, file.Name())
//         }
//     }
    
//     c.JSON(http.StatusOK, gin.H{"files": fileNames})
// }
func fileListHandler(c *gin.Context) {
	// Retrieve the string parameter
	selectedFolder := c.PostForm("folder")
	fmt.Println("Selected folder:", selectedFolder)

	if selectedFolder == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No folder selected"})
		return
	}

    targetFolder := "your-target-folder" // Replace with the folder name you want to list

    targetFolderPath := filepath.Join("./", targetFolder)

    files, err := os.ReadDir(targetFolderPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve file list"})
        return
    }

    var fileNames []string
    for _, file := range files {
        if !file.IsDir() {
            fileNames = append(fileNames, file.Name())
        }
    }

    c.JSON(http.StatusOK, gin.H{"files": fileNames})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	// r.ForwardedByClientIP = true
	// r.SetTrustedProxies([]string{"127.0.0.1"})
	r.MaxMultipartMemory = 28 << 20 //8mib
	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	// Upload endpoint
	r.POST("/upload", uploadHandler)

	// Upload endpoint
	r.GET("/files", fileListHandler)

	// Download endpoint
	r.GET("/download/:filename", downloadHandler)

	r.Run(":8080")
}
