package repository

import (
	"database/sql"

	"grouptracker/internal/models"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(name, description string, teamID uint64) (uint64, error) {
	res, err := r.db.Exec(
		`INSERT INTO projects (name, description, team_id) VALUES (?, ?, ?)`,
		name, description, teamID,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *ProjectRepo) GetByTeam(teamID uint64) ([]models.Project, error) {
	rows, err := r.db.Query(
		`SELECT id, name, description, team_id, status, created_at, updated_at
		 FROM projects WHERE team_id = ? ORDER BY created_at DESC`, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		p := models.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.TeamID, &p.Status,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepo) GetByID(id, teamID uint64) (*models.Project, error) {
	p := &models.Project{}
	err := r.db.QueryRow(
		`SELECT id, name, description, team_id, status, created_at, updated_at
		 FROM projects WHERE id = ? AND team_id = ?`, id, teamID,
	).Scan(&p.ID, &p.Name, &p.Description, &p.TeamID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProjectRepo) Update(id, teamID uint64, req models.UpdateProjectRequest) error {
	query := "UPDATE projects SET"
	args := []interface{}{}
	sets := []string{}

	if req.Name != nil {
		sets = append(sets, " name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		sets = append(sets, " description = ?")
		args = append(args, *req.Description)
	}
	if req.Status != nil {
		sets = append(sets, " status = ?")
		args = append(args, *req.Status)
	}

	if len(sets) == 0 {
		return nil
	}

	for i, s := range sets {
		if i > 0 {
			query += ","
		}
		query += s
	}
	query += " WHERE id = ? AND team_id = ?"
	args = append(args, id, teamID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ProjectRepo) GetMembers(projectID, teamID uint64) ([]models.ProjectMember, error) {
	rows, err := r.db.Query(
		`SELECT pm.id, pm.user_id, u.name, u.email, pm.share_percentage
		 FROM project_members pm
		 JOIN users u ON pm.user_id = u.id
		 WHERE pm.project_id = ? AND pm.team_id = ?`,
		projectID, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMember
	for rows.Next() {
		m := models.ProjectMember{}
		if err := rows.Scan(&m.ID, &m.UserID, &m.Name, &m.Email, &m.SharePercentage); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *ProjectRepo) AddMember(projectID, userID, teamID uint64, share float64) error {
	_, err := r.db.Exec(
		`INSERT INTO project_members (project_id, user_id, team_id, share_percentage) VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE share_percentage = ?`,
		projectID, userID, teamID, share, share,
	)
	return err
}

func (r *ProjectRepo) RemoveMember(projectID, userID, teamID uint64) error {
	_, err := r.db.Exec(
		`DELETE FROM project_members WHERE project_id = ? AND user_id = ? AND team_id = ?`,
		projectID, userID, teamID,
	)
	return err
}

func (r *ProjectRepo) GetTasks(projectID, teamID uint64) ([]models.Task, error) {
	rows, err := r.db.Query(
		`SELECT t.id, t.project_id, t.team_id, t.title, t.description, t.assigned_to,
		        u.name, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		 FROM tasks t
		 LEFT JOIN users u ON t.assigned_to = u.id
		 WHERE t.project_id = ? AND t.team_id = ?
		 ORDER BY FIELD(t.priority,'high','medium','low'), t.created_at`,
		projectID, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		t := models.Task{}
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.TeamID, &t.Title, &t.Description,
			&t.AssignedTo, &t.AssigneeName, &t.Status, &t.Priority, &t.DueDate,
			&t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *ProjectRepo) CreateTask(projectID, teamID uint64, req models.CreateTaskRequest) (uint64, error) {
	res, err := r.db.Exec(
		`INSERT INTO tasks (project_id, team_id, title, description, assigned_to, priority, due_date)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		projectID, teamID, req.Title, req.Description, req.AssignedTo, req.Priority, req.DueDate,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *ProjectRepo) UpdateTask(taskID, teamID uint64, req models.UpdateTaskRequest) error {
	query := "UPDATE tasks SET"
	args := []interface{}{}
	sets := []string{}

	if req.Title != nil {
		sets = append(sets, " title = ?")
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		sets = append(sets, " description = ?")
		args = append(args, *req.Description)
	}
	if req.Status != nil {
		sets = append(sets, " status = ?")
		args = append(args, *req.Status)
	}
	if req.Priority != nil {
		sets = append(sets, " priority = ?")
		args = append(args, *req.Priority)
	}
	if req.AssignedTo != nil {
		sets = append(sets, " assigned_to = ?")
		args = append(args, *req.AssignedTo)
	}
	if req.DueDate != nil {
		sets = append(sets, " due_date = ?")
		args = append(args, *req.DueDate)
	}

	if len(sets) == 0 {
		return nil
	}

	for i, s := range sets {
		if i > 0 {
			query += ","
		}
		query += s
	}
	query += " WHERE id = ? AND team_id = ?"
	args = append(args, taskID, teamID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ProjectRepo) GetProjectCount(teamID uint64) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM projects WHERE team_id = ?`, teamID,
	).Scan(&count)
	return count, err
}

func (r *ProjectRepo) GetTaskStats(teamID uint64) (int, int, error) {
	var total, done int
	err := r.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(status = 'done'), 0) FROM tasks WHERE team_id = ?`, teamID,
	).Scan(&total, &done)
	return total, done, err
}
