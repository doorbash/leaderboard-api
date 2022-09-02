-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Sep 02, 2022 at 02:06 PM
-- Server version: 8.0.30
-- PHP Version: 8.0.19

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `leaderboard`
--

DELIMITER $$
--
-- Procedures
--
CREATE DEFINER=`root`@`%` PROCEDURE `GET_LEADERBOARD` (IN `p_lid` VARCHAR(50) CHARSET utf8mb4, IN `p_offset` INT, IN `p_count` INT)   BEGIN
DECLARE _lid INT DEFAULT 0;
DECLARE v1_order TINYINT DEFAULT 0;
DECLARE v2_order TINYINT DEFAULT 0;
DECLARE v3_order TINYINT DEFAULT 0;

START TRANSACTION;

SELECT `id`, `value1_order`, `value2_order`, `value3_order`
INTO _lid, v1_order, v2_order, v3_order
FROM `leaderboards`
WHERE `uid` = p_lid;

IF _lid > 0 THEN

SELECT name, value1, value2, value3 FROM `leaderboard_data`
JOIN players ON `leaderboard_data`.`pid` = `players`.`uid`
WHERE `players`.`banned` = 0 AND `lid` = p_lid
ORDER BY

CASE WHEN v3_order = 1 THEN value3 END DESC,
CASE WHEN v3_order <> 1 THEN value3 END ASC

,

CASE WHEN v2_order = 1 THEN value2 END DESC,
CASE WHEN v2_order <> 1 THEN value2 END ASC

,

CASE WHEN v1_order = 1 THEN value1 END DESC,
CASE WHEN v1_order <> 1 THEN value1 END ASC,

`leaderboard_data`.id ASC

LIMIT p_offset, p_count;

END IF;

COMMIT;
END$$

CREATE DEFINER=`root`@`%` PROCEDURE `GET_PLAYER_POSITION` (IN `p_lid` VARCHAR(50) CHARSET utf8mb4, IN `p_pid` VARCHAR(30) CHARSET utf8mb4)   sp: BEGIN
DECLARE _lid INT DEFAULT 0;
DECLARE _ldid INT DEFAULT 0;
DECLARE v1_order TINYINT DEFAULT 0;
DECLARE v2_order TINYINT DEFAULT 0;
DECLARE v3_order TINYINT DEFAULT 0;
DECLARE v1 INT DEFAULT 0;
DECLARE v2 INT DEFAULT 0;
DECLARE v3 INT DEFAULT 0;

START TRANSACTION;

SELECT `id`, `value1`, `value2`, `value3`
INTO _ldid, v1, v2, v3
FROM `leaderboard_data`
WHERE `pid` = p_pid;

IF _ldid > 0 THEN

SELECT `id`, `value1_order`, `value2_order`, `value3_order`
INTO _lid, v1_order, v2_order, v3_order
FROM `leaderboards`
WHERE `uid` = p_lid;

IF _lid > 0 THEN

IF v3_order = 0 AND v2_order = 0 AND v1_order = 0 THEN
LEAVE sp;
END IF;

