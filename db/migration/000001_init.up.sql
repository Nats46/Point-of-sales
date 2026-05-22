CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    spv BIGINT REFERENCES users(id) ON DELETE SET NULL,
    "group" VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS approvals (
    id BIGSERIAL PRIMARY KEY,
    transaction VARCHAR(255) NOT NULL,
    requester BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    approver BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS inventories (
    id BIGSERIAL PRIMARY KEY,
    item_code VARCHAR(100) NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(100) NOT NULL,
    updated_by VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS batches (
    id BIGSERIAL PRIMARY KEY,
    item_code VARCHAR(100) NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(100) NOT NULL,
    updated_by VARCHAR(100) NOT NULL,
    stock DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    batch_number VARCHAR(100) NOT NULL,
    expired_date TIMESTAMP NOT NULL,
    stock_date TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS discounts (
    id BIGSERIAL PRIMARY KEY,
    discount_code VARCHAR(100) UNIQUE NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    discount_type VARCHAR(50) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(100) NOT NULL,
    updated_by VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS sales_headers (
    transaction_id BIGSERIAL PRIMARY KEY,
    transaction_no VARCHAR(100) UNIQUE NOT NULL,
    customer_id BIGINT,
    cashier_id BIGINT NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) NOT NULL DEFAULT 0,
    tax DECIMAL(10,2) NOT NULL DEFAULT 0,
    grand_total DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(100) NOT NULL,
    updated_by VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS sales_details (
    sales_detail_id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES sales_headers(transaction_id) ON DELETE CASCADE,
    item_id BIGINT NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    item_code VARCHAR(100) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) NOT NULL DEFAULT 0,
    tax DECIMAL(10,2) NOT NULL DEFAULT 0,
    grand_total DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(100) NOT NULL,
    updated_by VARCHAR(100) NOT NULL
);
