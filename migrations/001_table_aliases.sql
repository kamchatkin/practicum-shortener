-- Write your migrate up statements here

create table if not exists aliases
(
    alias      varchar(10) primary key,
    source     text not null,
    quantity   bigint    default 0,
    created_at timestamp default now()
);

comment on table aliases is 'long to short and vice versa';

create unique index aliases_source_uindex
    on aliases (source);

comment on column aliases.quantity is 'redirects';


---- create above / drop below ----

drop table aliases;
