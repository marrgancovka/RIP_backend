drop table if exists Users CASCADE;
drop table if exists Ships CASCADE;
drop table if exists Flights CASCADE;
drop table if exists FlightsShip CASCADE;
drop table if exists Cosmodroms CASCADE;

create table if not exists Users(
    id serial PRIMARY KEY,
    firstName varchar(20) not null,
    secondName varchar(30) not null,
    phone varchar(12) UNIQUE not null,
    username varchar(30) UNIQUE not null,
    userpassword varchar(30) not null,
    userrole varchar(30) not null
);

create table if not exists Ships(
    id serial PRIMARY KEY,
    title varchar(30) not null,
    rocket varchar(50) not null,
    type varchar(50) not null,
    description text not null,
    image_url varchar(200) not null,
    is_delete boolean no null
);

create table if not exists Flights(
    id serial PRIMARY KEY,
    status varchar(30) not null,
    date_creation date default now() not null,
    date_formation date not null,
    date_end date not null,
    date_flight date not null,
    id_user int(10) not null REFERENCES Users (id)
);

create table if not exists Cosmodroms(
    id serial PRIMARY KEY,
    title varchar(40) not null UNIQUE,
    city varchar(50),
    country varchar(50)
);

create table if not exists FlightsShip(
    id serial PRIMARY KEY,
    id_ship int(10) not null REFERENCES Ships (id),
    id_flight int(10) not null REFERENCES Flights (id),
    id_cosmodrom_begin int(5) REFERENCES Cosmodroms (id) not null,
    id_cosmodrom_end int(5) REFERENCES Cosmodroms (id) not null
);

