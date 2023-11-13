package services

import (
	"fmt"
	"log"
	"time"

	// "time"

	"routepayfs.com/upload/config"
	"routepayfs.com/upload/db"
	"routepayfs.com/upload/errors"
	"routepayfs.com/upload/models"
	// "github.com/go-co-op/gocron"
)

//go:generate mockgen -destination=../mocks/medication_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationService

type DocumentService interface {
	CreateDocument(request *models.DocumentRequest, foldername string, filelocation string) (*models.DocumentResponse, *errors.Error)
	GetAllDocuments() ([]models.DocumentResponse, *errors.Error)
	GetDocumentByFolderName(folderName string) ([]models.DocumentResponse, *errors.Error)
	DeleteDocument(id string) error

	// DeleteUserDocument(userID uint) ([]models.DocumentResponse, *errors.Error)
	// UpdateDocument(request *models.UpdateDocumentRequest, documentID uint, userID uint) *errors.Error
	// FindDocument(documentName string, userId int) (*[]models.Document, error)
}

// medicationService struct
type documentService struct {
	Config        *config.Config
	UploadDocRepo db.DocumentRepository
}

// NewMedicationService instantiate an authService
func NewDocumentService(UploadDocRepo db.DocumentRepository, conf *config.Config) DocumentService {
	return &documentService{
		Config:        conf,
		UploadDocRepo: UploadDocRepo,
	}
}


// CreateDocument implements DocumentService.
func (d *documentService) CreateDocument(request *models.DocumentRequest, foldername string, filepath string) (*models.DocumentResponse, *errors.Error) {
	    fmt.Println("uploadHandlerService called")
	
		document := request.ReqToDocumentModel()
		response, err := d.UploadDocRepo.CreateDocument(document)   
		if err != nil {
			return nil, errors.ErrInternalServerError
		}
		return response.DocumentToResponse(), nil
}

func formatDateAgo(createdAtStr string) string {
	// Define the layout for the timestamp format
	layout := "2006-01-02 15:04:05 -0700 MST"

	// Parse the created_at string into a time.Time object
	createdAt, err := time.Parse(layout, createdAtStr)
	if err != nil {
		// Handle the error if the provided format is invalid
		return "Invalid date format"
	}

	// Get the current date
	currentDate := time.Now()

	// Calculate the difference in days
	difference := currentDate.Sub(createdAt)
	differenceInDays := int(difference.Hours() / 24)

	// Create the "X days ago" format
	if differenceInDays == 1 {
		return "1 day ago"
	} else {
		return fmt.Sprintf("%d days ago", differenceInDays)
	}
}


func (m *documentService) GetAllDocuments() ([]models.DocumentResponse, *errors.Error) {
	var documentResponses []models.DocumentResponse

	documents, err := m.UploadDocRepo.GetAllDocuments()
	if err != nil {
		log.Printf("error getting all documents: %v", err)
		return nil, errors.ErrInternalServerError
	}

	for _, document := range documents {
		createdAtString := time.Unix(document.CreatedAt, 0).Format("2006-01-02 15:04:05 -0700 MST")
		formattedDate := formatDateAgo(createdAtString)  // Format the date
		documentResponse := *document.DocumentToResponse()
		documentResponse.CreatedAt = formattedDate  
		documentResponses = append(documentResponses, documentResponse)
	}

	return documentResponses, nil
}



func (m *documentService) GetDocumentByFolderName(folderName string) ([]models.DocumentResponse, *errors.Error) {
	var documentResponses []models.DocumentResponse

	documents, err := m.UploadDocRepo.GetDocumentByFolderName(folderName)
	if err != nil {
		log.Printf("error getting documents by folder name: %v", err)
		return nil, errors.ErrInternalServerError
	}

	for _, document := range documents {
		createdAtString := time.Unix(document.CreatedAt, 0).Format("2006-01-02 15:04:05 -0700 MST")
		formattedDate := formatDateAgo(createdAtString)  // Format the date
		folder := models.Folder{Foldername: folderName}
		documentResponse := *document.DocumentToResponse()
		documentResponse.CreatedAt = formattedDate 
		documentResponse.Folder = folder 
		documentResponses = append(documentResponses, documentResponse)
	}

	return documentResponses, nil
}

func (s *documentService) DeleteDocument(id string) error{
	err := s.UploadDocRepo.DeleteDocument(id)
	if err != nil {
		return err
	}
	return nil
}