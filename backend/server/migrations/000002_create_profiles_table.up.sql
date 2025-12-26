    CREATE TABLE profiles (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        full_name VARCHAR(255) NOT NULL,
        date_of_birth DATE NOT NULL,
        aadhaar_number VARCHAR(20) NOT NULL UNIQUE,
        phone_number VARCHAR(20) NOT NULL UNIQUE,
        address TEXT,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        
        -- Links profile to the users table
        CONSTRAINT fk_user 
            FOREIGN KEY(user_id) 
            REFERENCES users(id) 
            ON DELETE CASCADE
    );

    -- Index for performance on queries looking up a user's profile
    CREATE INDEX idx_profiles_user_id ON profiles(user_id);