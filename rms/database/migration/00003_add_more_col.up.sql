begin ;
alter table address
    add column lat float8,
    add column log float8;


alter table users
drop column created_by;

commit;