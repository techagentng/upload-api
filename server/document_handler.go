package server

import (
	// "encoding/json"
	_ "encoding/json"
	"fmt"
	"log"
	_ "log"
	"math/rand"
	_ "math/rand"
	"mime/multipart"

	// "mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/google/uuid"
	_ "routepayfs.com/upload/errors"
	"routepayfs.com/upload/models"
	"routepayfs.com/upload/server/response"
)
type FileWithTime struct {
    FileName       string    `json:"fileName"`
    DateCreated   time.Time `json:"dateCreated"`
    ModificationTime time.Time `json:"modificationTime"`
}
type ErrorResponse struct {
    Error string `json:"error"`
}

type Folder struct {
	ID   uint
	Name string
}


var folderMap = map[string]uint{
	"NDPR": 1,
	"PCIDSS": 2,
    "ISO-27001": 3,
    "Compliance": 4,
    "Organogram": 5,
    "SLA": 6,
    "TAT": 7,
    "culture": 8,
    "Information security management": 9,
	// Add more folder mappings as needed
}
// func getFolderID(folderName string) (uint, error) {
// 	folderID, ok := folderMap[folderName]
// 	if !ok {
// 		return 0, fmt.Errorf("folder not found")
// 	}
// 	return folderID, nil
// }
func saveUploadedFile(c *gin.Context, fileHeader *multipart.FileHeader, folderPath string) error {
    dst := filepath.Join(folderPath, fileHeader.Filename)

    _, err := os.Stat(dst)
    if err == nil {
        fmt.Println("File already exists:", fileHeader.Filename)
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "File already exists"})
        return err
    }

    if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
        fmt.Println("Failed to save file:", err)
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save file"})
        return err
    }

    return nil
}
func (s *Server) handleFileUpload() gin.HandlerFunc {
    return func(c *gin.Context) {
        var uploadRequest models.DocumentRequest
    		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		uploadRequestUserId := user.ID
        UploaderName := user.Name
        
        uploadRequest.Filename = c.PostForm("filename")
        uploadRequest.DocumentType = c.PostForm("doctype")
        uploadRequest.Folder= c.PostForm("folder")
        uploadRequest.DocumentNumber= concatenateString(c.PostForm("doctype"))
        uploadRequest.Department = c.PostForm("department")
        uploadRequest.Division = c.PostForm("division")
        uploadRequest.Docclass = c.PostForm("docclass")
        uploadRequest.UserID = uploadRequestUserId
        uploadRequest.UploaderName = UploaderName
    
           // Validate the request datavalidate
           v := validator.New()
           if err := v.Struct(uploadRequest); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        selectedFolder := uploadRequest.Folder

        if selectedFolder == "" {
            c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No folder selected"})
            return
        }

        folderPath := filepath.Join("./uploads", selectedFolder)

        if err := os.MkdirAll(folderPath, 0755); err != nil {
            fmt.Println("Error creating subdirectory:", err)
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create subdirectory"})
            return
        }

      files := c.Request.MultipartForm.File["file"]
           fmt.Println("formFfile", files)


           for _, fileHeader := range files {
            // Open the file for reading
            file, err := fileHeader.Open()
            if err != nil {
                fmt.Println("Error opening uploaded file:", err)
                c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to open uploaded file"})
                return
            }
            defer file.Close()
           
            // Use the original filename as the destination filename
            filePath := filepath.Join(folderPath, fileHeader.Filename)
        
            // Save the uploaded file to the specified folder
            if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
                fmt.Println("Error saving file:", err)
                c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save uploaded file"})
                return
            }
        }
           
            createDocument, err := s.DocumentService.CreateDocument(&uploadRequest)
            if err != nil {
                fmt.Println("Error creating document:", err)
                c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to parse service"})
                return
            }

        response.JSON(c, "document created successfully", http.StatusCreated, createDocument, nil)
    }
}

func generateRandomNumbers(length int) string {
    rand.Seed(time.Now().UnixNano())
    result := ""
    for i := 0; i < length; i++ {
        result += fmt.Sprintf("%d", rand.Intn(10)) 
    }
    return result
}

func concatenateString(input string) string {
    if len(input) < 3 {
        return input
    }

    firstThree := input[:3]

    // Generate 4 random digits
    randomDigits := generateRandomNumbers(4)

    result := firstThree + "-" + randomDigits

    return result
}

func (s *Server) handleDownloadDocument() gin.HandlerFunc {
	return func(c *gin.Context){
        filename := c.Param("filename")
        filepath := filepath.Join("uploads", filename)
        _, err := os.Stat(filepath)
        if os.IsNotExist(err) {
            c.JSON(http.StatusNotFound, ErrorResponse{Error: "File not found"})
            return
        }
    
        c.File(filepath)
	}
}

func (s *Server) handleGetFolderList() gin.HandlerFunc {
	return func(c *gin.Context){
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
}

func (s *Server) handleFindDocument() gin.HandlerFunc {
	return func(c *gin.Context){

	}
}

func (s *Server) handleDeleteDocument() gin.HandlerFunc {
	return func(c *gin.Context){
        log.Println("Delete document called")
        fmt.Println("Delete document called")

    folderName := c.Param("folderName")
    fileName := c.Param("fileName")

            // Decode the encoded file name
            encodedFileName := c.Param("fileName")
            fileName, err := url.QueryUnescape(encodedFileName)
            if err != nil {
                c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid file name"})
                return
            }
            
    filePath := filepath.Join("./uploads", folderName, fileName)
    fmt.Println("Received delete request for folder:", folderName, "and filename:", fileName)

    err = os.Remove(filePath)

    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete file"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "File deleted successfullyxxx"})
	}
}

func (s *Server) handleGetAllDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		// _, _, err := GetValuesFromContext(c)
		// if err != nil {
		// 	err.Respond(c)
		// 	return
		// }
		allDocument, err := s.DocumentService.GetAllDocuments()
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "documents retrieved successfully", http.StatusOK, allDocument, nil)
	}
}
func (s *Server) handleEditDocument() gin.HandlerFunc {
	return func(c *gin.Context){

	}
}

