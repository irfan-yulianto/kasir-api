-- ============================================================================
-- SQL SCHEMA untuk Supabase (English naming convention)
-- ============================================================================

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Create products table with foreign key to categories
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL
);

-- Insert initial categories data
INSERT INTO categories (name, description) VALUES 
    ('Makanan', 'Produk makanan dan snack'),
    ('Minuman', 'Produk minuman'),
    ('Bumbu Dapur', 'Produk bumbu masak');

-- Insert initial products data with category_id
INSERT INTO products (name, price, stock, category_id) VALUES 
    ('Indomie Goreng', 3500, 10, 1),
    ('Vit 1000ml', 3000, 40, 2),
    ('Kecap ABC', 12000, 20, 3);

-- ============================================================================
-- TRANSACTION TABLES (Session 3)
-- ============================================================================

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create transaction_details table
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);

-- ============================================================================
-- MIGRATION: Rename produk to products (run this if you have existing data)
-- ============================================================================
-- ALTER TABLE produk RENAME TO products;
-- ALTER TABLE products RENAME COLUMN nama TO name;
-- ALTER TABLE products RENAME COLUMN harga TO price;
-- ALTER TABLE products RENAME COLUMN stok TO stock;
