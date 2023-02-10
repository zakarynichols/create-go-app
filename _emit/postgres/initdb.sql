CREATE TYPE thing_type AS ENUM ('abstract', 'concrete');

CREATE TABLE things (
  thing_id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  type thing_type NOT NULL
);

INSERT INTO things (name, description, type) VALUES
    ('iPhone', 'A smartphone made by Apple', 'concrete'),
    ('Honor', 'High respect; great esteem', 'abstract'),
    ('The Great Gatsby', 'A novel by F. Scott Fitzgerald', 'concrete'),
    ('Photosynthesis', 'How plants convert light into energy', 'abstract');

SELECT * FROM things;