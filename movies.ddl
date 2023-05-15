CREATE DATABASE movies;

CREATE TABLE directors (
    id          INT AUTO_INCREMENT NOT NULL,
    first_name  VARCHAR(128) NOT NULL,
    last_name   VARCHAR(128) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=INNODB;

CREATE TABLE movies (
    id              INT AUTO_INCREMENT NOT NULL,
    isbn            VARCHAR(128) NOT NULL,
    title           VARCHAR(128) NOT NULL,
    director_id     INT,
    PRIMARY KEY (id),
    INDEX dir_ind (director_id),
    FOREIGN KEY (director_id)
        REFERENCES directors(id)
        ON DELETE CASCADE
) ENGINE=INNODB;

insert into directors(first_name, last_name) values ('Steven', 'Spielberg');
insert into movies(isbn, title, director_id) values('9788401490040', 'Jaws', 1);