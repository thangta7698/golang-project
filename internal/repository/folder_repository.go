package repository

import (
	"context"

	"go-training-system/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *model.Folder) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Folder, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Folder, error)
	Update(ctx context.Context, folder *model.Folder) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSharedWithUser(ctx context.Context, userID uuid.UUID) ([]model.Folder, error)
}

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{db: db}
}

func (r *folderRepository) Create(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

func (r *folderRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	var folder model.Folder
	err := r.db.WithContext(ctx).Preload("Owner").Preload("Notes").Preload("Shares").First(&folder, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *folderRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Folder, error) {
	var folders []model.Folder
	err := r.db.WithContext(ctx).Preload("Notes").Where("owner_id = ?", ownerID).Find(&folders).Error
	return folders, err
}

func (r *folderRepository) Update(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

func (r *folderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Folder{}, "id = ?", id).Error
}

func (r *folderRepository) GetSharedWithUser(ctx context.Context, userID uuid.UUID) ([]model.Folder, error) {
	var folders []model.Folder
	err := r.db.WithContext(ctx).
		Joins("JOIN folder_shares ON folders.id = folder_shares.folder_id").
		Where("folder_shares.user_id = ?", userID).
		Preload("Owner").
		Find(&folders).Error
	return folders, err
}
