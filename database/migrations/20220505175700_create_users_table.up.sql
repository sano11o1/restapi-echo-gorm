CREATE TABLE users (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255),
  `line_user_id` VARCHAR(255),
  created_at DATETIME(6),
  updated_at DATETIME(6),
  PRIMARY KEY (`id`)
);
