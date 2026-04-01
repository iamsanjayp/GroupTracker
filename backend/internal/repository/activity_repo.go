package repository

import (
	"database/sql"

	"grouptracker/internal/models"
)

type ActivityRepo struct {
	db *sql.DB
}

func NewActivityRepo(db *sql.DB) *ActivityRepo {
	return &ActivityRepo{db: db}
}

func (r *ActivityRepo) GetDayLog(userID, teamID uint64, date string) ([]models.Activity, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, team_id, activity_date, hour_slot, activity_type,
		        description, activity_points, reward_points, project_id, created_at, updated_at
		 FROM activities
		 WHERE user_id = ? AND team_id = ? AND activity_date = ?
		 ORDER BY hour_slot`,
		userID, teamID, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		a := models.Activity{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.TeamID, &a.ActivityDate, &a.HourSlot,
			&a.ActivityType, &a.Description, &a.ActivityPoints, &a.RewardPoints,
			&a.ProjectID, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}

func (r *ActivityRepo) UpsertActivity(userID, teamID uint64, date string, entry models.ActivityEntry) error {
	_, err := r.db.Exec(
		`INSERT INTO activities (user_id, team_id, activity_date, hour_slot, activity_type,
		                         description, activity_points, reward_points, project_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE
		   activity_type = VALUES(activity_type),
		   description = VALUES(description),
		   activity_points = VALUES(activity_points),
		   reward_points = VALUES(reward_points),
		   project_id = VALUES(project_id)`,
		userID, teamID, date, entry.HourSlot, entry.ActivityType,
		entry.Description, entry.ActivityPoints, entry.RewardPoints, entry.ProjectID,
	)
	return err
}

func (r *ActivityRepo) GetTeamActivities(teamID uint64, date string) ([]models.Activity, error) {
	rows, err := r.db.Query(
		`SELECT a.id, a.user_id, a.team_id, a.activity_date, a.hour_slot, a.activity_type,
		        a.description, a.activity_points, a.reward_points, a.project_id, a.created_at, a.updated_at
		 FROM activities a
		 WHERE a.team_id = ? AND a.activity_date = ?
		 ORDER BY a.user_id, a.hour_slot`,
		teamID, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		a := models.Activity{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.TeamID, &a.ActivityDate, &a.HourSlot,
			&a.ActivityType, &a.Description, &a.ActivityPoints, &a.RewardPoints,
			&a.ProjectID, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}

func (r *ActivityRepo) GetUserTodayPoints(userID, teamID uint64, date string) (float64, float64, error) {
	var actTotal, rewTotal float64
	err := r.db.QueryRow(
		`SELECT COALESCE(SUM(activity_points), 0), COALESCE(SUM(reward_points), 0)
		 FROM activities WHERE user_id = ? AND team_id = ? AND activity_date = ?`,
		userID, teamID, date,
	).Scan(&actTotal, &rewTotal)
	return actTotal, rewTotal, err
}

func (r *ActivityRepo) GetUserLoggedDays(userID, teamID uint64, month string) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(DISTINCT activity_date)
		 FROM activities
		 WHERE user_id = ? AND team_id = ? AND activity_date LIKE ?`,
		userID, teamID, month+"%",
	).Scan(&count)
	return count, err
}

func (r *ActivityRepo) GetTeamActiveToday(teamID uint64, date string) (int, int, error) {
	var active, total int
	err := r.db.QueryRow(
		`SELECT COUNT(DISTINCT user_id) FROM activities WHERE team_id = ? AND activity_date = ?`,
		teamID, date,
	).Scan(&active)
	if err != nil {
		return 0, 0, err
	}
	err = r.db.QueryRow(
		`SELECT COUNT(*) FROM users WHERE team_id = ? AND is_active = TRUE`,
		teamID,
	).Scan(&total)
	return active, total, err
}
