begin;

drop index if exists id_unique_user;

create unique index id_unique_users
    on users(email)
    where archived_at is null;



end;