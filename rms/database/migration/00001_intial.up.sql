begin ;

create table if not exists users(
                                    uid serial primary key ,
                                    role text default 'user' not null ,
                                    created_by int ,
                                    name text not null,
                                    email text not null ,
                                    password text not null,
                                    created_at timestamp with time zone default now(),
                                    updated_at timestamp with time zone default now(),
                                    archived_at timestamp with time zone default null
);

create table if not exists restaurant(
                                         rid serial primary key,
                                         name text not null,
                                         created_by int not null,
                                         created_at timestamp with time zone default now(),
                                         updated_at timestamp with time zone default now(),
                                         archived_at timestamp with time zone default null
);

create table if not exists address(
                                      address_id serial primary key ,
                                      uid int references users(uid),
                                      rid int references restaurant(rid),
                                      address text not null ,
                                      coordinates POINT NOT NULL,
                                      created_at timestamp with time zone default now(),
                                      updated_at timestamp with time zone default now(),
                                      archived_at timestamp with time zone default null

);

create table if not exists dishes(
                                     did serial primary key ,
                                     rid int references restaurant(rid),
                                     name text not null ,
                                     rate int not null,
                                     created_at timestamp with time zone default now(),
                                     updated_at timestamp with time zone default now(),
                                     archived_at timestamp with time zone default null
);

commit ;