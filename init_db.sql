CREATE DATABASE fuji;
USE fuji;

CREATE TABLE events(
    id int not null auto_increment primary key,
    name tinytext not null,
    attributes text not null,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users(
    id int not null auto_increment primary key,
    os tinytext,
    device tinytext,
    locale tinytext,
    voiceover boolean,
    bold_text boolean,
    reduce_motion boolean,
    reduce_transparency boolean,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
