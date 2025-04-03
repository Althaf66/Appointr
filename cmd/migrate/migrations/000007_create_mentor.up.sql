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

-- Social Media table
-- CREATE TABLE socialmedia (
--     id SERIAL PRIMARY KEY,
--     userid BIGINT REFERENCES mentors(userid),
--     name VARCHAR(100),
--     link VARCHAR(255),
--     CONSTRAINT unique_social UNIQUE (userid, name)
-- );

-- Add indexes for better query performance
CREATE INDEX idx_mentor_userid ON mentors(userid);
CREATE INDEX idx_gigs_userid ON gigs(userid);
CREATE INDEX idx_education_userid ON education(userid);
CREATE INDEX idx_experience_userid ON experience(userid);
-- CREATE INDEX idx_socialmedia_userid ON socialmedia(userid);