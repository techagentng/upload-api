package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	_ "net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"routepayfs.com/upload/config"
	"routepayfs.com/upload/db"
	"routepayfs.com/upload/server"
	"routepayfs.com/upload/services"
)
type FileWithTime struct {
    FileName       string    `json:"fileName"`
    DateCreated   time.Time `json:"dateCreated"`
    ModificationTime time.Time `json:"modificationTime"`
}
// Define a custom response struct
type FileUploadResponse struct {
    Message     string `json:"message"`
    FileName    string `json:"filename"`
    FileSize    int64  `json:"filesize"`
    ContentType string `json:"content_type"`
    Folder      string `json:"folder"`
}

// func uploadHandler(c *gin.Context) {
//     fmt.Println("uploadHandler called")

//     // Retrieve the string parameter
//     selectedFolder := c.PostForm("folder")
//     fmt.Println("Selected folder:", selectedFolder)

//     if selectedFolder == "" {
//         c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No folder selected"})
//         return
//     }

//     // Define the absolute path to the "uploads" directory
//     uploadsBasePath := "./uploads"

// 	    // Define the folder path
// 		folderPath := filepath.Join(uploadsBasePath, selectedFolder)

//     // Create the "uploads" folder if it doesn't exist
//     if err := os.MkdirAll(uploadsBasePath, 0755); err != nil {
//         fmt.Println("Error creating subdirectory:", err)
//         c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create subdirectory"})
//         return
//     }

//     // Retrieve the multipart form
//     form, err := c.MultipartForm()

//     if err != nil {
//         c.JSON(http.StatusBadRequest, nil)
//         return
//     }

//     files := form.File["file"]
//     selectedFolder = form.Value["folder"][0] // Retrieve the selected folder
//     fmt.Println("formFfile", files)
//     fmt.Println(c.Request.MultipartForm.File)
//     var uploadedFiles []string // To store the uploaded file names
//     fmt.Println("Entering loop")

//     // Check if no files were uploaded, and create an empty PDF if necessary
//     if len(files) == 0 {
//         emptyPDFPath := filepath.Join(folderPath, selectedFolder+".pdf")
//         emptyPDF, err := os.Create(emptyPDFPath)
//         if err != nil {
//             fmt.Println("Error creating empty PDF:", err)
//             c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create empty PDF"})
//             return
//         }
//         emptyPDF.Close()
//         uploadedFiles = append(uploadedFiles, selectedFolder+".pdf")
//     }

//     for _, file := range files {
//         // Save the uploaded file to the specified path
//         dst := filepath.Join(folderPath, file.Filename)
		
// 		_, err := os.Stat(dst)
//         if err == nil {
//             // File already exists, handle the error
//             fmt.Println("File already exists:", file.Filename)
//             c.JSON(http.StatusBadRequest, ErrorResponse{Error: "File already exists"})
//             return
//         }

//         if err := c.SaveUploadedFile(file, dst); err != nil {
//             fmt.Println("Failed to save file:", err)
//             c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save file"})
//             return
//         }
//         // fmt.Println("File uploaded:", file)
//         fmt.Println("uploads aloop", uploadedFiles)

//         uploadedFiles = append(uploadedFiles, file.Filename) // Store the file name
//     }
// 	submissionTime := time.Now()
//     c.JSON(http.StatusOK, gin.H{
// 		"message": "Files uploaded successfully", 
// 		"files": uploadedFiles,
// 		"submissionTime": submissionTime.Format(time.RFC3339),
// 	})
// }

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

func folderListHandler(c *gin.Context) {
    uploadsBasePath := "./uploads"
    selectedFolder := c.Param("folderName")
    folderPath := filepath.Join(uploadsBasePath, selectedFolder)

    fmt.Println("Reading folder:", folderPath)
    files, err := os.ReadDir(folderPath)
    if err != nil {
        fmt.Println("Error reading folder:", err)
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "No file or folder created"})
        return
    }

    var fileList []FileWithTime
    for _, file := range files {
        if !file.IsDir() {
            filePath := filepath.Join(folderPath, file.Name())

            fileInfo, err := os.Stat(filePath)
            if err != nil {
                fmt.Println("Error getting file info:", err)
                continue
            }

            dateCreated := fileInfo.ModTime()
            modificationTime := fileInfo.ModTime()

            fileData := FileWithTime{
                FileName:         file.Name(),
                DateCreated:     dateCreated,
                ModificationTime: modificationTime,
            }

            fileList = append(fileList, fileData)
        }
    }

    c.JSON(http.StatusOK, gin.H{"files": fileList})
}

func fileListHandler(c *gin.Context) {
    uploadsBasePath := "./uploads"
    
    files, err := os.ReadDir(uploadsBasePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve file list. it does not exist"})
        return
    }
    
    var fileList []FileWithTime
    for _, file := range files {
        if !file.IsDir() {
            filePath := filepath.Join(uploadsBasePath, file.Name())

            fileInfo, err := os.Stat(filePath)
            if err != nil {
                fmt.Println("Error getting file info:", err)
                continue
            }

            dateCreated := fileInfo.ModTime()
            modificationTime := fileInfo.ModTime()

            fileData := FileWithTime{
                FileName:         file.Name(),
                DateCreated:     dateCreated,
                ModificationTime: modificationTime,
            }

            fileList = append(fileList, fileData)
        }
    }
    
    c.JSON(http.StatusOK, gin.H{"files": fileList})
}

func deleteFileHandler(c *gin.Context) {
	fmt.Println("called")
    folderName := c.Param("folderName")
    fileName := c.Param("fileName")
    // Decode the folder name
	decodedFolderName, err := url.QueryUnescape(fileName)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
    filePath := filepath.Join("./uploads", folderName, decodedFolderName)
    
    err = os.Remove(filePath)

    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete file"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
    conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

    gormDB := db.GetDB(conf)
    authRepo := db.NewAuthRepo(gormDB)
    authService := services.NewAuthService(authRepo, conf)
    documentRepo := db.NewDocumentRepo(gormDB)
    documentService := services.NewDocumentService(documentRepo, conf)

    s := &server.Server{
		Config:                   conf,
		AuthRepository:           authRepo,
		AuthService:              authService,
        DocumentService:         documentService,
	}
	r := gin.Default()
	r.Use(cors.Default())
	// r.ForwardedByClientIP = true
	// r.SetTrustedProxies([]string{"127.0.0.1"})
	r.MaxMultipartMemory = 28 << 20 //8mib
	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	// r.Run(":8080")
    s.Start()
}
