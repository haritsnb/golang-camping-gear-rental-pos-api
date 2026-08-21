-- +migrate Up
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL REFERENCES roles(id),
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

ALTER TABLE roles ADD CONSTRAINT fk_roles_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE roles ADD CONSTRAINT fk_roles_modified_by FOREIGN KEY (modified_by) REFERENCES users(id);
ALTER TABLE roles ADD CONSTRAINT fk_roles_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    identity_type VARCHAR(20) NOT NULL,
    identity_number VARCHAR(50) NOT NULL UNIQUE,
    identity_photo_url VARCHAR(255),
    phone VARCHAR(20) NOT NULL,
    emergency_contact VARCHAR(20),
    email VARCHAR(100),
    address TEXT NOT NULL,
    is_blacklisted BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id),
    brand_id INT NOT NULL REFERENCES brands(id),
    name VARCHAR(150) NOT NULL,
    rental_price_per_day DECIMAL(12,2) NOT NULL,
    default_deposit DECIMAL(12,2) DEFAULT 0,
    late_fee_per_day DECIMAL(12,2) NOT NULL,
    lost_compensation_fee DECIMAL(12,2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE product_units (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    unit_code VARCHAR(50) NOT NULL UNIQUE,
    serial_number VARCHAR(100),
    condition VARCHAR(30) DEFAULT 'good',
    status VARCHAR(30) DEFAULT 'available',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id INT NOT NULL REFERENCES customers(id),
    user_id INT NOT NULL REFERENCES users(id),
    booking_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    start_date TIMESTAMP NOT NULL,
    expected_return_date TIMESTAMP NOT NULL,
    actual_return_date TIMESTAMP,
    total_rental_fee DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_deposit DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_penalty_fee DECIMAL(12,2) NOT NULL DEFAULT 0,
    grand_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(30) DEFAULT 'booked',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE rental_items (
    id SERIAL PRIMARY KEY,
    rental_id INT NOT NULL REFERENCES rentals(id) ON DELETE RESTRICT,
    product_unit_id INT NOT NULL REFERENCES product_units(id),
    daily_rate DECIMAL(12,2) NOT NULL,
    duration_days INT NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    condition_out VARCHAR(50) DEFAULT 'good',
    condition_in VARCHAR(50),
    item_penalty_fee DECIMAL(12,2) DEFAULT 0,
    return_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    rental_id INT NOT NULL REFERENCES rentals(id),
    user_id INT NOT NULL REFERENCES users(id),
    payment_type VARCHAR(30) NOT NULL,
    payment_method VARCHAR(30) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    reference_number VARCHAR(100),
    notes TEXT,
    paid_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE maintenance (
    id SERIAL PRIMARY KEY,
    product_unit_id INT NOT NULL REFERENCES product_units(id),
    issue_description TEXT NOT NULL,
    cost DECIMAL(12,2) DEFAULT 0,
    start_date DATE NOT NULL,
    completed_date DATE,
    status VARCHAR(30) DEFAULT 'in_progress',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_by INT REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id)
);

CREATE TABLE revoked_tokens (
    id SERIAL PRIMARY KEY,
    token_jti VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexing
CREATE INDEX idx_customers_deleted ON customers(deleted_at);
CREATE INDEX idx_products_deleted ON products(deleted_at);
CREATE INDEX idx_product_units_deleted ON product_units(deleted_at);
CREATE INDEX idx_rentals_deleted ON rentals(deleted_at);
CREATE INDEX idx_product_units_code ON product_units(unit_code);
CREATE INDEX idx_product_units_status ON product_units(status);
CREATE INDEX idx_rentals_invoice ON rentals(invoice_number);

-- +migrate Down
DROP TABLE IF EXISTS revoked_tokens CASCADE;
DROP TABLE IF EXISTS maintenance CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS rental_items CASCADE;
DROP TABLE IF EXISTS rentals CASCADE;
DROP TABLE IF EXISTS product_units CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS brands CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
