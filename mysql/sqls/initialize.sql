DROP DATABASE IF EXISTS chat;

CREATE DATABASE chat;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

DROP TABLE IF EXISTS `usrs`;

CREATE TABLE `usrs` (`id` int NOT NULL);

-- https://devpress.csdn.net/cloudnative/63055e53c67703293080f68c.html
-- 이거 보고 mysql-server 커넥션 공부 더 하기, 근데 이거도 정답은 아니고 누가 질문 올린거임
