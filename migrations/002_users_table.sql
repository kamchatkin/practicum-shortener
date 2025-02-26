-- Write your migrate up statements here

create table if not exists cookie_users
(
    id serial8 primary key,
    created_at timestamp default now()
);

comment on table aliases is 'cookie users';

alter table aliases
    add user_id bigint default 0;

---- create above / drop below ----

drop table cookie_users;
alter table aliases
    drop column user_id;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
