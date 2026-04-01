package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type SkillRepo struct {
	db *sql.DB
}

func NewSkillRepo(db *sql.DB) *SkillRepo {
	return &SkillRepo{db: db}
}

// Fixed skill lists — these are the only options available
var normalSkills = []string{
	"Additive Manufacturing (3D Printing)",
	"Agentic AI & LLM Optimization",
	"Autonomous Mobile Robotics (AMR)",
	"Battery Management Systems (BMS)",
	"Big Data Analytics and machine learning",
	"Bio-Process Engineering",
	"Bioinformatics and Data Analytics",
	"Blockchain Technology",
	"Cloud Computing",
	"Computational Fluid Dynamics (CFD)",
	"Computer Vision and Image Processing",
	"Control System",
	"Cyber Security and Cryptography",
	"Data Acquisition System",
	"Design for Manufacturing and Assembly",
	"DevOps and IT Infra",
	"Digital Signal Processing",
	"Edge AI",
	"Embedded Systems & Firmware",
	"FPGA Prototyping",
	"Full-Stack Software Development",
	"IoT and Sensor Integration",
	"Manufacturing and Fabrication",
	"Mechanical Engineering CAD and FEA",
	"Mechanical Modelling",
	"Mechanisms Design",
	"Microbial and Plant Bioprospecting",
	"Molecular Biology and Genetic Engineering",
	"Natural Language Processing",
	"PCB Design and Development",
	"PLC and Industrial Control",
	"Pneumatics & Electro-Pneumatics",
	"Power Electronics & Grid Integration",
	"Power System",
	"Precision Agriculture (Agri-Tech)",
	"Robot Systems Integration",
	"Servo-Drives & Motion Control",
	"Unmanned Aerial Systems",
	"VLSI & Circuit Design",
}

var FixedSkills = map[string][]string{
	"primary":   normalSkills,
	"secondary": normalSkills,
	"special": {
		"Augmented Reality (AR) & Virtual Reality (VR) Development",
		"Business Process Intelligence (BPI)",
		"Continuous Improvement (Lean/Kaizen)",
		"Intellectual Property Rights (IPR)",
		"Prompt Engineering",
		"Quality Tools (Six Sigma/TQM)",
		"Report writing",
		"Research methodology",
		"Generative AI (Gen AI)",
		"User Experience (UI/UX) Design",
	},
}

type UserSkill struct {
	ID        uint64  `json:"id"`
	UserID    uint64  `json:"user_id"`
	TeamID    uint64  `json:"team_id"`
	SkillName string  `json:"skill_name"`
	Category  string  `json:"category"`
	Level     int     `json:"level"`
	Validated bool    `json:"validated"`
}

type SetSkillsRequest struct {
	Primary   []string `json:"primary"`   // exactly 2
	Secondary []string `json:"secondary"` // exactly 2
	Special   []string `json:"special"`   // exactly 2
}

// GetUserSkills returns the skills for a specific user
func (r *SkillRepo) GetUserSkills(userID, teamID uint64) ([]UserSkill, error) {
	rows, err := r.db.Query(
		`SELECT us.id, us.user_id, us.team_id, s.name, s.category, us.level, us.validated
		 FROM user_skills us
		 JOIN skills s ON us.skill_id = s.id
		 WHERE us.user_id = ? AND us.team_id = ?
		 ORDER BY FIELD(s.category, 'primary', 'secondary', 'special'), s.name`,
		userID, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []UserSkill
	for rows.Next() {
		var s UserSkill
		if err := rows.Scan(&s.ID, &s.UserID, &s.TeamID, &s.SkillName, &s.Category, &s.Level, &s.Validated); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}

// GetTeamSkills returns all user skills grouped by user for a team
func (r *SkillRepo) GetTeamSkills(teamID uint64) (map[uint64][]UserSkill, error) {
	rows, err := r.db.Query(
		`SELECT us.id, us.user_id, us.team_id, s.name, s.category, us.level, us.validated
		 FROM user_skills us
		 JOIN skills s ON us.skill_id = s.id
		 WHERE us.team_id = ?
		 ORDER BY us.user_id, FIELD(s.category, 'primary', 'secondary', 'special'), s.name`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uint64][]UserSkill)
	for rows.Next() {
		var s UserSkill
		if err := rows.Scan(&s.ID, &s.UserID, &s.TeamID, &s.SkillName, &s.Category, &s.Level, &s.Validated); err != nil {
			return nil, err
		}
		result[s.UserID] = append(result[s.UserID], s)
	}
	return result, nil
}

// SetUserSkills replaces all skills for a user. Validates the request first.
func (r *SkillRepo) SetUserSkills(userID, teamID uint64, req SetSkillsRequest) error {
	// Validate counts
	if len(req.Primary) != 2 {
		return fmt.Errorf("exactly 2 primary skills required")
	}
	if len(req.Secondary) != 2 {
		return fmt.Errorf("exactly 2 secondary skills required")
	}
	if len(req.Special) != 2 {
		return fmt.Errorf("exactly 2 special skills required")
	}

	// Validate skill names exist in fixed list
	for _, name := range req.Primary {
		if !isValidSkill("primary", name) {
			return fmt.Errorf("invalid primary skill: %s", name)
		}
	}
	for _, name := range req.Secondary {
		if !isValidSkill("secondary", name) {
			return fmt.Errorf("invalid secondary skill: %s", name)
		}
	}
	for _, name := range req.Special {
		if !isValidSkill("special", name) {
			return fmt.Errorf("invalid special skill: %s", name)
		}
	}

	// Use a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing user skills for this team
	_, err = tx.Exec(`DELETE us FROM user_skills us
		JOIN skills s ON us.skill_id = s.id
		WHERE us.user_id = ? AND us.team_id = ?`, userID, teamID)
	if err != nil {
		return err
	}

	// Insert all skills
	allSkills := []struct {
		name     string
		category string
	}{}
	for _, name := range req.Primary {
		allSkills = append(allSkills, struct{ name, category string }{name, "primary"})
	}
	for _, name := range req.Secondary {
		allSkills = append(allSkills, struct{ name, category string }{name, "secondary"})
	}
	for _, name := range req.Special {
		allSkills = append(allSkills, struct{ name, category string }{name, "special"})
	}

	for _, skill := range allSkills {
		// Ensure skill exists in skills table for this team
		var skillID uint64
		err = tx.QueryRow(`SELECT id FROM skills WHERE name = ? AND category = ? AND team_id = ?`,
			skill.name, skill.category, teamID).Scan(&skillID)
		if err == sql.ErrNoRows {
			// Create it
			res, err := tx.Exec(`INSERT INTO skills (name, category, team_id) VALUES (?, ?, ?)`,
				skill.name, skill.category, teamID)
			if err != nil {
				return err
			}
			id, _ := res.LastInsertId()
			skillID = uint64(id)
		} else if err != nil {
			return err
		}

		// Insert user_skill
		_, err = tx.Exec(`INSERT INTO user_skills (user_id, skill_id, team_id, level, validated) VALUES (?, ?, ?, 0, FALSE)
			ON DUPLICATE KEY UPDATE level = level`, userID, skillID, teamID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// HasSkillsSet checks if a user has selected their skills
func (r *SkillRepo) HasSkillsSet(userID, teamID uint64) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM user_skills WHERE user_id = ? AND team_id = ?`, userID, teamID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count >= 6, nil
}

func isValidSkill(category, name string) bool {
	skills, ok := FixedSkills[category]
	if !ok {
		return false
	}
	for _, s := range skills {
		if strings.EqualFold(s, name) {
			return true
		}
	}
	return false
}
