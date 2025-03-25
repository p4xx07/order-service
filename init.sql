-- Use the correct database
USE test;

-- Create a sample users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inserting sample user data
INSERT INTO users (name, email) VALUES
('John Doe', 'john.doe@example.com' ),
('Jane Smith', 'jane.smith@example.com'),
('Emily Johnson', 'emily.johnson@example.com'),
('Michael Brown', 'michael.brown@example.com'),
('Sarah Davis', 'sarah.davis@example.com');

-- Creating the products table
CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(50)
);

-- Inserting sample product data
INSERT INTO products (name, description, price, category) VALUES
('Laptop', 'High-performance laptop for gaming and development', 1200.00, 'Electronics'),
('Smartphone', 'Latest model with powerful camera and fast performance', 800.00, 'Electronics'),
('Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 25.50, 'Accessories'),
('Headphones', 'Noise-canceling headphones for immersive sound', 150.00, 'Accessories'),
('Keyboard', 'Mechanical keyboard with RGB lighting', 100.00, 'Accessories'),
('Coffee Mug', 'Ceramic mug with a funny quote', 15.00, 'Home & Kitchen'),
('Blender', 'High-speed blender for smoothies and shakes', 60.00, 'Home & Kitchen'),
('Desk Chair', 'Comfortable ergonomic chair for home office', 200.00, 'Furniture');

-- Creating the inventory table
CREATE TABLE IF NOT EXISTS inventories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT UNSIGNED,
    stock BIGINT NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Inserting sample inventory data
INSERT INTO inventories (product_id, stock) VALUES
(1, 50),  -- Laptop
(2, 100), -- Smartphone
(3, 200), -- Wireless Mouse
(4, 150), -- Headphones
(5, 120), -- Keyboard
(6, 300), -- Coffee Mug
(7, 80),  -- Blender
(8, 60);  -- Desk Chair

-- Creating the orders table
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inserting sample order data
INSERT INTO orders (user_id, status) VALUES
(1, 'pending'),
(2, 'completed'),
(3, 'shipped'),
(4, 'cancelled'),
(5, 'pending');

-- Creating the order_items table
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED,
    product_id BIGINT UNSIGNED,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Inserting sample order items data
INSERT INTO order_items (order_id, product_id, quantity, price) VALUES
(1, 1, 1, 1200.00),  -- Order 1, Laptop
(1, 3, 2, 25.50),    -- Order 1, Wireless Mouse
(2, 2, 1, 800.00),   -- Order 2, Smartphone
(2, 4, 1, 150.00),   -- Order 2, Headphones
(3, 5, 1, 100.00),   -- Order 3, Keyboard
(3, 6, 3, 15.00),    -- Order 3, Coffee Mug
(4, 7, 1, 60.00),    -- Order 4, Blender
(5, 8, 1, 200.00);   -- Order 5, Desk Chair
