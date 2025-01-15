package models

import(
    "time"
    "github.com/google/uuid"
)


type File struct{
    ID              string      `json:"id"`
    UserID          string      `json:"user_id"`
    Filename        string      `json:"filename"`
    Type            string      `json:"type"`
    EncryptedData   string      `json:"encrypted_data"`
    Tags            []string    `json:"tags"`
    CreatedAt       time.Time   `json:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at"`
}


func (f *File) AddDefaults() {
    f.ID = uuid.New().String()
    f.CreatedAt = time.Now()
    f.UpdatedAt = time.Now()

}

type Folder struct {
    ID              string      `json:"id"`
    Name            string      `json:"name"`
    Type            string      `json:"type"`
    ParentFolderID  string      `json:"parent_folder_id"`
    Files           []File      `json:"files"`
    CreatedAt       time.Time   `json:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at"`
}

func (f *Folder) AddDefaults() {
    f.ID = uuid.New().String()
    f.Type = "folder"
    f.CreatedAt = time.Now()
    f.UpdatedAt = time.Now()
}

