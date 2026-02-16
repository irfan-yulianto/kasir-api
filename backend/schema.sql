CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL DEFAULT 0 CHECK (price >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
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
    quantity INT NOT NULL CHECK (quantity > 0),
    subtotal INT NOT NULL CHECK (subtotal >= 0)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'kasir' CHECK (role IN ('admin', 'kasir')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_td_transaction_id ON transaction_details(transaction_id);
CREATE INDEX IF NOT EXISTS idx_td_product_id ON transaction_details(product_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Seed data
INSERT INTO categories (name, description) VALUES
    ('Makanan', 'Produk makanan dan snack'),
    ('Minuman', 'Produk minuman'),
    ('Bumbu Dapur', 'Produk bumbu masak');

INSERT INTO products (name, price, stock, category_id) VALUES
    ('Indomie Goreng', 3500, 10, 1),
    ('Vit 1000ml', 3000, 40, 2),
    ('Kecap ABC', 12000, 20, 3);
