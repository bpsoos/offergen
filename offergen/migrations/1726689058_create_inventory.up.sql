CREATE TABLE IF NOT EXISTS inventories(
    owner_id uuid references users(id) PRIMARY KEY,
    title VARCHAR (500) NOT NULL,
    is_published BOOLEAN NOT NULL
);

