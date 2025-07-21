-- Users table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       full_name VARCHAR(100) NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       phone_number VARCHAR(20) NOT NULL,
                       password_hash TEXT NOT NULL,
                       role VARCHAR(20) CHECK (role IN ('student', 'agent', 'admin')) NOT NULL,
                       is_email_verified BOOLEAN NOT NULL DEFAULT false,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Properties table
CREATE TABLE properties (
                            id SERIAL PRIMARY KEY,
                            agent_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                            title VARCHAR(150) NOT NULL,
                            description TEXT,
                            price INTEGER NOT NULL,
                            location TEXT NOT NULL,
                            property_type VARCHAR(50) CHECK (property_type IN ('self-contained', 'shared', 'hostel', 'flat')),
                            image_url TEXT,
                            video_url TEXT,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inspection Requests
CREATE TABLE inspection_requests (
                                     id SERIAL PRIMARY KEY,
                                     property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                                     student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                     status VARCHAR(20) CHECK (status IN ('pending', 'confirmed', 'completed', 'cancelled')) DEFAULT 'pending',
                                     requested_date TIMESTAMP NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Reviews table
CREATE TABLE reviews (
                         id SERIAL PRIMARY KEY,
                         student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         rating INTEGER CHECK (rating BETWEEN 1 AND 5),
                         comment TEXT,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Favorites table (optional)
CREATE TABLE favorites (
                           student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                           PRIMARY KEY (student_id, property_id)
);

-- Sessions table
CREATE TABLE sessions (
                       id uuid PRIMARY KEY,
                       refresh_token VARCHAR(100) NOT NULL,
                       user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       user_agent VARCHAR(100) NOT NULL,
                       client_ip VARCHAR(100) NOT NULL,
                       is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
                       expires_at TIMESTAMP NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);

-- Verify Email Table
CREATE TABLE verify_emails (
                               id BIGSERIAL PRIMARY KEY,
                               user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               email TEXT NOT NULL,
                               token TEXT NOT NULL,
                               is_used BOOLEAN NOT NULL DEFAULT false,
                               expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + interval '15 minutes'),
                               created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
