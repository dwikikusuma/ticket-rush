-- name: GetTicketByID :one
SELECT id, event_name, stadium, price, seat_id, status, event_date
FROM tickets
WHERE id = $1;