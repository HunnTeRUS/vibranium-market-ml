CREATE TABLE IF NOT EXISTS orders (
    orderId VARCHAR(255) PRIMARY KEY UNIQUE,
    userId VARCHAR(255) NOT NULL,
    type INT NOT NULL,
    amount INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    status VARCHAR(255) NOT NULL,
    symbol VARCHAR(255) NOT NULL,
    INDEX idx_orders_symbol_status_type (symbol, status, type),
    INDEX idx_orders_orderId (orderId)
);

CREATE TABLE IF NOT EXISTS wallet (
    userId VARCHAR(255) PRIMARY KEY UNIQUE,
    balance DECIMAL(10,2) NOT NULL DEFAULT 0,
    vibranium INT DEFAULT 0 NOT NULL
);