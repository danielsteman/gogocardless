## Development

Local postgres:

```
docker run -d \
    --name gogocardless-postgres \
    -e POSTGRES_DB=gogocardless \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=admin \
    -p 5432:5432 \
    postgres:latest
```