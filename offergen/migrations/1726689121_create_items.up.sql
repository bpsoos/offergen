CREATE TABLE IF NOT EXISTS items(
    id uuid PRIMARY KEY,
    owner_id uuid references users(id) ON DELETE CASCADE,
    name VARCHAR (150) NOT NULL,
    description VARCHAR (500),
    category VARCHAR (50),
    price NUMERIC CHECK (price > 0)
);

