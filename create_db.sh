#!/usr/bin/env bash

echo "
create database s3mycloud CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
grant all privileges on s3mycloud.* to 's3mycloud'@'%' identified by 's3mycloud';
flush privileges;
" | mysql -uroot -p
