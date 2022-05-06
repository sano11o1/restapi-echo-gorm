CREATE TABLE invoices (
    `id` INT NOT NULL AUTO_INCREMENT,
    `target_year_and_month` varchar(50) NOT NULL,
    `name`  varchar(100) NOT NULL,
    `price` INT NOT NULL,
    `created_at` DATETIME(6),
    `updated_at` DATETIME(6),
    `sent_at` DATETIME(6) Default NULL,
    `user_id` INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (`id`)
);
