CREATE TABLE IF NOT EXISTS users (
       id UUID PRIMARY KEY,
       first_name VARCHAR(150) NOT NULL,
       last_name VARCHAR(150) NOT NULL,
       age INT NOT NULL,
       email VARCHAR(200) UNIQUE NOT NULL,
       password TEXT NOT NULL ,
       refresh_token VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP,
       deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_unique_email ON users(email) WHERE deleted_at IS NULL;