begin;

drop index if exists id_user_active;

create unique index id_unique_user
on users(email)
where archived_at is not null;



end;