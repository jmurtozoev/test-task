CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) UNIQUE,
    price numeric,
    update_count integer DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE EXTENSION pg_trgm;

INSERT INTO products(name, price) VALUES('Coca-cola', 60),
                                         ('Bread', 15),
                                         ('Water', 10),
                                         ('Chocolate', 100);