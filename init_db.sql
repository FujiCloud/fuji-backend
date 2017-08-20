CREATE DATABASE fuji;
USE fuji;

CREATE TABLE events(
    id int not null auto_increment primary key,
    name text not null,
    attributes text not null,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
