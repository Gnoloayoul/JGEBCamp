create database if not exists testBase;

create table if not exists testBase.testTables(
    id bigint auto_increment primary key,
    A bigint null,
    B bigint null,
    C bigint null,
    constraint test_id unique (id, A)
);