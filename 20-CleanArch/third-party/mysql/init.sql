-- Initial database setup for Order System
-- This script will be executed when MySQL container starts for the first time

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    price DECIMAL(10,2) NOT NULL,
    tax DECIMAL(10,2) NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_final_price ON orders(final_price);

-- Insert some sample data for testing
INSERT INTO orders (id, price, tax, final_price) VALUES 
('order-001', 100.00, 10.00, 110.00),
('order-002', 250.50, 25.05, 275.55),
('order-003', 75.25, 7.53, 82.78)
ON DUPLICATE KEY UPDATE 
    price = VALUES(price),
    tax = VALUES(tax),
    final_price = VALUES(final_price); 