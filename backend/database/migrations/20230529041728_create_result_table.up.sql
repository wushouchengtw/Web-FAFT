CREATE TABLE IF NOT EXISTS Result (
    result_id int NOT NULL PRIMARY KEY,
    time TIMESTAMP,
    duration varchar(10),
    suite varchar(40),
    dut_id int not NULL,
    build_version varchar(20),
    host varchar(60) NOT NULL,
    test_id int not NULL,
    result boolean not NULL,
    firmware_RO_Version varchar(50),
    firmware_RW_version varchar(50),
    FOREIGN KEY (dut_id) REFERENCES DUT(dut_id)
    FOREIGN KEY (test_id) REFERENCES Test(test_id)
);