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

## PG cluster

```bash
./setup_pg_cluster.sh
```

## pg master

```postgresql
select * from pg_stat_replication;
select * from pg_replication_slots;
SHOW transaction_read_only;
```

## pg slave

```postgresql
SHOW transaction_read_only;

select * from pg_stat_wal_receiver;
SELECT pg_wal_replay_pause();
SELECT pg_wal_replay_resume();
select now() - pg_last_xact_replay_timestamp();

SELECT
  pg_is_in_recovery() AS is_slave,
  pg_last_wal_receive_lsn() AS receive,
  pg_last_wal_replay_lsn() AS replay,
  pg_last_wal_receive_lsn() = pg_last_wal_replay_lsn() AS synced,
  (
   EXTRACT(EPOCH FROM now()) -
   EXTRACT(EPOCH FROM pg_last_xact_replay_timestamp())
  )::int AS lag;
```
