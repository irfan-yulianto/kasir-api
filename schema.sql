-- ============================================================================
-- SQL SCHEMA untuk Supabase
-- ============================================================================

-- Buat tabel categories
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Buat tabel produk dengan foreign key ke categories
CREATE TABLE IF NOT EXISTS produk (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    harga INTEGER NOT NULL DEFAULT 0,
    stok INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL
);

-- Insert data awal categories
INSERT INTO categories (name, description) VALUES 
    ('Makanan', 'Produk makanan dan snack'),
    ('Minuman', 'Produk minuman'),
    ('Bumbu Dapur', 'Produk bumbu masak');

-- Insert data awal produk dengan category_id
INSERT INTO produk (nama, harga, stok, category_id) VALUES 
    ('Indomie Goreng', 3500, 10, 1),
    ('Vit 1000ml', 3000, 40, 2),
    ('Kecap ABC', 12000, 20, 3);

-- ============================================================================
-- TRANSACTION TABLES (Session 3)
-- ============================================================================

-- Buat tabel transactions
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Buat tabel transaction_details
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES produk(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);
