-- name: GetTicketByID :one
SELECT id, event_name, stadium, price, seat_id, status, event_date
FROM tickets
WHERE id = $1;

-- name: GetTicketBySeatAndEvent :one
SELECT id, event_name, stadium, price, seat_id, status, event_date
FROM tickets
WHERE seat_id = $1 AND event_name = $2;