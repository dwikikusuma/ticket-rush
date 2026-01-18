CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    event_name VARCHAR(100),
    stadium VARCHAR(100),
    price INT,
    seat_id VARCHAR(50),
    status VARCHAR(20) DEFAULT 'AVAILABLE',
    event_date TIMESTAMP
);

-- Add an index for the search optimization (we will see why later)
CREATE INDEX idx_tickets_status ON tickets(status);
