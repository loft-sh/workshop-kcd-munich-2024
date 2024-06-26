-- Create the database
CREATE DATABASE IF NOT EXISTS ecommerce;
USE ecommerce;

-- Drop tables if they already exist
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- Create products table
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

-- Create orders table
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    product_id INT,
    quantity INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Insert data into users table
INSERT INTO users (name, email) VALUES 
('John Doe', 'john@example.com'),
('Jane Smith', 'jane@example.com'),
('Alice Johnson', 'alice@example.com');

-- Insert data into products table
INSERT INTO products (name, price) VALUES 
('Laptop', 999.99),
('Smartphone', 499.99),
('Tablet', 299.99);

-- Insert data into orders table
INSERT INTO orders (user_id, product_id, quantity) VALUES 
(1, 1, 1),
(2, 2, 2),
(3, 3, 1),
(1, 3, 3);

-- Confirm data inserted successfully
SELECT * FROM users;
SELECT * FROM products;
SELECT * FROM orders;

