# graphQlDemo

In order to get started you need to set up some pre-reqs.

1. Install postgres sql `brew install postgresql@15`
   - Once installed on a terminal type in `psql postgres`
   - Then set up a user with root access using below
     - `CREATE ROLE [user] WITH LOGIN PASSWORD ['password'];`
     - `ALTER ROLE user CREATEDB;`
   - After set up quit by typing `\q`
   - Now log in with the user you just created by typing in: `psql postgres -U user`
   - Verify everything was created as expected by typing in `\du`. You should see your new user listed as a super user and other roles
2. Install [pg admin](https://www.pgadmin.org/download/)
   - Use pgadmin to try to connect to your local postgres to verify you have it up and running by starting the pg admin server
   - host: "localhost"
   - user: [user]
   - password: [password]
   - maintenance database: "postgres"
   - Now create a database called "tools-back"
3. In order to develop locally their are a couple of environment variables you will want to set within your .bash_profile or equivalent files. These are used in order to develop using a local postgres sql db.
   Make a .env file with this in the root project directory. Use the command `. ./.env` to set the environment variables. Use `printenv` to verify environment variables are correctly set.

```bash
export POSTGRES_DB_HOST="localhost"
export POSTGRES_DB_PORT=5432
export POSTGRES_DB_USER="[user]"
export POSTGRES_DB_PASSWORD="[password]"
export POSTGRES_DB="tools-back"
export PORT=3546
```
4. Make sure postgres is running `brew services start postgresql@15`
5. Start the server by running `go run main.go`
6. In a web browser navigate to (https://localhost:3546/playground). If you see the graphql playground you have a running server.

Some other useful commands:

- Codegen changes after graphql or a schema change `go generate .`
- Build `go build -v`
- To create a new schema `go run -mod=mod entgo.io/ent/cmd/ent new User`