IF v3_order = 0 AND v2_order = 0 AND v1_order = 1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value1 > v1 OR (value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 0 AND v2_order = 0 AND v1_order = -1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value1 < v1 OR (value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 0 AND v2_order = 1 AND v1_order = 1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value2 > v2 OR (value2 = v2 AND value1 > v1) OR (value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 0 AND v2_order = 1 AND v1_order = -1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value2 > v2 OR (value2 = v2 AND value1 < v1) OR (value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 0 AND v2_order = -1 AND v1_order = 1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value2 < v2 OR (value2 = v2 AND value1 > v1) OR (value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 0 AND v2_order = -1 AND v1_order = -1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value2 < v2 OR (value2 = v2 AND value1 < v1) OR (value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 1 AND v2_order = 1 AND v1_order = 1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 > v3 OR (value3 = v3 AND value2 > v2) OR (value3 = v3 AND value2 = v2 AND value1 > v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 1 AND v2_order = 1 AND v1_order = -1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 > v3 OR (value3 = v3 AND value2 > v2) OR (value3 = v3 AND value2 = v2 AND value1 < v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 1 AND v2_order = -1 AND v1_order = 1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 > v3 OR (value3 = v3 AND value2 < v2) OR (value3 = v3 AND value2 = v2 AND value1 > v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = 1 AND v2_order = -1 AND v1_order = -1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 > v3 OR (value3 = v3 AND value2 < v2) OR (value3 = v3 AND value2 = v2 AND value1 < v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = -1 AND v2_order = 1 AND v1_order = 1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 < v3 OR (value3 = v3 AND value2 > v2) OR (value3 = v3 AND value2 = v2 AND value1 > v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = -1 AND v2_order = 1 AND v1_order = -1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 < v3 OR (value3 = v3 AND value2 > v2) OR (value3 = v3 AND value2 = v2 AND value1 < v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = -1 AND v2_order = -1 AND v1_order = 1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 < v3 OR (value3 = v3 AND value2 < v2) OR (value3 = v3 AND value2 = v2 AND value1 > v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

IF v3_order = -1 AND v2_order = -1 AND v1_order = -1 THEN
SELECT COUNT(id) FROM `leaderboard_data` WHERE value3 < v3 OR (value3 = v3 AND value2 < v2) OR (value3 = v3 AND value2 = v2 AND value1 < v1) OR (value3 = v3 AND value2 = v2 AND value1 = v1 AND id < _ldid);
LEAVE sp;
END IF;

END IF;

END IF;

COMMIT;
END$$

CREATE DEFINER=`root`@`%` PROCEDURE `NEW_LEADERBOARD_DATA` (IN `p_lid` VARCHAR(50) CHARSET utf8mb4, IN `p_pid` VARCHAR(30) CHARSET utf8mb4, IN `p_v1` INT ZEROFILL, IN `p_v2` INT ZEROFILL, IN `p_v3` INT ZEROFILL)   sp:BEGIN
DECLARE `ldid` INT DEFAULT 0;
DECLARE `v1` INT DEFAULT 0;
DECLARE `v2` INT DEFAULT 0;
DECLARE `v3` INT DEFAULT 0;
DECLARE `v1_order` TINYINT DEFAULT 0;
DECLARE `v2_order` TINYINT DEFAULT 0;
DECLARE `v3_order` TINYINT DEFAULT 0;
DECLARE `v1_limit` INT DEFAULT 0;
DECLARE `v2_limit` INT DEFAULT 0;
DECLARE `v3_limit` INT DEFAULT 0;

START TRANSACTION;

SELECT `value1_order`, `value2_order`, `value3_order`,
`value1_limit`, `value2_limit`, `value3_limit`
INTO v1_order, v2_order, v3_order, 
v1_limit, v2_limit, v3_limit
FROM `leaderboards`
WHERE `uid` = p_lid;


-- CHECK LIMITS

IF v3_limit <> 0 THEN
IF v3_order = 1 THEN
IF p_v3 > v3_limit THEN
SELECT 3;
LEAVE sp;
END IF;
ELSEIF v3_order = -1 THEN
IF p_v3 < v3_limit THEN
SELECT 3;
LEAVE sp;
END IF;
END IF;
END IF;

IF v2_limit <> 0 THEN
IF v2_order = 1 THEN
IF p_v2 > v2_limit THEN
SELECT 3;
LEAVE sp;
END IF;
ELSEIF v2_order = -1 THEN
IF p_v2 < v2_limit THEN
SELECT 3;
LEAVE sp;
END IF;
END IF;
END IF;

IF v1_limit <> 0 THEN
IF v1_order = 1 THEN
IF p_v1 > v1_limit THEN
SELECT 3;
LEAVE sp;
END IF;
ELSEIF v1_order = -1 THEN
IF p_v1 < v1_limit THEN
SELECT 3;
LEAVE sp;
END IF;
END IF;
END IF;



SELECT `id`, `value1`, `value2`, `value3`
INTO ldid, v1, v2, v3
FROM `leaderboard_data`
WHERE `lid` = p_lid AND `pid` = p_pid;

IF ldid > 0 THEN
-- EXISTS

IF p_v3 = v3 THEN

-- CHECK v2
IF p_v2 = v2 THEN

-- CHECK V1
IF p_v1 = v1 THEN
SELECT 2;
LEAVE sp;
ELSEIF v1_order = 1 THEN
IF p_v1 < v1 THEN
SELECT 2;
LEAVE sp;
END IF;
ELSEIF v1_order = -1 THEN
IF p_v1 > v1 THEN
SELECT 2;
LEAVE sp;
END IF;
END IF;

ELSEIF v2_order = 1 THEN
IF p_v2 < v2 THEN
SELECT 2;
LEAVE sp;
END IF;
ELSEIF v2_order = -1 THEN
IF p_v2 > v2 THEN
SELECT 2;
LEAVE sp;
END IF;
END IF;

ELSEIF v3_order = 1 THEN
IF p_v3 < v3 THEN
SELECT 2;
LEAVE sp;
END IF;
ELSEIF v3_order = -1 THEN
IF p_v3 > v3 THEN
SELECT 2;
LEAVE sp;
END IF;
END IF;

UPDATE `leaderboard_data`
SET `value1` = p_v1, `value2` = p_v2, `value3` = p_v3
WHERE `lid` = p_lid AND `pid` = p_pid;

SELECT 1;

ELSE
-- NOT EXISTS

INSERT INTO `leaderboard_data`(`lid`, `pid`, `value1`, `value2`, `value3`)
VALUES(p_lid, p_pid, p_v1, p_v2, p_v3);

SELECT 1;

END IF;
COMMIT;
END$$

CREATE DEFINER=`root`@`localhost` PROCEDURE `NEW_PLAYER` (IN `name` VARCHAR(200) CHARSET utf8mb4)   BEGIN
DECLARE lid INT DEFAULT 0;
START TRANSACTION;
INSERT INTO `players`(`name`) VALUES(name);
SET lid = LAST_INSERT_ID();
IF lid > 0 THEN
SELECT `id`, `uid` FROM `players` WHERE `id` = lid;
COMMIT;
ELSE
ROLLBACK;
END IF;
END$$

--
-- Functions
--
CREATE DEFINER=`root`@`localhost` FUNCTION `RANDSTRING` (`length` SMALLINT) RETURNS VARCHAR(100) CHARSET utf8mb4  begin
    SET @returnStr = '';
    SET @allowedChars = 'abcdefghijklmnopqrstuvwxyz0123456789';
    SET @i = 0;

    WHILE (@i < length) DO
        SET @returnStr = CONCAT(@returnStr, substring(@allowedChars, FLOOR(RAND() * LENGTH(@allowedChars) + 1), 1));
        SET @i = @i + 1;
    END WHILE;

    RETURN @returnStr;
END$$

DELIMITER ;

-- --------------------------------------------------------

--
-- Table structure for table `games`
--

CREATE TABLE IF NOT EXISTS `games` (
  `id` varchar(30) NOT NULL,
  `name` varchar(200) NOT NULL,
  `data` json NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `leaderboards`
--

CREATE TABLE IF NOT EXISTS `leaderboards` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` varchar(50) NOT NULL,
  `gid` varchar(30) NOT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'default',
  `value1_order` tinyint NOT NULL DEFAULT '1',
  `value2_order` tinyint NOT NULL DEFAULT '0',
  `value3_order` tinyint NOT NULL DEFAULT '0',
  `value1_limit` int NOT NULL DEFAULT '0',
  `value2_limit` int NOT NULL DEFAULT '0',
  `value3_limit` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`),
  KEY `gid` (`gid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Triggers `leaderboards`
--
DELIMITER $$
CREATE TRIGGER `leaderboards_before_insert` BEFORE INSERT ON `leaderboards` FOR EACH ROW BEGIN
SET @id = 0;
WHILE (@id IS NOT NULL) DO 
    SET NEW.`uid` = RANDSTRING(50);
    SET @id = (SELECT `id` FROM `leaderboards` WHERE `uid` = NEW.`uid`);
END WHILE;
END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Table structure for table `leaderboard_data`
--

CREATE TABLE IF NOT EXISTS `leaderboard_data` (
  `id` int NOT NULL AUTO_INCREMENT,
  `lid` varchar(50) NOT NULL,
  `pid` varchar(30) NOT NULL,
  `value1` int NOT NULL DEFAULT '0',
  `value2` int NOT NULL DEFAULT '0',
  `value3` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `lid` (`lid`,`pid`),
  KEY `pid` (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `players`
--

CREATE TABLE IF NOT EXISTS `players` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` varchar(30) NOT NULL,
  `name` varchar(200) NOT NULL,
  `banned` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Triggers `players`
--
DELIMITER $$
CREATE TRIGGER `players_before_insert` BEFORE INSERT ON `players` FOR EACH ROW BEGIN
SET @id = 0;
WHILE (@id IS NOT NULL) DO 
    SET NEW.`uid` = RANDSTRING(30);
    SET @id = (SELECT `id` FROM `players` WHERE `uid` = NEW.`uid`);
END WHILE;
END
$$
DELIMITER ;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `leaderboards`
--
ALTER TABLE `leaderboards`
  ADD CONSTRAINT `leaderboards_ibfk_1` FOREIGN KEY (`gid`) REFERENCES `games` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `leaderboard_data`
--
ALTER TABLE `leaderboard_data`
  ADD CONSTRAINT `leaderboard_data_ibfk_1` FOREIGN KEY (`lid`) REFERENCES `leaderboards` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `leaderboard_data_ibfk_2` FOREIGN KEY (`pid`) REFERENCES `players` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
