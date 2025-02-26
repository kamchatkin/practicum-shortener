-- Write your migrate up statements here

alter table aliases
    add column deleted_at timestamp;

---- create above / drop below ----

alter table aliases
    drop column deleted_at;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
