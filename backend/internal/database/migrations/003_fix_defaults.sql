-- Phase 3: Fix NULL defaults for existing rows and clean up skills ENUM
USE grouptracker;

-- Fix NULL join_status for users created before Phase 2
UPDATE users SET join_status = 'approved' WHERE join_status IS NULL;

-- Fix NULL roll_no 
UPDATE users SET roll_no = '' WHERE roll_no IS NULL;

-- Update skills category ENUM to remove non_ps (keep primary, secondary, special)
ALTER TABLE skills MODIFY COLUMN category ENUM('primary','secondary','special') NOT NULL;
