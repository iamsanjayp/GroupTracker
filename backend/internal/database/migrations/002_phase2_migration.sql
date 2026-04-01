-- Phase 2 Migrations: Attendance, Missed OTP, and RBAC

USE grouptracker;

-- 1. Update Users Table
ALTER TABLE users ADD COLUMN IF NOT EXISTS roll_no VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS join_status ENUM('pending', 'approved') DEFAULT 'approved';

-- 2. Attendances Table
CREATE TABLE IF NOT EXISTS attendances (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    team_id     BIGINT UNSIGNED NOT NULL,
    user_id     BIGINT UNSIGNED NOT NULL,
    date        DATE NOT NULL,
    hour_slot   TINYINT NOT NULL CHECK (hour_slot BETWEEN 1 AND 7),
    status      ENUM('Present', 'Absent', 'PS Slot', 'Event', 'OnDuty', 'Class') NOT NULL DEFAULT 'Present',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uniq_attendance (team_id, date, hour_slot, user_id),
    INDEX idx_attendance_team_date (team_id, date)
) ENGINE=InnoDB;

-- 3. Missed Attendances (OTP) Table
CREATE TABLE IF NOT EXISTS missed_attendances (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    team_id     BIGINT UNSIGNED NOT NULL,
    user_id     BIGINT UNSIGNED NOT NULL,
    date        DATE NOT NULL,
    hour_slot   TINYINT NOT NULL CHECK (hour_slot BETWEEN 1 AND 7),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uniq_missed_attendance (user_id, date, hour_slot),
    INDEX idx_missed_team_date (team_id, date)
) ENGINE=InnoDB;
