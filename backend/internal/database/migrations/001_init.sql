-- GroupTracker Database Schema
-- Phase 1: Full initialization

CREATE DATABASE IF NOT EXISTS grouptracker;
USE grouptracker;

-- ============================================
-- TEAMS
-- ============================================
CREATE TABLE IF NOT EXISTS teams (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    invite_code VARCHAR(20) UNIQUE NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- ============================================
-- USERS
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    name          VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255),
    avatar_url    VARCHAR(500),
    google_id     VARCHAR(100),
    team_id       BIGINT UNSIGNED,
    role          ENUM('captain','vice_captain','manager','strategist','member') DEFAULT 'member',
    is_active     BOOLEAN DEFAULT TRUE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL,
    INDEX idx_users_team (team_id),
    INDEX idx_users_email (email),
    INDEX idx_users_google (google_id)
) ENGINE=InnoDB;

-- ============================================
-- REFRESH TOKENS
-- ============================================
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_refresh_user (user_id)
) ENGINE=InnoDB;

-- ============================================
-- SKILLS
-- ============================================
CREATE TABLE IF NOT EXISTS skills (
    id       BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name     VARCHAR(100) NOT NULL,
    category ENUM('primary','secondary','special','non_ps') NOT NULL,
    team_id  BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    INDEX idx_skills_team (team_id)
) ENGINE=InnoDB;

-- ============================================
-- USER SKILLS
-- ============================================
CREATE TABLE IF NOT EXISTS user_skills (
    id        BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id   BIGINT UNSIGNED NOT NULL,
    skill_id  BIGINT UNSIGNED NOT NULL,
    team_id   BIGINT UNSIGNED NOT NULL,
    level     INT DEFAULT 0,
    validated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)  REFERENCES users(id)  ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id)  REFERENCES teams(id)  ON DELETE CASCADE,
    UNIQUE KEY uniq_user_skill (user_id, skill_id),
    INDEX idx_uskills_team (team_id)
) ENGINE=InnoDB;

-- ============================================
-- PROJECTS
-- ============================================
CREATE TABLE IF NOT EXISTS projects (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    team_id     BIGINT UNSIGNED NOT NULL,
    status      ENUM('active','completed','on_hold') DEFAULT 'active',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    INDEX idx_projects_team (team_id)
) ENGINE=InnoDB;

-- ============================================
-- PROJECT MEMBERS
-- ============================================
CREATE TABLE IF NOT EXISTS project_members (
    id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id       BIGINT UNSIGNED NOT NULL,
    user_id          BIGINT UNSIGNED NOT NULL,
    team_id          BIGINT UNSIGNED NOT NULL,
    share_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id)    REFERENCES users(id)    ON DELETE CASCADE,
    FOREIGN KEY (team_id)    REFERENCES teams(id)    ON DELETE CASCADE,
    UNIQUE KEY uniq_project_member (project_id, user_id),
    INDEX idx_pmembers_team (team_id)
) ENGINE=InnoDB;

-- ============================================
-- TASKS
-- ============================================
CREATE TABLE IF NOT EXISTS tasks (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id  BIGINT UNSIGNED NOT NULL,
    team_id     BIGINT UNSIGNED NOT NULL,
    title       VARCHAR(200) NOT NULL,
    description TEXT,
    assigned_to BIGINT UNSIGNED,
    status      ENUM('todo','in_progress','review','done') DEFAULT 'todo',
    priority    ENUM('low','medium','high') DEFAULT 'medium',
    due_date    DATE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id)  REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id)     REFERENCES teams(id)    ON DELETE CASCADE,
    FOREIGN KEY (assigned_to) REFERENCES users(id)    ON DELETE SET NULL,
    INDEX idx_tasks_team (team_id),
    INDEX idx_tasks_project (project_id)
) ENGINE=InnoDB;

-- ============================================
-- ACTIVITIES
-- ============================================
CREATE TABLE IF NOT EXISTS activities (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    team_id         BIGINT UNSIGNED NOT NULL,
    activity_date   DATE NOT NULL,
    hour_slot       TINYINT NOT NULL,
    activity_type   ENUM('project_work','ps_slot','self_study','event','class_participation') NOT NULL,
    description     VARCHAR(500) NOT NULL,
    activity_points DECIMAL(6,2) DEFAULT 0.00,
    reward_points   DECIMAL(6,2) DEFAULT 0.00,
    project_id      BIGINT UNSIGNED,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)    REFERENCES users(id)    ON DELETE CASCADE,
    FOREIGN KEY (team_id)    REFERENCES teams(id)    ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
    UNIQUE KEY uniq_user_date_hour (user_id, activity_date, hour_slot),
    INDEX idx_activities_team (team_id),
    INDEX idx_activities_user_date (user_id, activity_date)
) ENGINE=InnoDB;

-- ============================================
-- PS RECORDS
-- ============================================
CREATE TABLE IF NOT EXISTS ps_records (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    team_id         BIGINT UNSIGNED NOT NULL,
    course_name     VARCHAR(200) NOT NULL,
    level           INT NOT NULL,
    reward_points   DECIMAL(6,2) NOT NULL,
    activity_points DECIMAL(6,2) NOT NULL,
    completed_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    INDEX idx_ps_team (team_id),
    INDEX idx_ps_user (user_id)
) ENGINE=InnoDB;

-- ============================================
-- POINTS (Aggregated)
-- ============================================
CREATE TABLE IF NOT EXISTS points (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    team_id         BIGINT UNSIGNED NOT NULL,
    total_activity  DECIMAL(10,2) DEFAULT 0.00,
    total_reward    DECIMAL(10,2) DEFAULT 0.00,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    UNIQUE KEY uniq_user_points (user_id),
    INDEX idx_points_team (team_id)
) ENGINE=InnoDB;
