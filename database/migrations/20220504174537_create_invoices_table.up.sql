CREATE TABLE invoices
(
    id serial,
    target_month varchar(50) NOT NULL,
    pricel int NOT NULL,
    PRIMARY KEY (id)
);
