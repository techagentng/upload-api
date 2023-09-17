package services

import (
	"fmt"
	// "time"

	"routepayfs.com/upload/config"
	"routepayfs.com/upload/db"
	"routepayfs.com/upload/errors"
	"routepayfs.com/upload/models"
	// "github.com/go-co-op/gocron"
)

//go:generate mockgen -destination=../mocks/medication_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationService

type DocumentService interface {
	CreateDocument(request *models.DocumentRequest) (*models.DocumentResponse, *errors.Error)
	// GetAllUserDocuments(userID uint) ([]models.DocumentResponse, *errors.Error)
	// GetAllDocuments() ([]models.DocumentResponse, *errors.Error)
	// DeleteUserDocument(userID uint) ([]models.DocumentResponse, *errors.Error)
	// UpdateDocument(request *models.UpdateDocumentRequest, documentID uint, userID uint) *errors.Error
	// FindDocument(documentName string, userId int) (*[]models.Document, error)
}

// NewMedicationService instantiate an authService
func NewDocumentService(UploadDocRepo db.DocumentRepository, conf *config.Config) DocumentService {
	return &documentService{
		Config:        conf,
		UploadDocRepo: UploadDocRepo,
	}
}

// medicationService struct
type documentService struct {
	Config        *config.Config
	UploadDocRepo db.DocumentRepository
}

// CreateDocument implements DocumentService.
func (d *documentService) CreateDocument(request *models.DocumentRequest) (*models.DocumentResponse, *errors.Error) {
	    fmt.Println("uploadHandlerService called")
	
		document := request.ReqToDocumentModel()
		
		response, err := d.UploadDocRepo.CreateDocument(document)   
		if err != nil {
			return nil, errors.ErrInternalServerError
		}
		return response.DocumentToResponse(), nil
}

// func uploadHandler(c *gin.Context) {
//     fmt.Println("uploadHandler called")

//     // Retrieve the string parameter
// 	var documentRequest models.DocumentRequest
// 	_, user, err := GetValuesFromContext(c)
// 	if err != nil {
// 		err.Respond(c)
// 		return
// 	}
// 	userId := user.ID

// 	if err := decode(c, &documentRequest); err != nil {
// 		response.JSON(c, "", http.StatusBadRequest, nil, err)
// 		return
// 	}
// 	documentRequest.UserID = userId
// 	createdDocument, err := s.DocumentService.CreateDocument(&documentRequest)
// 	if err != nil {
// 		err.Respond(c)
// 		return
// 	}

// 	response.JSON(c, "document created successful", http.StatusCreated, createdDocument, nil)

//     selectedFolder := c.PostForm("folder")  //h
//     fmt.Println("Selected folder:", selectedFolder)

//     if selectedFolder == "" {
//         c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No folder selected"})
//         return
//     }

//     // Define the absolute path to the "uploads" directory //h
//     uploadsBasePath := "./uploads"

// 	    // Define the folder path
// 		folderPath := filepath.Join(uploadsBasePath, selectedFolder) //h

//     // Create the "uploads" folder if it doesn't exist
//     if err := os.MkdirAll(uploadsBasePath, 0755); err != nil {
//         fmt.Println("Error creating subdirectory:", err)
//         c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create subdirectory"}) //h
//         return
//     }

//     // Retrieve the multipart form
//     form, err := c.MultipartForm() //h

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


// // DeleteUserDocument implements DocumentService.
// func (d *documentService) DeleteUserDocument(userID uint) ([]models.DocumentResponse, *errors.Error) {
// 	panic("unimplemented")
// }

// // FindDocument implements DocumentService.
// func (d *documentService) FindDocument(documentName string, userId int) (*[]models.Document, error) {
// 	panic("unimplemented")
// }

// // GetAllDocuments implements DocumentService.
// func (d *documentService) GetAllDocuments() ([]models.DocumentResponse, *errors.Error) {
// 	panic("unimplemented")
// }

// // GetAllUserDocuments implements DocumentService.
// func (d *documentService) GetAllUserDocuments(userID uint) ([]models.DocumentResponse, *errors.Error) {
// 	panic("unimplemented")
// }



