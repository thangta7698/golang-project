package repository

import (
	"context"
	"fmt"

	"go-training-system/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error)
	GetTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error)
	AddMemberToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error
	AddManagerToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error
	RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
	RemoveManagerFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *teamRepository) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error) {
    var team model.Team
    err := r.db.WithContext(ctx).
        Preload("CreatedBy").
        Preload("Users").
        Preload("Users.User").
        Preload("Users.AddedBy").
        First(&team, "id = ?", teamID).Error

    if err != nil {
        return nil, err
    }
    return &team, nil
}

func (r *teamRepository) GetTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error) {
    var teams []model.Team
    err := r.db.WithContext(ctx).
        Joins("LEFT JOIN team_user ON team_user.team_id = teams.id").
        Where("team_user.user_id = ?", userID).
        Preload("CreatedBy").
        Preload("Users").
        Preload("Users.User").
        Group("teams.id").
        Find(&teams).Error

    if err != nil {
        return nil, err
    }
    return teams, nil
}

func (r *teamRepository) AddMemberToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error {
    // Kiểm tra user có tồn tại không
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
        return fmt.Errorf("user not found")
    }

    // Kiểm tra user đã trong team chưa (bất kể role gì)
    var count int64
    r.db.WithContext(ctx).Model(&model.TeamUser{}).
        Where("team_id = ? AND user_id = ?", teamID, userID).
        Count(&count)

    if count > 0 {
        return fmt.Errorf("user is already in the team")
    }

    // Thêm user vào team với role MEMBER
    member := model.TeamUser{
        TeamID:    teamID,
        UserID:    userID,
        Role:      model.UserRoleMember,
        AddedByID: addedBy,
    }
    return r.db.WithContext(ctx).Create(&member).Error
}

func (r *teamRepository) AddManagerToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error {
    // Kiểm tra user có tồn tại không
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
        return fmt.Errorf("user not found")
    }

    // Kiểm tra user đã trong team chưa (bất kể role gì)
    var count int64
    r.db.WithContext(ctx).Model(&model.TeamUser{}).
        Where("team_id = ? AND user_id = ?", teamID, userID).
        Count(&count)

    if count > 0 {
        return fmt.Errorf("user is already in the team")
    }

    // Thêm user vào team với role MANAGER
    manager := model.TeamUser{
        TeamID:    teamID,
        UserID:    userID,
        Role:      model.UserRoleManager,
        AddedByID: addedBy,
    }
    return r.db.WithContext(ctx).Create(&manager).Error
}

func (r *teamRepository) RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
    result := r.db.WithContext(ctx).
        Where("team_id = ? AND user_id = ? AND role = ?", teamID, userID, model.UserRoleMember).
        Delete(&model.TeamUser{})

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("member not found in team")
    }

    return nil
}

func (r *teamRepository) RemoveManagerFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
    result := r.db.WithContext(ctx).
        Where("team_id = ? AND user_id = ? AND role = ?", teamID, userID, model.UserRoleManager).
        Delete(&model.TeamUser{})

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("manager not found in team")
    }

    return nil
}
