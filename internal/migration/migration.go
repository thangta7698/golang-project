package migration

import (
	"go-training-system/internal/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Team{},
		// &model.TeamMember{},
		// &model.TeamManager{},
		&model.Folder{},
		&model.Note{},
		&model.FolderShare{},
		&model.NoteShare{},
        &model.TeamUser{},
	)
}
