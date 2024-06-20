-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id BIGINT UNIQUE NOT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS channels (
    name VARCHAR(50)
);

-- Create admins table
CREATE TABLE IF NOT EXISTS admins (
    id BIGINT UNIQUE NOT NULL
);

