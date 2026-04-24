-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price numeric(12, 2) NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

	create_at TIMESTAMP DEFAULT NOW(),
	update_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS subscriptions;