/* ADD available roll to users table and add 5 roll */
ALTER TABLE users ADD availableRoll int(11) DEFAULT 0 NOT NULL;
UPDATE table users SET availableRoll = 4;