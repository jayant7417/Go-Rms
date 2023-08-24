begin ;
alter table address
    add column created_at timestamp with time zone default now(),
    add column updated_at timestamp with time zone default now(),
    add column archived_at timestamp with time zone;




alter table dishes
    add column created_at timestamp with time zone default now(),
    add column updated_at timestamp with time zone default now(),
    add column archived_at timestamp with time zone;


commit;