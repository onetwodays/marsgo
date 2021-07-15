signal
CREATE USER im WITH PASSWORD 'im' SUPERUSER ;
```api
CREATE KEYSPACE signal WITH REPLICATION = 
 {'class': 'SimpleStrategy', 'replication_factor': 1};
 
 USE signal;
 
 CREATE TABLE test_table (
   id text,
   test_value text,
   PRIMARY KEY (id)
 );
 
 INSERT INTO test_table (id, test_value) VALUES ('1', 'one');
 INSERT INTO test_table (id, test_value) VALUES ('2', 'two');
 INSERT INTO test_table (id, test_value) VALUES ('3', 'three');
 
 SELECT * FROM test_table;
```