-- Add book ratings table
-- This migration adds a new table for storing book ratings and reviews

CREATE TABLE IF NOT EXISTS book_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL,
    user_id UUID NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_book_ratings_book 
        FOREIGN KEY (book_id) 
        REFERENCES books(id) 
        ON UPDATE CASCADE 
        ON DELETE CASCADE,
    
    -- Unique constraint to prevent duplicate ratings from same user
    CONSTRAINT unique_user_book_rating UNIQUE (user_id, book_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_book_ratings_book_id ON book_ratings(book_id);
CREATE INDEX IF NOT EXISTS idx_book_ratings_user_id ON book_ratings(user_id);
CREATE INDEX IF NOT EXISTS idx_book_ratings_rating ON book_ratings(rating);
CREATE INDEX IF NOT EXISTS idx_book_ratings_deleted_at ON book_ratings(deleted_at);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_book_ratings_updated_at 
    BEFORE UPDATE ON book_ratings 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
