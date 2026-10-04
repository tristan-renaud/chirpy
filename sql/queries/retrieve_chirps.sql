-- name: RetrieveChirps :one
SELECT * FROM chirps
ORDER BY created_at ASC;
