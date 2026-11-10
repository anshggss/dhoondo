-- name: FindWord :one
select id
from words
where word = $1;
