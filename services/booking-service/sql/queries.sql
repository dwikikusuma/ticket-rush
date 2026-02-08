-- name: CreateBooking :one
INSERT INTO bookings (user_id, ticket_id, status)
VALUES ($1, $2, $3)
    RETURNING id, user_id, ticket_id, status, created_at;

-- name: UpdateBookingStatus :exec
UPDATE bookings
SET status = $1
WHERE id = $2;