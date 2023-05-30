CREATE DATABASE movies;

CREATE TABLE movies (
    id              INT AUTO_INCREMENT NOT NULL,
    isbn            VARCHAR(128) NOT NULL,
    title           VARCHAR(128) NOT NULL,
    director        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=INNODB;

insert into movies(isbn, title, director) values('9788401490040', 'Jaws', 'Steven Spielberg');