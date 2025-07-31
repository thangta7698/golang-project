package repository

import (
	"context"

	"go-training-system/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Note, error)
	GetByFolderID(ctx context.Context, folderID uuid.UUID) ([]model.Note, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSharedWithUser(ctx context.Context, userID uuid.UUID) ([]model.Note, error)
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *noteRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	var note model.Note
	err := r.db.WithContext(ctx).Preload("Owner").Preload("Folder").Preload("Shares").First(&note, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) GetByFolderID(ctx context.Context, folderID uuid.UUID) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.WithContext(ctx).Where("folder_id = ?", folderID).Find(&notes).Error
	return notes, err
}

func (r *noteRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&notes).Error
	return notes, err
}

func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *noteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Note{}, "id = ?", id).Error
}

func (r *noteRepository) GetSharedWithUser(ctx context.Context, userID uuid.UUID) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.WithContext(ctx).
		Joins("JOIN note_shares ON notes.id = note_shares.note_id").
		Where("note_shares.user_id = ?", userID).
		Preload("Owner").
		Preload("Folder").
		Find(&notes).Error
	return notes, err
}
