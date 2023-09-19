drop table if exists users; 

create table users (
    id integer primary key, 
    username text, 
    balance numeric
);

INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 2000);
INSERT INTO users (id, username, balance) VALUES (2, 'second_user', 0);
INSERT INTO users (id, username, balance) VALUES (3, 'third_user', 1000);