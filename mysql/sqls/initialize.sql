DROP DATABASE IF EXISTS chatdb;

CREATE DATABASE chatdb;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

USE chatdb;

DROP TABLE IF EXISTS `usrs`;

CREATE TABLE `usrs` (
        `id` VARCHAR(20) PRIMARY KEY, 
        `password` VARCHAR(255) NOT NULL, 
        `conn_id` VARCHAR(255));

DROP TABLE IF EXISTS `chat`;

CREATE TABLE `chat` (
        `chat_id` INT AUTO_INCREMENT PRIMARY KEY,
        `conn_id` VARCHAR(255) NOT NULL,
        `writer_id` VARCHAR(255) NOT NULL,
        `write_time` VARCHAR(100) NOT NULL,
        `text_body` TEXT NOT NULL);


-- 이거 작동 안함 왜 그런겨? chat DB create까지만 작동함

-- https://devpress.csdn.net/cloudnative/63055e53c67703293080f68c.html
-- 이거 보고 mysql-server 커넥션 공부 더 하기, 근데 이거도 정답은 아니고 누가 질문 올린거임

-- mysql 코드 수정하면 mysql_data 폴더 지우고 다시 docker-compose 빌드해야함
-- volumes 때문에 저걸 참조해서 db가 작동하기 때문