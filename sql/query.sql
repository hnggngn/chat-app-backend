-- name: GetUserByUsername :one
select *
from users
where username = $1;

-- name: CreateNewUser :exec
insert into users (username, password, avatar)
values ($1, $2, $3);

-- name: DeleteUserByUsername :exec
delete
from users
where username = $1;