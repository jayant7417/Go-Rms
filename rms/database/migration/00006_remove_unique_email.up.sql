begin ;
alter table users
drop constraint users_email_key;
end;