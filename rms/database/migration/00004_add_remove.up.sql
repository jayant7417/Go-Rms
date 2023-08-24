begin ;
alter table restaurant
    add column address text not null,
    add column lat float8 not null ,
    add column log float8 not null ;


alter table address
    drop column rid;

commit;