CREATE TABLE achievement_references (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    mongo_achievement_id VARCHAR(24) NOT NULL,
    status VARCHAR(10) NOT NULL CHECK (status IN ('draft', 'submitted', 'verified', 'rejected')),
    submitted_at TIMESTAMP,
    verified_at TIMESTAMP,
    verified_by UUID,
    rejection_note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    -- Foreign keys
    CONSTRAINT fk_student
        FOREIGN KEY (student_id) REFERENCES students(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_verified_by
        FOREIGN KEY (verified_by) REFERENCES users(id)
        ON DELETE SET NULL
);