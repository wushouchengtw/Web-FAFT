CREATE TABLE IF NOT EXISTS DUT (
    dut_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
    model varchar(32) NOT NULL,
    board varchar(32) NOT NULL
);