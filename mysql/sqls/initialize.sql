DROP DATABASE IF EXISTS chatdb;

CREATE DATABASE chatdb;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

USE chatdb;

CREATE TABLE `usrs` (
        `uuid` VARCHAR(255) NOT NULL PRIMARY KEY,
        `id` VARCHAR(20) NOT NULL, 
        `password` VARCHAR(255) NOT NULL, 
        `conn_id` INT NOT NULL,
        `order_usr` INT(2) DEFAULT 0);

CREATE TABLE `chat` (
        `chat_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `writer_id` VARCHAR(255) NOT NULL,
        `write_time` DATETIME NOT NULL,
        `text_body` TEXT NOT NULL,
        `is_answer` TINYINT(1) DEFAULT 0);

CREATE TABLE `request` (
        `request_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `requester_uuid` VARCHAR(255) NOT NULL,
        `requester_id` VARCHAR(255) NOT NULL,
        `target_uuid` VARCHAR(255) NOT NULL,
        `target_id` VARCHAR(255) NOT NULL,
        `request_time` VARCHAR(100) NOT NULL);

CREATE TABLE `connection` (
        `connection_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `first_usr` VARCHAR(255) NOT NULL,
        `second_usr` VARCHAR(255) NOT NULL,
        `start_date` VARCHAR(100) NOT NULL);

CREATE TABLE `beabouttodelete` (
        `delete_date` DATETIME DEFAULT '0000-00-00 00:00:00',
        `connection_id` INT,
        FOREIGN KEY (`connection_id`) REFERENCES `connection`(`connection_id`) ON DELETE CASCADE);

CREATE TABLE `question` (
        `question_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `target_word` VARCHAR(255) NOT NULL,
        `question_contents` VARCHAR(255) NOT NULL);

CREATE TABLE `answer` (
        `answer_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `connection_id` INT NOT NULL,
        `first_answer` VARCHAR(255) DEFAULT 'not-written',
        `second_answer` VARCHAR(255) DEFAULT 'not-written',
        `answer_date` VARCHAR(255) NOT NULL,
        `question_id` INT,
        FOREIGN KEY (`question_id`) REFERENCES `question`(`question_id`) ON UPDATE CASCADE ON DELETE CASCADE);

CREATE TABLE `exceptionword` (
        `exception_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `connection_id` INT NOT NULL,
        `except_word` TEXT NOT NULL);

CREATE TABLE `anniversary` (
        `anniversary_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `connection_id` INT NOT NULL,
        `year` INT NOT NULL,
        `month` INT NOT NULL,
        `date` INT NOT NULL,
        `contents` VARCHAR(255) NOT NULL,
        `d_day` TINYINT(1) NOT NULL);        

-- 이거 작동 안함 왜 그런겨? chat DB create까지만 작동함

-- https://devpress.csdn.net/cloudnative/63055e53c67703293080f68c.html
-- 이거 보고 mysql-server 커넥션 공부 더 하기, 근데 이거도 정답은 아니고 누가 질문 올린거임

-- mysql 코드 수정하면 mysql_data 폴더 지우고 다시 docker-compose 빌드해야함
-- volumes 때문에 저걸 참조해서 db가 작동하기 때문