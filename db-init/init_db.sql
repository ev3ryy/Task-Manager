CREATE TABLE Tasks(
    id SERIAl PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE
)