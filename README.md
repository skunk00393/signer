1. Create a postgres database named signer
2. Execute the follwoing queries:
```
create table users (
	id serial primary key,
	username varchar(50) not null,
	password varchar(65) not null,
	created_on timestamp not null
)
```
```
create table signature (
	id serial primary key,
	user_id integer not null references users(id),
	hashed_sig varchar(500) not null,
	signed_on timestamp not null
)
```
```
create table qas (
	id serial primary key,
	question varchar(500) not null,
	answer varchar(500) not null,
	answered_on timestamp not null,
	sig_id integer not null references signature(id)
)
```
3. Run go run . while in the project signer folder
 
NOTE:  I tested it using Postman because I couldn't get it working with grpcurl (searched for the issue and some answers said the issue is connected to Windows)

PS: Will dockerize ASAP
