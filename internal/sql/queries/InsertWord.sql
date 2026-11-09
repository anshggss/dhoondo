-- name: InsertWord :one
insert into words(
    word
) values (
$1
) returning id;

