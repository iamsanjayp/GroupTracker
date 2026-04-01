package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"

	"grouptracker/internal/models"
)

type TeamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(name string) (*models.Team, error) {
	code := generateInviteCode()
	res, err := r.db.Exec(
		`INSERT INTO teams (name, invite_code) VALUES (?, ?)`,
		name, code,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &models.Team{
		ID:         uint64(id),
		Name:       name,
		InviteCode: code,
	}, nil
}

func (r *TeamRepo) FindByID(id uint64) (*models.Team, error) {
	team := &models.Team{}
	err := r.db.QueryRow(
		`SELECT id, name, invite_code, created_at, updated_at FROM teams WHERE id = ?`, id,
	).Scan(&team.ID, &team.Name, &team.InviteCode, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *TeamRepo) FindByInviteCode(code string) (*models.Team, error) {
	team := &models.Team{}
	err := r.db.QueryRow(
		`SELECT id, name, invite_code, created_at, updated_at FROM teams WHERE invite_code = ?`, code,
	).Scan(&team.ID, &team.Name, &team.InviteCode, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *TeamRepo) GetMemberCount(teamID uint64) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM users WHERE team_id = ?`, teamID,
	).Scan(&count)
	return count, err
}

func generateInviteCode() string {
	b := make([]byte, 5)
	rand.Read(b)
	return hex.EncodeToString(b)
}
