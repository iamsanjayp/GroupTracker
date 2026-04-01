package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"grouptracker/internal/models"
)

type AttendanceRepo struct {
	db *sql.DB
}

func NewAttendanceRepo(db *sql.DB) *AttendanceRepo {
	return &AttendanceRepo{db: db}
}

func (r *AttendanceRepo) SaveBulk(teamID uint64, date string, records []models.Attendance) error {
	if len(records) == 0 {
		return nil
	}

	query := `INSERT INTO attendances (team_id, user_id, date, hour_slot, status) VALUES `
	vals := []interface{}{}
	placeholders := []string{}

	for _, rec := range records {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
		vals = append(vals, teamID, rec.UserID, date, rec.HourSlot, rec.Status)
	}

	query += strings.Join(placeholders, ", ")
	query += ` ON DUPLICATE KEY UPDATE status = VALUES(status)`

	_, err := r.db.Exec(query, vals...)
	return err
}

func (r *AttendanceRepo) GetByDateAndSession(teamID uint64, date string, session string) ([]models.Attendance, error) {
	// Morning: hours 1-4, Afternoon: hours 5-7
	var query string
	var rows *sql.Rows
	var err error

	if session == "morning" {
		query = `SELECT id, team_id, user_id, date, hour_slot, status FROM attendances 
				 WHERE team_id = ? AND date = ? AND hour_slot BETWEEN 1 AND 4`
	} else {
		query = `SELECT id, team_id, user_id, date, hour_slot, status FROM attendances 
				 WHERE team_id = ? AND date = ? AND hour_slot BETWEEN 5 AND 7`
	}

	rows, err = r.db.Query(query, teamID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.Attendance
	for rows.Next() {
		var rec models.Attendance
		// We format the date out as string using time pkg but mysql driver usually maps DATE to string if specified
		if err := rows.Scan(&rec.ID, &rec.TeamID, &rec.UserID, &rec.Date, &rec.HourSlot, &rec.Status); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *AttendanceRepo) LogMissedOTP(teamID, userID uint64, date string, hourSlot int) error {
	query := `INSERT INTO missed_attendances (team_id, user_id, date, hour_slot) 
			  VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id`
	_, err := r.db.Exec(query, teamID, userID, date, hourSlot)
	return err
}

func (r *AttendanceRepo) GetMissedOTPExports(teamID uint64) ([]models.MissedAttendanceExport, error) {
	query := `
		SELECT 
			DATE_FORMAT(m.date, '%Y-%m-%d') as Date,
			COALESCE(u.roll_no, '') as 'Roll No',
			u.name as Name,
			u.email as 'Mail Id',
			m.hour_slot as Hour
		FROM missed_attendances m
		JOIN users u ON m.user_id = u.id
		WHERE m.team_id = ?
		ORDER BY m.date DESC, m.hour_slot ASC, u.name ASC
	`
	rows, err := r.db.Query(query, teamID)
	if err != nil {
		fmt.Printf("GetMissedOTPExports Error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var exports []models.MissedAttendanceExport
	for rows.Next() {
		var exp models.MissedAttendanceExport
		if err := rows.Scan(&exp.Date, &exp.RollNo, &exp.Name, &exp.Email, &exp.HourSlot); err != nil {
			fmt.Printf("GetMissedOTPExports Scan Error: %v\n", err)
			return nil, err
		}
		exports = append(exports, exp)
	}
	return exports, nil
}

func (r *AttendanceRepo) GetUserAttendanceStats(teamID, userID uint64) (float64, error) {
	var totalSlots int
	var presentSlots int

	query := `
		SELECT 
			COUNT(*) as total_slots,
			COALESCE(SUM(CASE WHEN status IN ('Present', 'PS Slot', 'Event', 'OnDuty', 'Class') THEN 1 ELSE 0 END), 0) as present_slots
		FROM attendances
		WHERE team_id = ? AND user_id = ?
	`
	
	err := r.db.QueryRow(query, teamID, userID).Scan(&totalSlots, &presentSlots)
	if err != nil {
		if err == sql.ErrNoRows {
			return 100.0, nil // No attendance taken yet, assume 100%
		}
		return 0, err
	}

	if totalSlots == 0 {
		return 100.0, nil
	}

	percentage := float64(presentSlots) / float64(totalSlots) * 100.0
	return percentage, nil
}
