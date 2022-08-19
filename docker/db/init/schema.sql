-- functions

DELIMITER $$
CREATE FUNCTION RandString(length SMALLINT) RETURNS varchar(100) CHARSET utf8mb4
begin
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

-- tables

CREATE TABLE IF NOT EXISTS games (
    id varchar(30) NOT NULL,
    name varchar(200) NOT NULL,
    version int NOT NULL DEFAULT '1',
    version_name varchar(50) NOT NULL DEFAULT '1.0.0',
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS levels (
    gid varchar(30) NOT NULL,
    number int NOT NULL,
    data json NOT NULL,
    min_valid_time int NOT NULL DEFAULT '0',
    max_valid_score int NOT NULL DEFAULT '0',
    PRIMARY KEY (gid,number)
);

CREATE TABLE IF NOT EXISTS players (
    id varchar(30) NOT NULL,
    name varchar(200) NOT NULL,
    PRIMARY KEY (id)
);

DELIMITER $$
CREATE TRIGGER players_before_insert BEFORE INSERT ON players FOR EACH ROW BEGIN
SET @pid = "";
WHILE (@pid IS NOT NULL) DO 
    SET NEW.id = RANDSTRING(30);
    SET @pid = (SELECT id FROM players WHERE id = NEW.id);
END WHILE;
END
$$
DELIMITER ;

CREATE TABLE IF NOT EXISTS scores (
    gid varchar(30) NOT NULL,
    level int NOT NULL,
    pid varchar(30) NOT NULL,
    time int NOT NULL DEFAULT '0',
    score int NOT NULL DEFAULT '0',
    PRIMARY KEY (gid,level,pid),
    KEY pid (pid)
);

-- foreign keys

ALTER TABLE levels
    ADD CONSTRAINT levels_ibfk_1 FOREIGN KEY (gid) REFERENCES games (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE scores
	ADD CONSTRAINT scores_ibfk_1 FOREIGN KEY (gid) REFERENCES games (id) ON DELETE CASCADE ON UPDATE CASCADE,
	ADD CONSTRAINT scores_ibfk_2 FOREIGN KEY (pid) REFERENCES players (id) ON DELETE CASCADE ON UPDATE CASCADE;