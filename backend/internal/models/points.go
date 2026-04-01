package models

type Points struct {
	ID            uint64  `json:"id"`
	UserID        uint64  `json:"user_id"`
	TeamID        uint64  `json:"team_id"`
	TotalActivity float64 `json:"total_activity"`
	TotalReward   float64 `json:"total_reward"`
}

type PointsSummary struct {
	UserID        uint64  `json:"user_id"`
	Name          string  `json:"name"`
	TotalActivity float64 `json:"total_activity"`
	TotalReward   float64 `json:"total_reward"`
	TotalPoints   float64 `json:"total_points"`
}

type PSRecord struct {
	ID             uint64  `json:"id"`
	UserID         uint64  `json:"user_id"`
	TeamID         uint64  `json:"team_id"`
	CourseName     string  `json:"course_name"`
	Level          int     `json:"level"`
	RewardPoints   float64 `json:"reward_points"`
	ActivityPoints float64 `json:"activity_points"`
	CompletedAt    string  `json:"completed_at"`
}

type CreatePSRecordRequest struct {
	CourseName     string  `json:"course_name"`
	Level          int     `json:"level"`
	RewardPoints   float64 `json:"reward_points"`
	ActivityPoints float64 `json:"activity_points"`
}

type PointTransaction struct {
	ID             string  `json:"id"`
	Source         string  `json:"source"`
	Reason         string  `json:"reason"`
	Date           string  `json:"date"`
	ActivityPoints float64 `json:"activity_points"`
	RewardPoints   float64 `json:"reward_points"`
	CreatedAt      string  `json:"created_at"`
}

type PointHistoryResponse struct {
	Transactions []PointTransaction `json:"transactions"`
	TotalCount   int                `json:"total_count"`
	Page         int                `json:"page"`
	TotalPages   int                `json:"total_pages"`
}

type BulkPointRecord struct {
	Email          string  `json:"email"`
	RollNo         string  `json:"roll_no"`
	Reason         string  `json:"reason"`
	ActivityPoints float64 `json:"activity_points"`
	RewardPoints   float64 `json:"reward_points"`
}

type BulkPointUploadRequest struct {
	Records []BulkPointRecord `json:"records"`
}

type BulkPointResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedRows   []string `json:"failed_rows"`
}
