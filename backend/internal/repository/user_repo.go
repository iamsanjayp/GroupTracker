package repository

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"grouptracker/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) (uint64, error) {
	res, err := r.db.Exec(
		`INSERT INTO users (email, name, password_hash, avatar_url, team_id, role, roll_no, join_status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Email, user.Name, user.PasswordHash, user.AvatarURL, user.TeamID, user.Role, user.RollNo, user.JoinStatus,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, avatar_url, team_id, role, is_active, roll_no, join_status, created_at, updated_at
		 FROM users WHERE email = ?`, email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.AvatarURL,
		&user.TeamID, &user.Role, &user.IsActive, &user.RollNo, &user.JoinStatus, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) FindByID(id uint64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, avatar_url, team_id, role, is_active, roll_no, join_status, created_at, updated_at
		 FROM users WHERE id = ?`, id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.AvatarURL,
		&user.TeamID, &user.Role, &user.IsActive, &user.RollNo, &user.JoinStatus, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}



func (r *UserRepo) UpdateTeamAndRole(userID, teamID uint64, role string, joinStatus string) error {
	_, err := r.db.Exec(
		`UPDATE users SET team_id = ?, role = ?, join_status = ? WHERE id = ?`,
		teamID, role, joinStatus, userID,
	)
	return err
}

func (r *UserRepo) UpdateRole(userID, teamID uint64, role string) error {
	_, err := r.db.Exec(
		`UPDATE users SET role = ? WHERE id = ? AND team_id = ?`,
		role, userID, teamID,
	)
	return err
}

func (r *UserRepo) RemoveFromTeam(userID, teamID uint64) error {
	_, err := r.db.Exec(
		`UPDATE users SET team_id = NULL, role = 'member' WHERE id = ? AND team_id = ?`,
		userID, teamID,
	)
	return err
}

func (r *UserRepo) GetTeamMembers(teamID uint64) ([]models.TeamMember, error) {
	rows, err := r.db.Query(
		`SELECT id, name, email, role, avatar_url, is_active, roll_no, join_status
		 FROM users WHERE team_id = ? ORDER BY
		 FIELD(role,'captain','vice_captain','manager','strategist','member'), name`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		m := models.TeamMember{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Role, &m.AvatarURL, &m.IsActive, &m.RollNo, &m.JoinStatus); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// Refresh tokens
func (r *UserRepo) SaveRefreshToken(userID uint64, token string, expiresAt time.Time) error {
	hash := hashToken(token)
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)`,
		userID, hash, expiresAt,
	)
	return err
}

func (r *UserRepo) ValidateRefreshToken(userID uint64, token string) (bool, error) {
	hash := hashToken(token)
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM refresh_tokens WHERE user_id = ? AND token_hash = ? AND expires_at > NOW()`,
		userID, hash,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepo) DeleteRefreshToken(userID uint64, token string) error {
	hash := hashToken(token)
	_, err := r.db.Exec(
		`DELETE FROM refresh_tokens WHERE user_id = ? AND token_hash = ?`,
		userID, hash,
	)
	return err
}

func (r *UserRepo) DeleteAllRefreshTokens(userID uint64) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE user_id = ?`, userID)
	return err
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
