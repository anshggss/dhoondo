-- name: FindWebsite :one
Select count(*) from websites where url = $1;
