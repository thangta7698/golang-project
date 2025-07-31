package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamUser struct {
	TeamID    uuid.UUID `json:"team_id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key"`
	AddedAt   time.Time `json:"added_at" gorm:"default:CURRENT_TIMESTAMP"`
	AddedByID uuid.UUID `json:"added_by_id" gorm:"type:uuid"`
    Role      UserRole  `json:"role" gorm:"type:varchar(20);not null;check:role IN ('MANAGER', 'MEMBER')"`

    // Relationships
    Team     Team `gorm:"foreignKey:TeamID"`
	User     User `gorm:"foreignKey:UserID"`
	AddedBy User `json:"added_by" gorm:"foreignKey:AddedByID"`
}
func (TeamUser) TableName() string {
    return "team_user"
}
