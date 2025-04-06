CREATE TABLE mentors (
    id bigserial PRIMARY KEY,
    userid BIGINT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100),
  	language VARCHAR(100) [],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE gigs (
    id bigserial PRIMARY KEY,
    userid BIGINT REFERENCES mentors(userid),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    expertise TEXT,
    discipline VARCHAR(100) [],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE education (
    id SERIAL PRIMARY KEY,
    userid BIGINT REFERENCES mentors(userid),
    year_from VARCHAR(4),
    year_to VARCHAR(4),
    degree VARCHAR(255),
    field VARCHAR(255),
    institute VARCHAR(255)
);

-- Experience table
CREATE TABLE experience (
    id SERIAL PRIMARY KEY,
    userid BIGINT REFERENCES mentors(userid),
    year_from VARCHAR(4),
    year_to VARCHAR(4),
    title VARCHAR(255),
    company VARCHAR(255),
    description TEXT
);

CREATE TABLE workingat (
    id SERIAL PRIMARY KEY,
    userid BIGINT REFERENCES mentors(userid),
    title VARCHAR(255),
    company VARCHAR(255),
    totalyear BIGINT,
    month BIGINT,
    linkedin VARCHAR(255),
    github VARCHAR(255),
    instagram VARCHAR(255)
);

CREATE TABLE bookingslots (
    id SERIAL PRIMARY KEY,
    userid INTEGER NOT NULL,
    days VARCHAR(50)[] NOT NULL,
    start_time VARCHAR(5) NOT NULL,
    start_period VARCHAR(2) NOT NULL CHECK (start_period IN ('AM', 'PM')),
    end_time VARCHAR(5) NOT NULL,
    end_period VARCHAR(2) NOT NULL CHECK (end_period IN ('AM', 'PM')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

ALTER TABLE gigs
ADD COLUMN amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
ADD CONSTRAINT positive_amount CHECK (amount >= 0);

CREATE TABLE meetings (
    id SERIAL PRIMARY KEY,
    userid VARCHAR(255) NOT NULL,
    mentorid VARCHAR(255),
    day VARCHAR(50) NOT NULL,
    date VARCHAR(50) NOT NULL,
    start_time VARCHAR(5) NOT NULL,
    start_period VARCHAR(2) NOT NULL CHECK (start_period IN ('AM', 'PM')),
    isconfirm BOOLEAN DEFAULT FALSE,
    ispaid BOOLEAN DEFAULT FALSE,
    iscompleted BOOLEAN DEFAULT FALSE,
    amount DECIMAL(10, 2) NOT NULL,
    link VARCHAR(255),
    -- CONSTRAINT valid_time_format CHECK (start_time ~ '^[0-2][0-9]:[0-5][0-9]$')
);

-- Add indexes for better query performance
CREATE INDEX idx_mentor_userid ON mentors(userid);
CREATE INDEX idx_gigs_userid ON gigs(userid);
CREATE INDEX idx_education_userid ON education(userid);
CREATE INDEX idx_experience_userid ON experience(userid);
-- CREATE INDEX idx_socialmedia_userid ON socialmedia(userid);