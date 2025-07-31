package dto

import (
	"time"

	"go-training-system/internal/model"

	"github.com/google/uuid"
)

type CreateFolderRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateFolderRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

type FolderResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateNoteRequest struct {
	Title    string    `json:"title" validate:"required,min=1,max=255"`
	Body     string    `json:"body"`
	FolderID uuid.UUID `json:"folder_id" validate:"required"`
}

type UpdateNoteRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
	Body  string `json:"body"`
}

type NoteResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	FolderID  uuid.UUID `json:"folder_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShareRequest struct {
	UserID uuid.UUID           `json:"user_id" validate:"required"`
	Access model.AccessLevel   `json:"access" validate:"required,oneof=read write"`
}
