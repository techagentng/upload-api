package models

import "time"

type Folder struct {
    FolderID        uint      `gorm:"primaryKey" json:"folder_id" gorm:"autoIncrement"`
    Foldername      string      `json:"foldername"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    Document []Document   `json:"document" gorm:"foreignKey:id"`
}