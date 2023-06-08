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

-- name: GetUserByID :one
select username, avatar, created_at, updated_at
from users
where id = $1;

-- name: UpdateUser :exec
update users as u
set username   = coalesce(nullif($1, ''), u.username),
    password   = coalesce(nullif($2, ''), u.password),
    avatar     = coalesce(nullif($3, ''), u.avatar),
    updated_at = timezone('utc', now())
where id = $4;

-- name: DeleteUser :exec
delete
from users
where id = $1;