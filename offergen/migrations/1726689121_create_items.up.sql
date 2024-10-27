CREATE TABLE IF NOT EXISTS items(
    id uuid PRIMARY KEY,
    owner_id uuid references users(id),
    name VARCHAR (500) NOT NULL,
    price NUMERIC CHECK (price > 0)
);

