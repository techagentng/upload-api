package services

import (
	"fmt"
	"log"
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
func (d *documentService) CreateDocument(request *models.DocumentRequest, foldername string, filepath string) (*models.DocumentResponse, *errors.Error) {
	    fmt.Println("uploadHandlerService called")
	
		document := request.ReqToDocumentModel()
		// document.FileID = fileFolderLink.FileID
		// document.FolderID = fileFolderLink.FolderID
		response, err := d.UploadDocRepo.CreateDocument(document)   
		if err != nil {
			return nil, errors.ErrInternalServerError
		}
		return response.DocumentToResponse(), nil
}

func (m *documentService) GetAllDocuments() ([]models.DocumentResponse, *errors.Error) {
	var documentResponses []models.DocumentResponse

	documents, err := m.UploadDocRepo.GetAllDocuments()
	if err != nil {
		log.Printf("error getting all medication history of user: %v", err)
		return nil, errors.ErrInternalServerError
	}

	for _, document := range documents {
		documentResponses = append(documentResponses, *document.DocumentToResponse())    
	}
	return documentResponses, nil
}

func (m *documentService) GetDocumentByFolderName(folderName string) ([]models.DocumentResponse, *errors.Error) {
	var documentResponses []models.DocumentResponse

	documents, err := m.UploadDocRepo.GetDocumentByFolderName(folderName)
	if err != nil {
		log.Printf("error getting all medication history of user: %v", err)
		return nil, errors.ErrInternalServerError
	}

	for _, document := range documents {
		documentResponses = append(documentResponses, *document.DocumentToResponse())    
	}
	return documentResponses, nil
}