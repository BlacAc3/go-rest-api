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

