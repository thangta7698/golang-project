package service

import (
	"context"
	"errors"

	"go-training-system/internal/dto"
	"go-training-system/internal/model"
	"go-training-system/internal/repository"

	"github.com/google/uuid"
)

type NoteService interface {
	CreateNote(ctx context.Context, req *dto.CreateNoteRequest, ownerID uuid.UUID) (*dto.NoteResponse, error)
	GetNote(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.NoteResponse, error)
	GetUserNotes(ctx context.Context, userID uuid.UUID) ([]dto.NoteResponse, error)
	GetFolderNotes(ctx context.Context, folderID uuid.UUID, userID uuid.UUID) ([]dto.NoteResponse, error)
	UpdateNote(ctx context.Context, id uuid.UUID, req *dto.UpdateNoteRequest, userID uuid.UUID) (*dto.NoteResponse, error)
	DeleteNote(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ShareNote(ctx context.Context, noteID uuid.UUID, req *dto.ShareRequest, sharedByID uuid.UUID) error
}

type noteService struct {
	noteRepo   repository.NoteRepository
	folderRepo repository.FolderRepository
}

func NewNoteService(noteRepo repository.NoteRepository, folderRepo repository.FolderRepository) NoteService {
	return &noteService{
		noteRepo:   noteRepo,
		folderRepo: folderRepo,
	}
}

func (s *noteService) CreateNote(ctx context.Context, req *dto.CreateNoteRequest, ownerID uuid.UUID) (*dto.NoteResponse, error) {
	// Check if user has access to the folder
	folder, err := s.folderRepo.GetByID(ctx, req.FolderID)
	if err != nil {
		return nil, err
	}

	hasAccess := folder.OwnerID == ownerID
	if !hasAccess {
		for _, share := range folder.Shares {
			if share.UserID == ownerID && share.Access == model.AccessLevelWrite {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		return nil, errors.New("access denied to folder")
	}

	note := &model.Note{
		Title:    req.Title,
		Body:     req.Body,
		FolderID: req.FolderID,
		OwnerID:  ownerID,
	}

	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, err
	}

	return &dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Body:      note.Body,
		FolderID:  note.FolderID,
		OwnerID:   note.OwnerID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}

func (s *noteService) GetNote(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.NoteResponse, error) {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check access
	if note.OwnerID != userID {
		hasAccess := false
		for _, share := range note.Shares {
			if share.UserID == userID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return nil, errors.New("access denied")
		}
	}

	return &dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Body:      note.Body,
		FolderID:  note.FolderID,
		OwnerID:   note.OwnerID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}

func (s *noteService) GetUserNotes(ctx context.Context, userID uuid.UUID) ([]dto.NoteResponse, error) {
	ownedNotes, err := s.noteRepo.GetByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sharedNotes, err := s.noteRepo.GetSharedWithUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	allNotes := append(ownedNotes, sharedNotes...)
	response := make([]dto.NoteResponse, len(allNotes))
	for i, note := range allNotes {
		response[i] = dto.NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Body:      note.Body,
			FolderID:  note.FolderID,
			OwnerID:   note.OwnerID,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		}
	}

	return response, nil
}

func (s *noteService) GetFolderNotes(ctx context.Context, folderID uuid.UUID, userID uuid.UUID) ([]dto.NoteResponse, error) {
	// Check folder access first
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	hasAccess := folder.OwnerID == userID
	if !hasAccess {
		for _, share := range folder.Shares {
			if share.UserID == userID {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		return nil, errors.New("access denied to folder")
	}

	notes, err := s.noteRepo.GetByFolderID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.NoteResponse, len(notes))
	for i, note := range notes {
		response[i] = dto.NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Body:      note.Body,
			FolderID:  note.FolderID,
			OwnerID:   note.OwnerID,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		}
	}

	return response, nil
}

func (s *noteService) UpdateNote(ctx context.Context, id uuid.UUID, req *dto.UpdateNoteRequest, userID uuid.UUID) (*dto.NoteResponse, error) {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check write access
	hasWriteAccess := note.OwnerID == userID
	if !hasWriteAccess {
		for _, share := range note.Shares {
			if share.UserID == userID && share.Access == model.AccessLevelWrite {
				hasWriteAccess = true
				break
			}
		}
	}

	if !hasWriteAccess {
		return nil, errors.New("write access denied")
	}

	note.Title = req.Title
	note.Body = req.Body

	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, err
	}

	return &dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Body:      note.Body,
		FolderID:  note.FolderID,
		OwnerID:   note.OwnerID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if note.OwnerID != userID {
		return errors.New("access denied")
	}

	return s.noteRepo.Delete(ctx, id)
}

func (s *noteService) ShareNote(ctx context.Context, noteID uuid.UUID, req *dto.ShareRequest, sharedByID uuid.UUID) error {
	note, err := s.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}

	if note.OwnerID != sharedByID {
		return errors.New("access denied")
	}

	// Implementation for sharing logic would go here
	// You'd need a NoteShareRepository for this
	return nil
}
