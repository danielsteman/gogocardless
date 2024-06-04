# gogocardless

Go to [bank account data portal](https://bankaccountdata.gocardless.com/overview/)

Check out the [docs](https://developer.gocardless.com/bank-account-data/quick-start-guide)

## Development

Ripped [these](https://www.arhea.net/posts/2023-08-24-golang-vscode-configuration/) vscode settings

Run a local Postgres instance:

```bash
docker run -d \          
    --name gogocardless-postgres \
    -e POSTGRES_DB=gogocardless \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=admin \
    -p 5432:5432 \
    postgres:latest
```

Run service that fetches data from gocardless:

```
go run ./cmd
```

Run frontend:

```
cd web
yarn install
yarn dev
```

## Tests

Fire up the test database:

```bash
docker run -d \          
    --name gogocardless-postgres-test \
    -e POSTGRES_DB=gogocardless-test \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=admin \
    -p 5431:5432 \
    postgres:latest
```

Run tests:

```bash
go test ./...
```