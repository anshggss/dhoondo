-- name: IndexWord :exec
insert into reverseindex (
    wordid,
    websiteid
) values (
    $1,
    $2
);
