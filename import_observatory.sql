
CREATE DATABASE IF NOT EXISTS observatory;

CREATE TABLE IF NOT EXISTS `observatory`.`sector` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `coordinates` VARCHAR(255) NULL DEFAULT NULL,
  `light_intensity` FLOAT NULL DEFAULT NULL,
  `foreign_objects` INT NULL DEFAULT NULL,
  `star_objects` INT NULL DEFAULT NULL,
  `unknown_objects` INT NULL DEFAULT NULL,
  `known_objects` INT NULL DEFAULT NULL,
  `notes` TEXT NULL DEFAULT NULL,
  `date_update` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `observatory`.`position` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `earth_position` VARCHAR(255) NULL DEFAULT NULL,
  `sun_position` VARCHAR(255) NULL DEFAULT NULL,
  `moon_position` VARCHAR(255) NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `observatory`.`naturalobjects` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(255) NULL DEFAULT NULL,
  `galaxy` VARCHAR(255) NULL DEFAULT NULL,
  `accuracy` FLOAT NULL DEFAULT NULL,
  `flux` FLOAT NULL DEFAULT NULL,
  `related_objects` VARCHAR(255) NULL DEFAULT NULL,
  `notes` TEXT NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `observatory`.`objects` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(255) NULL DEFAULT NULL,
  `accuracy` FLOAT NULL DEFAULT NULL,
  `quantity` INT NULL DEFAULT NULL,
  `time` TIME NULL DEFAULT NULL,
  `date` DATE NULL DEFAULT NULL,
  `notes` TEXT NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `observatory`.`link` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `sector_id` INT NULL DEFAULT NULL,
  `objects_id` INT NULL DEFAULT NULL,
  `naturalobjects_id` INT NULL DEFAULT NULL,
  `position_id` INT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `sector_id_fk_idx` (`sector_id` ASC) VISIBLE,
  INDEX `object_id_fk_idx` (`objects_id` ASC) VISIBLE,
  INDEX `position_id_fk_idx` (`position_id` ASC) VISIBLE,
  INDEX `natural_object_id_fk_idx` (`naturalobjects_id` ASC) VISIBLE,
  CONSTRAINT `natural_object_id_fk`
    FOREIGN KEY (`naturalobjects_id`)
    REFERENCES `observatory`.`naturalobjects` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `object_id_fk`
    FOREIGN KEY (`objects_id`)
    REFERENCES `observatory`.`objects` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `position_id_fk`
    FOREIGN KEY (`position_id`)
    REFERENCES `observatory`.`position` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `sector_id_fk`
    FOREIGN KEY (`sector_id`)
    REFERENCES `observatory`.`sector` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;

COMMIT;

DELIMITER //

CREATE PROCEDURE `observatory`.`multi_select`(IN table1 VARCHAR(255), IN table2 VARCHAR(255))
BEGIN
    DECLARE col_list1 TEXT;
    DECLARE col_list2 TEXT;
    DECLARE combined_col_list TEXT;

    SELECT GROUP_CONCAT(
        CASE WHEN COLUMN_NAME = 'id' THEN CONCAT(table1, '.', COLUMN_NAME, ' AS ', table1, '_id')
             ELSE CONCAT(table1, '.', COLUMN_NAME)
        END
    ) INTO col_list1
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_NAME = table1 AND TABLE_SCHEMA = 'observatory';

    SELECT GROUP_CONCAT(
        CASE WHEN COLUMN_NAME = 'id' THEN CONCAT(table2, '.', COLUMN_NAME, ' AS ', table2, '_id')
             ELSE CONCAT(table2, '.', COLUMN_NAME)
        END
    ) INTO col_list2
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_NAME = table2 AND TABLE_SCHEMA = 'observatory';

    SET combined_col_list = CONCAT(col_list1, ',', col_list2);

    SET @query = CONCAT('SELECT l1.id,', combined_col_list, ' FROM ', table1, ' ',
                        'JOIN link l1 ON ', table1, '.id = l1.', table1, '_id ',
                        'JOIN ', table2, ' ON l1.', table2, '_id = ', table2, '.id');


    PREPARE stmt FROM @query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END //

DELIMITER ;

DELIMITER //

CREATE TRIGGER `observatory`.`update_sector_trigger`
BEFORE UPDATE ON `sector`
FOR EACH ROW
BEGIN
    SET NEW.date_update = NOW();
END //

DELIMITER ;


INSERT INTO `observatory`.sector (coordinates, light_intensity, foreign_objects, star_objects, unknown_objects, known_objects, notes) VALUES
('N35 W75', 3.5, 2, 5, 1, 3, 'Notes about this sector'),
('N45 E85', 2.8, 1, 3, 0, 2, 'Notes about another sector'),
('S25 W45', 4.2, 3, 4, 2, 1, 'Additional notes here'),
('N55 E95', 5.0, 4, 6, 2, 4, 'New notes about this sector');


INSERT INTO `observatory`.position (earth_position, sun_position, moon_position) VALUES
('Earth Position 1', 'Sun Position 1', 'Moon Position 1'),
('Earth Position 2', 'Sun Position 2', 'Moon Position 2'),
('Earth Position 3', 'Sun Position 3', 'Moon Position 3'),
('Earth Position 4', 'Sun Position 4', 'Moon Position 4');


INSERT INTO `observatory`.naturalobjects (type, galaxy, accuracy, flux, related_objects, notes) VALUES
('Galaxy', 'Milky Way', 0.95, 1.2, 'Object1, Object2', 'Notes about this natural object'),
('Star', 'Andromeda', 0.89, 0.8, 'Object3, Object4', 'Notes about another natural object'),
('Nebula', 'Orion', 0.92, 1.5, 'Object7, Object8', 'New notes about this natural object'),
('Planet', 'Solar System', 0.99, 1.0, 'Object5, Object6', 'Additional notes here');


INSERT INTO `observatory`.objects (type, accuracy, quantity, time, date, notes) VALUES
('Comet', 0.9, 5, '12:30:00', '2024-01-15', 'Notes about this object'),
('Asteroid', 0.85, 10, '13:45:00', '2024-02-20', 'Notes about another object'),
('Satellite', 0.95, 12, '15:30:00', '2024-04-10', 'New notes about this object'),
('Meteor', 0.8, 7, '14:00:00', '2024-03-25', 'Additional notes here');


INSERT INTO `observatory`.link (sector_id, objects_id, naturalobjects_id, position_id) VALUES
(1, 1, 1, 1),
(2, 2, 2, 2),
(3, 3, 3, 3),
(1, 3, 2, 2),
(4, 2, 2, 1);


COMMIT;

