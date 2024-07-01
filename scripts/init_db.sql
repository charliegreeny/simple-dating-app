CREATE DATABASE dating_users;

USE dating_users;

CREATE TABLE genders(
    gender VARCHAR(128),
    PRIMARY KEY (gender)
);

CREATE TABLE users(
    id varchar(256),
    name varchar(256) NOT NULL,
    email varchar(256) NOT NULL,
    date_of_birth date NOT NULL,
    gender VARCHAR(128) NOT NULL,
    password  varchar(256) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (email),
    CONSTRAINT fk_userGender FOREIGN KEY (gender) REFERENCES genders(gender)
);

INSERT INTO genders VALUES ('MALE');
INSERT INTO genders VALUES ('FEMALE');