DO $$ BEGIN
    RAISE NOTICE 'Creating tables...';
END $$;

-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    author_id VARCHAR(255) NOT NULL, 
    author_name VARCHAR(100) NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);
DO $$ BEGIN
    RAISE NOTICE 'Created posts table.';
END $$;

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    author_id VARCHAR(255) NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reply_to_comment_id INT,
    FOREIGN KEY (reply_to_comment_id) REFERENCES comments(id) ON DELETE CASCADE 
);
DO $$ BEGIN
    RAISE NOTICE 'Created comments table.';
END $$;

-- Archived posts table
CREATE TABLE IF NOT EXISTS archived_posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    author_id VARCHAR(255) NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    archived_at TIMESTAMP NOT NULL DEFAULT NOW()
);
DO $$ BEGIN
    RAISE NOTICE 'Created archived_posts table.';
END $$;

-- Archived comments table
CREATE TABLE IF NOT EXISTS archived_comments (
    id SERIAL PRIMARY KEY,
    post_id INT REFERENCES archived_posts(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    author_id VARCHAR(255) NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    reply_to_comment_id INT,
    FOREIGN KEY (reply_to_comment_id) REFERENCES archived_comments(id) ON DELETE CASCADE
);
DO $$ BEGIN
    RAISE NOTICE 'Created archived_comments table.';
END $$;
