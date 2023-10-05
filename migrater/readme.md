## Migrate data from MongoDB to RDBMS


### Sequential Data pipeline 
- Parsing Oplog entries
- Transformation Process(Oplog to SQL)
- Execute and Migrate 


INSERT INTO test.student (_id, date_of_birth, is_graduated, name, roll_no) VALUES ('635b79e231d82a8ab1de863b','2000-01-30',false,'Selena Miller',51)

INSERT INTO test.student (_id, date_of_birth, is_graduated, name, roll_no) VALUES ('635b79e231d82a8ab1de863b', '2000-01-30', false, 'Selena Miller', 51);