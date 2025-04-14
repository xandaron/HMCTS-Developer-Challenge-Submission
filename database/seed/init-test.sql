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

-- Add a demo users (password: 'demo123')
INSERT INTO users (id, name, password_hash) VALUES
(1, 'testuser1', '$argon2id$v=19$m=65536,t=3,p=4$pqJ2kWwyHs6Uszb0saO8wQ==$i93hewu5pDcYtrEjUSaKnd6yB00FLwIjWzpuOK5o9/Q='),
(2, 'testuser2', '$argon2id$v=19$m=65536,t=3,p=4$pqJ2kWwyHs6Uszb0saO8wQ==$i93hewu5pDcYtrEjUSaKnd6yB00FLwIjWzpuOK5o9/Q=');

-- Insert mock tasks for testing
INSERT INTO tasks (id, user_id, name, description, status, deadline) VALUES
(1, 1, 'Task 1', 'Description for Task 1', 'INCOMPLETE', '2025-12-31 00:00:00'),
(2, 1, 'Task 2', 'Description for Task 2', 'COMPLETE', '2025-11-30 00:00:00'),
(3, 2, 'Task 3', 'Description for Task 3', 'INCOMPLETE', '2025-10-15 00:00:00');