# go-microservices-test

## Requirements
1. Install Make. You can use choco in windows `choco install make` or in mac `brew install make`
2. Have docker compose (Docker Desktop) installed.

## Running the project
1. In the main folder of the repository, open a terminal.
2. Type `cd project`
3. Type `make up_build` and wait for everything to start.
4. Open the Database in your preferred db admin tool (Recommended DBeaver).
5. Execute the script `db-init.sql` located in the project folder.
6. Load the CSV following instructions in `db-init.sql file`
7. Test the microservice using the following curl:
```
curl --location 'localhost:8080/consumption?meters_ids=1&meters_ids=2&start_date=2023-06-01&end_date=2023-07-05&kind_period=daily'
```