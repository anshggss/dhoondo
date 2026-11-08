-- name: InsertWebsite :exec
insert into websites (
    id,
    url
) values (
$1,
$2
);

