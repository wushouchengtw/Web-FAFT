CREATE TABLE IF NOT EXISTS Result (
    result_id int NOT NULL PRIMARY KEY,
    time TIMESTAMPTZ,
    duration float,
    suite varchar(64),
    dut_id int not NULL,
    milestone int,
    version varchar(32),
    host varchar(64) NOT NULL,
    test_id int not NULL,
    status boolean not NULL,
    reason blob,
    firmware_RO_Version varchar(64),
    firmware_RW_version varchar(64),
    FOREIGN KEY (dut_id) REFERENCES DUT(dut_id)
    FOREIGN KEY (test_id) REFERENCES Test(test_id)
);