<h1 align="center">Money transfer API</h1>

## 📜 Summary
- [About](#About)
- [Libs/Dependencies](#Libs/Dependencies)
- [Run](#Run)
- [Tests](#Tests)
- [Endpoints](#Endpoints)


<a id="About"></a> 
## 📃 About
This code is a challenge of the https://app.devgym.com.br/ platform. Its a small API that simulates a money transfer from one user's 
bank account to another one. I used this challenge to practice a little bit about database transactions. There's only two endpoints in this app and you can check it out in this session <a href="#Endpoints">endpoints</a>. 

---
<a id="Libs/Dependencies"></a> 
## 🗄 Libs/Dependencies </br>

| Name        | Description | Documentation | Installation |
| ----------- | ----------- | ------------- | ----------- |     
| pgx      | postgres database driver       |  github.com/jackc/pgx/v4 |  go get github.com/jackc/pgx/v4      |
| chi               |  http router  lib | https://github.com/go-chi/chi                   | go get github.com/go-chi/chi   |
| godotenv             | .env vars manager              | github.com/joho/godotenv             | go get github.com/joho/godotenv    | 

---

<a id="Run"></a> 
## ⚙️ Run

There's two ways of starting this project: using docker to start the webserve and the database or starting the database via docker and
starting the server on your machine.


### Using docker-compose for database and webserver

Inside the root folder, go to the .ENV file and set the variable DB_HOST with the value of "db". <br>

Run one of the commands below to build golang image:

```bash
docker-compose -f docker-compose.production.yml build
```

```bash
make build
```

Then run one of the commands below to start the containers:

```bash
docker-compose -f docker-compose.production.yml up -d
```

```bash
make run_prod
```

If you want to destroy it all, run one of the commands below:

```bash
docker-compose -f docker-compose.production.yml down
```

```bash
make down
```

### Using docker-compose for database and webserver locally

Inside the root folder, go to the .ENV file and set the variable DB_HOST with the value of "localhost". <br>

Run one of the commands below to start the database:

```bash
docker compose up -d
```

```bash
make db_up
```

then start the api:

```bash
go run main.go
```

```bash
make run
```

If you want to destroy the database's docker, run one of the commands below:

```bash
docker compose down 
```

```bash
make db_down
```

<a id="Tests"></a> 
## 🧪 Tests

All the tests in this project are integration ones. So, befofe running it is required
to instantiate the database first. After you start it, run one of the commands below:

```bash
make tests
```

```bash
go test ./test -v
```

<a id="Endpoints"></a> 
## 💻 Endpoints

Consulting user balance: <br>
userId: user identifier. Example: 1<br>

Request: 

```bash
curl localhost:3000/balance/{userId}
```

Response: 

```json
{"balance":2000}
```

Transfer amount from one user to another:

Request body:<br>
debtorId - user identifier that is transfering the amount<br>
beneficiaryId - user identifier that is receiving the amount<br>
amount - value that is being transferred<br>

Request: 

```bash
curl --location 'http://localhost:3000/transfer' \
    --header 'Content-Type: application/json' \
    --data '{
        "amount": 100.0,
        "debtor_id": 1,
        "beneficiary_id": 2
    }'
```
