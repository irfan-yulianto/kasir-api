CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);

INSERT INTO categories (name, description) VALUES 
    ('Makanan', 'Produk makanan dan snack'),
    ('Minuman', 'Produk minuman'),
    ('Bumbu Dapur', 'Produk bumbu masak');

INSERT INTO products (name, price, stock, category_id) VALUES 
    ('Indomie Goreng', 3500, 10, 1),
    ('Vit 1000ml', 3000, 40, 2),
    ('Kecap ABC', 12000, 20, 3);
