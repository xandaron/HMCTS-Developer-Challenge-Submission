-- Database schema and initial data for HMCTS Developer Challenge

-- Create tables if they don't exist
CREATE TABLE IF NOT EXISTS users (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id INT UNSIGNED NOT NULL,
  name TINYTEXT NOT NULL,
  description TEXT,
  status ENUM('COMPLETE', 'INCOMPLETE') NOT NULL DEFAULT 'INCOMPLETE',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deadline TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Add a demo user (password: 'demo123')
INSERT INTO users (name, password_hash) VALUES  
('demo', '$argon2id$v=19$m=65536,t=3,p=4$pqJ2kWwyHs6Uszb0saO8wQ==$i93hewu5pDcYtrEjUSaKnd6yB00FLwIjWzpuOK5o9/Q=');