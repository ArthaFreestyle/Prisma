INSERT INTO roles (id, name, description, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'student', 'Regular student user', '2025-01-01 00:00:00'),
('22222222-2222-2222-2222-222222222222', 'lecturer', 'Lecturer with teaching privileges', '2025-01-01 00:00:00'),
('33333333-3333-3333-3333-333333333333', 'admin', 'Full system administrator', '2025-01-01 00:00:00');

INSERT INTO users (username, email, password_hash, full_name, role_id) VALUES
('alan', 'alan@example.com', '$2a$12$qNqVWvZsirkDGFVWVsB0ueKKV3vKxR2xf823dHWeqSp5PWXHxGOAS', 'Alan Pratama', '11111111-1111-1111-1111-111111111111'),
('budi', 'budi@example.com', '$2a$12$qNqVWvZsirkDGFVWVsB0ueKKV3vKxR2xf823dHWeqSp5PWXHxGOAS', 'Budi Santoso', '22222222-2222-2222-2222-222222222222');
