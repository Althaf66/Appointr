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
    -- mentor_id BIGINT REFERENCES mentors(userid),
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

-- Working At table (one-to-one relationship with mentor)
CREATE TABLE working_at (
    id SERIAL PRIMARY KEY,
    userid BIGINT UNIQUE REFERENCES mentors(userid),
    title VARCHAR(255),
    company VARCHAR(255)
);

-- Social Media table
CREATE TABLE social_media (
    id SERIAL PRIMARY KEY,
    userid BIGINT REFERENCES mentors(userid),
    name VARCHAR(100),
    link VARCHAR(255),
    CONSTRAINT unique_social UNIQUE (userid, name)
);

-- Years of Experience table (one-to-one relationship with mentor)
CREATE TABLE years_of_experience (
    id SERIAL PRIMARY KEY,
    userid BIGINT UNIQUE REFERENCES mentors(userid),
    year BIGINT,
    month BIGINT
);

-- Add indexes for better query performance
CREATE INDEX idx_mentor_userid ON mentors(userid);
CREATE INDEX idx_gigs_userid ON gigs(userid);
CREATE INDEX idx_education_userid ON education(userid);
CREATE INDEX idx_experience_userid ON experience(userid);
CREATE INDEX idx_social_media_userid ON social_media(userid);