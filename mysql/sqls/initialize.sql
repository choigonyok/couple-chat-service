DROP DATABASE IF EXISTS chat;

CREATE DATABASE chat;

USE chat;

CREATE TABLE usrs (
        usr_uniqueid int,
        usrname varchar(255),
        usr_id varchar(255),
        usr_pw varchar(255),
        phone_num int,
        create_date datetime,
        is_conneted tinyint(1)
);

