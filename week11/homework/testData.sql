create database if not exists testBase;

create table if not exists testBase.testTables(
    id bigint auto_increment primary key,
    A bigint null,
    B bigint null,
    C bigint null,
    constraint test_id unique (id, A)
);

INSERT INTO testTables (id, A, B, C)
VALUES (1,2,975,724),
       (2,3,747,428),
       (3,4,204,954),
       (4,5,294,889),
       (5,6,871,161),
       (6,7,886,597),
       (7,8,700,191),
       (8,9,668,427),
       (9,10,950,330),
       (10,11,921,412),
       (11,12,190,538),
       (12,13,242,440),
       (13,14,684,6),
       (14,15,798,222),
       (15,16,231,921),
       (16,17,933,24),
       (17,18,200,383),
       (18,19,935,687),
       (19,20,286,157),
       (20,21,519,539);

CREATE DATABASE IF NOT EXISTS testBase_NewOne;