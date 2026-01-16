
CREATE TABLE IF NOT EXISTS employees (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    position VARCHAR(255) NOT NULL,
    department_id CHAR(36) NOT NULL,
    salary DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_employees_deleted_at (deleted_at),
    INDEX idx_employees_department_id (department_id),
    INDEX idx_employees_name (name),
    INDEX idx_employees_position (position)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
