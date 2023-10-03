// I removed the verify password by encryption
package models

import (
	"time"
)

// "errors"
// "fmt"
// "time"
// goval "github.com/go-passwd/validator"
// "github.com/go-playground/locales/en"
// ut "github.com/go-playground/universal-translator"
// "github.com/go-playground/validator/v10"
// enTranslations "github.com/go-playground/validator/v10"
// "github.com/leebenson/conform"
// "golang.org/x/crypto/bcrypt"

type Document struct {
	Model
	FileName           string    `json:"filename" binding:"required,min=2"`
	DocumentType       string    `json:"document_type"`
	DocumentNumber     string    `json:"document_number"`
	Department         string    `json:"department"`
	Division           string    `json:"division"`
	Docclass           string    `json:"docclass"`
	DocumentAuthor     string    `json:"document_author"`
	DocumentUploadDate time.Time `json:"document_upload_date"`
	DateCreated        string    `json:"date_created"`
	DocumentUploadNumber string   `json:"document_upload_number"`
	UploaderName       string    `json:"uploader_name"`
	UserID             uint      `json:"user_id"`
	User               User      `json:"user" gorm:"foreignKey:UserID"`
	Folder       	   string      `json:"folder" gorm:"foreignKey:Foldername"`
	Filepath           string    `json:"filepath"`
}
type DocumentRequest struct {
	Filename       string `json:"filename"`
    Foldername        string    `json:"foldername"`
	DocumentType  string    `json:"doctype"`
	DocumentNumber string `json:"document_number"`
    Department    string    `json:"department"`
    Division      string    `json:"division"`
    Docclass      string    `json:"docclass"`
	DocumentUploadDate    string `json:"document_upload_date"`
	UploaderName string `json:"uploader_name"`
	UserID        uint      `json:"user_id"`
	Filepath	  string	`json:"filepath"`		
}
type DocumentResponse struct {
	ID                     uint   `json:"id"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
	Filename       string `json:"filename"`
	DocumentType  string    `json:"doctype"`
	DocumentNumber string `json:"document_number"`
    Department    string    `json:"department"`
    Division      string    `json:"division"`
    Docclass      string    `json:"docclass"`
	DocumentUploadDate    string `json:"document_upload_date"`
	DocumentUploadNumber  string  `json:"documentuploadnumber"`
	UploaderName string `json:"uploader_name" binding:"required"`
	UserID        uint      `json:"user_id"`
	FolderID string   `json:"folder_id"` // Foreign key
	Folder   Folder `json:"folder" gorm:"foreignKey:FolderID"`
}

func (m *DocumentRequest) ReqToDocumentModel() *Document {
	return &Document{
		FileName: 		 m.Filename,
		DocumentType:       m.DocumentType,
		DocumentNumber:     m.DocumentNumber,
		Department:         m.Department,
		Division:           m.Division,
		Docclass: 		 	m.Docclass,
		UploaderName:       m.UploaderName,
		UserID:             m.UserID,
		Filepath: 			m.Filepath,
		Folder: 			m.Foldername,
	}
}


func (m *Document) DocumentToResponse() *DocumentResponse {
	return &DocumentResponse{
		ID: m.ID,
		CreatedAt: time.Unix(m.CreatedAt, 0).String(),
		Filename: m.FileName,
		DocumentType: m.DocumentType,
		DocumentNumber: m.DocumentNumber,
		Department: m.Department,
		Division: m.Division,
		Docclass: m.Docclass,
		DocumentUploadDate: m.DocumentUploadDate.String(),
		DocumentUploadNumber: m.DocumentUploadNumber,
		UploaderName: m.UploaderName,
		UserID: m.UserID,
	}
}