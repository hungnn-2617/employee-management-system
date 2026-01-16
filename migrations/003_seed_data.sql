
-- Insert sample departments
INSERT INTO departments (id, name) VALUES
    ('d1111111-1111-1111-1111-111111111111', 'Engineering'),
    ('d2222222-2222-2222-2222-222222222222', 'Marketing'),
    ('d3333333-3333-3333-3333-333333333333', 'Human Resources'),
    ('d4444444-4444-4444-4444-444444444444', 'Finance');

-- Insert sample employees
INSERT INTO employees (id, name, age, position, department_id, salary) VALUES
    ('e1111111-1111-1111-1111-111111111111', 'Nguyen Van A', 28, 'Backend Engineer', 'd1111111-1111-1111-1111-111111111111', 2000.00),
    ('e2222222-2222-2222-2222-222222222222', 'Tran Thi B', 25, 'Frontend Engineer', 'd1111111-1111-1111-1111-111111111111', 1800.00),
    ('e3333333-3333-3333-3333-333333333333', 'Le Van C', 30, 'Marketing Manager', 'd2222222-2222-2222-2222-222222222222', 2500.00),
    ('e4444444-4444-4444-4444-444444444444', 'Pham Thi D', 27, 'HR Specialist', 'd3333333-3333-3333-3333-333333333333', 1500.00),
    ('e5555555-5555-5555-5555-555555555555', 'Hoang Van E', 35, 'Finance Director', 'd4444444-4444-4444-4444-444444444444', 3500.00);
