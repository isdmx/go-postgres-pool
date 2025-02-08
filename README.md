# go-postgres-pool

## Init

```bash
docker volume create test-postgres

docker run -d -p 5432:5432 \
      -v test-postgres:/var/lib/postgresql/data \
      -e POSTGRES_USER=test \
      -e POSTGRES_DB=test \
      -e POSTGRES_PASSWORD=secret \
      --name test-postgres \
      --hostname test-postgres \
 postgres:alpine
```
