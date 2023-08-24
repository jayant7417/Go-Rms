begin ;
create index id_user_active on users(email)
where archived_at is not null;
end;