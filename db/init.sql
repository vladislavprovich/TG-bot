users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL
    -- other user-related columns
);

urls (
    url_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    original_url TEXT NOT NULL,
    short_url VARCHAR(255) UNIQUE NOT NULL
     -- other user-related columns
);