CREATE TABLE IF NOT EXISTS result (
    id int NOT NULL PRIMARY KEY,
    time TIMESTAMP,
    duration float,
    suite varchar(64),
    dut_id int not NULL,
    milestone int,
    version varchar(32),
    host varchar(64) NOT NULL,
    test_id int not NULL,
    status boolean not NULL,
    reason blob,
    firmware_ro_version varchar(64),
    firmware_rw_version varchar(64)
);
