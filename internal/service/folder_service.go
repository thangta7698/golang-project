package service

import (
	"context"
	"errors"

	"go-training-system/internal/dto"
	"go-training-system/internal/model"
	"go-training-system/internal/repository"

	"github.com/google/uuid"
)

type FolderService interface {
	CreateFolder(ctx context.Context, req *dto.CreateFolderRequest, ownerID uuid.UUID) (*dto.FolderResponse, error)
	GetFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.FolderResponse, error)
	GetUserFolders(ctx context.Context, userID uuid.UUID) ([]dto.FolderResponse, error)
	UpdateFolder(ctx context.Context, id uuid.UUID, req *dto.UpdateFolderRequest, userID uuid.UUID) (*dto.FolderResponse, error)
	DeleteFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ShareFolder(ctx context.Context, folderID uuid.UUID, req *dto.ShareRequest, sharedByID uuid.UUID) error
}

type folderService struct {
	folderRepo repository.FolderRepository
	userRepo   repository.UserRepository
}

func NewFolderService(folderRepo repository.FolderRepository, userRepo repository.UserRepository) FolderService {
	return &folderService{
		folderRepo: folderRepo,
		userRepo:   userRepo,
	}
}

func (s *folderService) CreateFolder(ctx context.Context, req *dto.CreateFolderRequest, ownerID uuid.UUID) (*dto.FolderResponse, error) {
	folder := &model.Folder{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
	}

	if err := s.folderRepo.Create(ctx, folder); err != nil {
		return nil, err
	}

	return &dto.FolderResponse{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		OwnerID:     folder.OwnerID,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}, nil
}

func (s *folderService) GetFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.FolderResponse, error) {
	folder, err := s.folderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user has access
	if folder.OwnerID != userID {
		hasAccess := false
		for _, share := range folder.Shares {
			if share.UserID == userID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return nil, errors.New("access denied")
		}
	}

	return &dto.FolderResponse{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		OwnerID:     folder.OwnerID,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}, nil
}

func (s *folderService) GetUserFolders(ctx context.Context, userID uuid.UUID) ([]dto.FolderResponse, error) {
	// Get owned folders
	ownedFolders, err := s.folderRepo.GetByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get shared folders
	sharedFolders, err := s.folderRepo.GetSharedWithUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Combine and convert to response
	allFolders := append(ownedFolders, sharedFolders...)
	response := make([]dto.FolderResponse, len(allFolders))
	for i, folder := range allFolders {
		response[i] = dto.FolderResponse{
			ID:          folder.ID,
			Name:        folder.Name,
			Description: folder.Description,
			OwnerID:     folder.OwnerID,
			CreatedAt:   folder.CreatedAt,
			UpdatedAt:   folder.UpdatedAt,
		}
	}

	return response, nil
}

func (s *folderService) UpdateFolder(ctx context.Context, id uuid.UUID, req *dto.UpdateFolderRequest, userID uuid.UUID) (*dto.FolderResponse, error) {
	folder, err := s.folderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if folder.OwnerID != userID {
		return nil, errors.New("access denied")
	}

	folder.Name = req.Name
	folder.Description = req.Description

	if err := s.folderRepo.Update(ctx, folder); err != nil {
		return nil, err
	}

	return &dto.FolderResponse{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		OwnerID:     folder.OwnerID,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}, nil
}

func (s *folderService) DeleteFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	folder, err := s.folderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if folder.OwnerID != userID {
		return errors.New("access denied")
	}

	return s.folderRepo.Delete(ctx, id)
}

func (s *folderService) ShareFolder(ctx context.Context, folderID uuid.UUID, req *dto.ShareRequest, sharedByID uuid.UUID) error {
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return err
	}

	if folder.OwnerID != sharedByID {
		return errors.New("access denied")
	}

	// Implementation for sharing logic would go here
	// You'd need a FolderShareRepository for this
	return nil
}
