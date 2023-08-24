begin ;

create unique index id_unique_address
    on restaurant(address)
    where archived_at is null;

end;