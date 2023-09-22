// I removed the verify password by encryption
package models

import "time"

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
	Folder             string    `json:"folder" binding:"required" gorm:"default:null"`
	DocumentType       string    `json:"document_type"`
	DocumentNumber     string    `json:"document_number"`
	Department         string    `json:"department"`
	Division           string    `json:"division"`
	DocClass           string    `json:"docclass"`
	DocumentAuthor     string    `json:"document_author"`
	DocumentUploadDate time.Time `json:"document_upload_date"`
	DateCreated        string    `json:"date_created"`
	DocumentUploadNumber string   `json:"document_upload_number"`
	UploaderName       string    `json:"uploader_name"`
	UserID             uint
}
type DocumentRequest struct {
	Filename       string `json:"filename"`
    Folder        string    `json:"folder"`
	DocumentType  string    `json:"doctype"`
	DocumentNumber string `json:"document_number"`
    Department    string    `json:"department"`
    Division      string    `json:"division"`
    Docclass      string    `json:"docclass"`
	DocumentUploadDate    string `json:"document_upload_date"`
	UploaderName string `json:"uploader_name"`
	UserID        uint      `json:"user_id"`
}
type DocumentResponse struct {
	ID                     uint   `json:"id"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
	Filename       string `json:"filename"`
    Folder        string    `json:"folder"`
	DocumentType  string    `json:"doctype"`
	DocumentNumber string `json:"document_number"`
    Department    string    `json:"department"`
    Division      string    `json:"division"`
    DocClass      string    `json:"docclass"`
	DocumentUploadDate    string `json:"document_upload_date"`
	DocumentUploadNumber  string  `json:"documentuploadnumber"`
	UploaderName string `json:"uploader_name" binding:"required"`
	UserID        uint      `json:"user_id"`
}

func (m *DocumentRequest) ReqToDocumentModel() *Document {
	return &Document{
		FileName: 		 m.Filename,
		Folder:             m.Folder,
		DocumentType:       m.DocumentType,
		DocumentNumber:     m.DocumentNumber,
		Department:         m.Department,
		Division:           m.Division,
		UploaderName:       m.UploaderName,
		UserID:             m.UserID,
	}
}


func (m *Document) DocumentToResponse() *DocumentResponse {
	return &DocumentResponse{
		ID: m.ID,
		Filename: m.FileName,
		Folder: m.Folder,
		DocumentType: m.DocumentType,
		DocumentNumber: m.DocumentNumber,
		Department: m.Department,
		Division: m.Division,
		DocClass: m.DocClass,
		DocumentUploadNumber: m.DocumentUploadNumber,
		UploaderName: m.UploaderName,
		UserID: m.UserID,
	}
}