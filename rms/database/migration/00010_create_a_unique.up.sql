begin ;

create unique index address_unique
    on restaurant(address)
    where archived_at is null;

end;