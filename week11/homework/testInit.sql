show databases;
use testBase;

TRUNCATE TABLE testBase.testTables;
TRUNCATE TABLE testBase_NewOne.testTables;

DROP TABLE testBase.testTables;
DROP TABLE testBase_NewOne.testTables;

ALTER TABLE testBase.testTables MODIFY id bigint AUTO_INCREMENT;
ALTER TABLE testBase_NewOne.testTables MODIFY id bigint AUTO_INCREMENT;

select * from testBase.testTables;
select * from testBase_NewOne.testTables;

create database if not exists testBase;

create table if not exists testBase.testTables(
    id bigint primary key,
    cid bigint null,
    biz_id bigint null,
    biz varchar(128) null,
    uid bigint null,
    ctime bigint null,
    utime bigint null,
    constraint biz_type_id_uid unique (biz_id, biz, uid)
);