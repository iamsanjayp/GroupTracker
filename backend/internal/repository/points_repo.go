package repository

import (
	"database/sql"
	"strings"

	"grouptracker/internal/models"
)

type PointsRepo struct {
	db *sql.DB
}

func NewPointsRepo(db *sql.DB) *PointsRepo {
	return &PointsRepo{db: db}
}

func (r *PointsRepo) RecalculateUser(userID, teamID uint64) error {
	_, err := r.db.Exec(
		`INSERT INTO points (user_id, team_id, total_activity, total_reward)
		 SELECT ?, ?,
		   COALESCE((SELECT SUM(activity_points) FROM activities WHERE user_id = ? AND team_id = ?), 0) +
		   COALESCE((SELECT SUM(activity_points) FROM ps_records WHERE user_id = ? AND team_id = ?), 0),
		   COALESCE((SELECT SUM(reward_points) FROM activities WHERE user_id = ? AND team_id = ?), 0) +
		   COALESCE((SELECT SUM(reward_points) FROM ps_records WHERE user_id = ? AND team_id = ?), 0)
		 ON DUPLICATE KEY UPDATE
		   total_activity = VALUES(total_activity),
		   total_reward = VALUES(total_reward)`,
		userID, teamID,
		userID, teamID, userID, teamID,
		userID, teamID, userID, teamID,
	)
	return err
}

func (r *PointsRepo) GetUserPoints(userID, teamID uint64) (*models.Points, error) {
	p := &models.Points{}
	err := r.db.QueryRow(
		`SELECT id, user_id, team_id, total_activity, total_reward
		 FROM points WHERE user_id = ? AND team_id = ?`,
		userID, teamID,
	).Scan(&p.ID, &p.UserID, &p.TeamID, &p.TotalActivity, &p.TotalReward)
	if err == sql.ErrNoRows {
		return &models.Points{UserID: userID, TeamID: teamID}, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PointsRepo) GetTeamLeaderboard(teamID uint64) ([]models.PointsSummary, error) {
	rows, err := r.db.Query(
		`SELECT p.user_id, u.name, p.total_activity, p.total_reward,
		        (p.total_activity + p.total_reward) as total_points
		 FROM points p
		 JOIN users u ON p.user_id = u.id
		 WHERE p.team_id = ?
		 ORDER BY total_points DESC`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []models.PointsSummary
	for rows.Next() {
		s := models.PointsSummary{}
		if err := rows.Scan(&s.UserID, &s.Name, &s.TotalActivity, &s.TotalReward, &s.TotalPoints); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, s)
	}
	return leaderboard, nil
}

// PS Records
func (r *PointsRepo) CreatePSRecord(userID, teamID uint64, req models.CreatePSRecordRequest) error {
	_, err := r.db.Exec(
		`INSERT INTO ps_records (user_id, team_id, course_name, level, reward_points, activity_points)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		userID, teamID, req.CourseName, req.Level, req.RewardPoints, req.ActivityPoints,
	)
	return err
}

func (r *PointsRepo) GetPSRecords(userID, teamID uint64) ([]models.PSRecord, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, team_id, course_name, level, reward_points, activity_points, completed_at
		 FROM ps_records WHERE user_id = ? AND team_id = ? ORDER BY completed_at DESC`,
		userID, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.PSRecord
	for rows.Next() {
		r := models.PSRecord{}
		if err := rows.Scan(&r.ID, &r.UserID, &r.TeamID, &r.CourseName, &r.Level,
			&r.RewardPoints, &r.ActivityPoints, &r.CompletedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

func (r *PointsRepo) GetPointsHistory(userID, teamID uint64, page, limit int) (*models.PointHistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var totalAct int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM activities WHERE user_id = ? AND team_id = ? AND (activity_points > 0 OR reward_points > 0)`, userID, teamID).Scan(&totalAct)
	if err != nil {
		return nil, err
	}

	var totalMan int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM ps_records WHERE user_id = ? AND team_id = ? AND (activity_points > 0 OR reward_points > 0)`, userID, teamID).Scan(&totalMan)
	if err != nil {
		return nil, err
	}

	totalCount := totalAct + totalMan
	totalPages := (totalCount + limit - 1) / limit

	query := `
		SELECT * FROM (
			SELECT 
				CONCAT('act_', id) as id,
				'daily_log' as source,
				activity_type as reason,
				DATE_FORMAT(activity_date, '%Y-%m-%d') as date,
				activity_points,
				reward_points,
				created_at
			FROM activities 
			WHERE user_id = ? AND team_id = ? AND (activity_points > 0 OR reward_points > 0)
			
			UNION ALL
			
			SELECT 
				CONCAT('man_', id) as id,
				'manual' as source,
				course_name as reason,
				DATE_FORMAT(completed_at, '%Y-%m-%d') as date,
				activity_points,
				reward_points,
				completed_at as created_at
			FROM ps_records 
			WHERE user_id = ? AND team_id = ? AND (activity_points > 0 OR reward_points > 0)
		) as t
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, teamID, userID, teamID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.PointTransaction
	for rows.Next() {
		var tx models.PointTransaction
		if err := rows.Scan(&tx.ID, &tx.Source, &tx.Reason, &tx.Date, &tx.ActivityPoints, &tx.RewardPoints, &tx.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return &models.PointHistoryResponse{
		Transactions: transactions,
		TotalCount:   totalCount,
		Page:         page,
		TotalPages:   totalPages,
	}, nil
}

func (r *PointsRepo) BulkAddPoints(teamID uint64, records []models.BulkPointRecord) (*models.BulkPointResponse, error) {
	resp := &models.BulkPointResponse{
		SuccessCount: 0,
		FailedRows:   []string{},
	}

	// Fetch all users in team to match by Email and RollNo
	rows, err := r.db.Query(`SELECT id, email, roll_no FROM users WHERE team_id = ?`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type userKey struct{ email, roll string }
	userMap := make(map[userKey]uint64)
	for rows.Next() {
		var id uint64
		var email string
		var rollNo *string
		if err := rows.Scan(&id, &email, &rollNo); err == nil {
			r := ""
			if rollNo != nil {
				r = *rollNo
			}
			userMap[userKey{strings.ToLower(email), r}] = id
		}
	}

	for _, rec := range records {
		userID, ok := userMap[userKey{strings.ToLower(rec.Email), rec.RollNo}]
		if !ok {
			resp.FailedRows = append(resp.FailedRows, rec.Email+" ("+rec.RollNo+") [User not found in DB]")
			continue
		}

		// Insert as manual bulk activity using ps_records to allow unlimited uploads per day
		_, err := r.db.Exec(
			`INSERT INTO ps_records (user_id, team_id, course_name, level, activity_points, reward_points) 
			 VALUES (?, ?, ?, 1, ?, ?)`,
			userID, teamID, "[Bulk] "+rec.Reason, rec.ActivityPoints, rec.RewardPoints,
		)
		if err != nil {
			resp.FailedRows = append(resp.FailedRows, rec.Email+" ("+rec.RollNo+") ["+err.Error()+"]")
			continue
		}
		
		_ = r.RecalculateUser(userID, teamID)
		resp.SuccessCount++
	}

	return resp, nil
}
