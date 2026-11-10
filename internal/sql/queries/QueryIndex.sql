-- name: QueryIndex :many
select url 
from websites 
where id
in 
(   select websiteid 
    from reverseindex 
    where wordid = 
    (   select id 
        from words 
        where word = $1));
