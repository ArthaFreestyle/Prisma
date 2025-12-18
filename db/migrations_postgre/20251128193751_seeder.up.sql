-- =========================
-- Roles
-- =========================

INSERT INTO roles (id, name, description, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'mahasiswa', 'Regular student user', '2025-01-01 00:00:00'),
('22222222-2222-2222-2222-222222222222', 'lecturer', 'Lecturer with teaching privileges', '2025-01-01 00:00:00'),
('33333333-3333-3333-3333-333333333333', 'admin', 'Full system administrator', '2025-01-01 00:00:00');

-- =========================
-- Users
-- =========================
INSERT INTO users (username, email, password_hash, full_name, role_id) VALUES
('alan', 'alan@example.com', '$2a$12$qNqVWvZsirkDGFVWVsB0ueKKV3vKxR2xf823dHWeqSp5PWXHxGOAS', 'Alan Pratama', '11111111-1111-1111-1111-111111111111'),
('budi', 'budi@example.com', '$2a$12$qNqVWvZsirkDGFVWVsB0ueKKV3vKxR2xf823dHWeqSp5PWXHxGOAS', 'Budi Santoso', '22222222-2222-2222-2222-222222222222');


-- =========================
-- Permissions
-- =========================
INSERT INTO permissions (name, resource, action, description)
VALUES
-- =========================
-- 5.1 Authentication
-- =========================
('auth:login', 'auth', 'login', 'Login untuk mendapatkan access token'),
('auth:refresh', 'auth', 'refresh', 'Refresh access token'),
('auth:logout', 'auth', 'logout', 'Logout dan revoke refresh token'),
('auth:profile', 'auth', 'read', 'Mengambil profil pengguna yang sedang login'),

-- =========================
-- 5.2 Users (Admin)
-- =========================
('users:list', 'users', 'list', 'Melihat semua pengguna'),
('users:detail', 'users', 'detail', 'Melihat detail pengguna'),
('users:create', 'users', 'create', 'Menambahkan pengguna baru'),
('users:update', 'users', 'update', 'Mengubah data pengguna'),
('users:delete', 'users', 'delete', 'Menghapus pengguna'),
('users:updateRole', 'users', 'updateRole', 'Mengubah role pengguna'),

-- =========================
-- 5.4 Achievements
-- =========================
('achievements:list', 'achievements', 'list', 'List semua achievement'),
('achievements:detail', 'achievements', 'detail', 'Detail achievement'),
('achievements:create', 'achievements', 'create', 'Mahasiswa membuat achievement'),
('achievements:update', 'achievements', 'update', 'Mahasiswa mengupdate achievement'),
('achievements:delete', 'achievements', 'delete', 'Mahasiswa menghapus achievement'),
('achievements:submit', 'achievements', 'submit', 'Submit achievement untuk verifikasi'),
('achievements:verify', 'achievements', 'verify', 'Dosen Wali memverifikasi achievement'),
('achievements:reject', 'achievements', 'reject', 'Dosen Wali menolak achievement'),
('achievements:history', 'achievements', 'readHistory', 'Melihat history status achievement'),
('achievements:uploadAttachment', 'achievements', 'upload', 'Upload attachment untuk achievement'),

-- =========================
-- 5.5 Students & Lecturers
-- =========================
('students:list', 'students', 'list', 'List semua mahasiswa'),
('students:detail', 'students', 'detail', 'Detail mahasiswa'),
('students:achievements', 'students', 'readAchievements', 'List achievement mahasiswa'),
('students:updateAdvisor', 'students', 'updateAdvisor', 'Update dosen wali mahasiswa'),

('lecturers:list', 'lecturers', 'list', 'List dosen'),
('lecturers:advisees', 'lecturers', 'readAdvisees', 'List mahasiswa bimbingan'),

-- =========================
-- 5.8 Reports & Analytics
-- =========================
('reports:statistics', 'reports', 'readStatistics', 'Laporan statistik'),
('reports:studentDetail', 'reports', 'readStudentDetail', 'Laporan detail perkembangan mahasiswa');

INSERT INTO role_permissions (role_id, permission_id)
SELECT '11111111-1111-1111-1111-111111111111', id
FROM permissions
WHERE name IN (
               'auth:login',
               'auth:logout',
               'auth:profile',

               'achievements:list',
               'achievements:detail',
               'achievements:create',
               'achievements:update',
               'achievements:delete',
               'achievements:submit',
               'achievements:uploadAttachment',

               'students:detail',           -- lihat info dirinya sendiri
               'students:achievements'      -- lihat achievement dirinya sendiri
    );


INSERT INTO role_permissions (role_id, permission_id)
SELECT '22222222-2222-2222-2222-222222222222', id
FROM permissions
WHERE name IN (
               'auth:login',
               'auth:logout',
               'auth:profile',

               'achievements:list',
               'achievements:detail',
               'achievements:verify',
               'achievements:reject',
               'achievements:history',

               'lecturers:list',
               'lecturers:advisees',

               'students:list',
               'students:detail',
               'students:achievements',

               'reports:statistics',
               'reports:studentDetail'
    );


INSERT INTO role_permissions (role_id, permission_id)
SELECT '33333333-3333-3333-3333-333333333333', id
FROM permissions;


INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at) VALUES
    ('aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
     (SELECT id FROM users WHERE username='budi'),
     'LECT-2025-001',
     'Computer Science',
     '2025-01-01 00:00:00');


INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at) VALUES
    ('bbbbbbb1-bbbb-bbbb-bbbb-bbbbbbbbbbb1',
     (SELECT id FROM users WHERE username='alan'),
     'STUD-2025-001',
     'Information Systems',
     '2025/2026',
     'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1',      -- advisor = lecturer budi
     '2025-01-01 00:00:00');