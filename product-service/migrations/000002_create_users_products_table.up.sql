CREATE TABLE IF NOT EXISTS users_products (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL ON DELETE CASCADE,
    product_id INT NOT NULL ON DELETE CASCADE,
    amount INT NOT NULL DEFAULT 1
);
