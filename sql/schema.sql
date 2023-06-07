create table users
(
    id         uuid primary key         default gen_random_uuid()      not null,
    username   varchar(30) unique                                      not null,
    password   varchar(100)                                            not null,
    avatar     varchar(254),
    created_at timestamp with time zone default timezone('utc', now()) not null,
    updated_at timestamp with time zone default timezone('utc', now()) not null
);