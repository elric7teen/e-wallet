# Descriptions
This is a simple e-wallet app.

# REQUIREMENT

 - Golang v1.13
 - Soda CLI (https://gobuffalo.io/en/docs/db/toolbox)
 - Postgres 11 or above

# Setup DB

 1. change directory to db_migrations/db_user
 2. update database.yml to fit your local machine or remote server
 3. to create database type `soda create -e $(your_env)`
 4. to populate db type `soda migrate -e $(your_env)`


## Update dependency
On project directory simply type `make vendor`

## Run Unit test
On project directory simply type `make test`
