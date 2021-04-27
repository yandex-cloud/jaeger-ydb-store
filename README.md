# YDB storage plugin for Jaeger

This is a storage backend using [Yandex.Database](https://cloud.yandex.ru/services/ydb) for [Jaeger](https://github.com/jaegertracing/jaeger)

## components

- collector/query: these are docker images with jaeger-collector/jaeger-query and plugin built-in into image
- watcher (see cmd/schema/watcher): this is schema-watcher that creates new tables for spans/indexes and removes old ones

## docker images

- `cr.yandex/yc/jaeger-ydb-collector`
- `cr.yandex/yc/jaeger-ydb-query`
- `cr.yandex/yc/jaeger-ydb-watcher`

## how to run

```
cp docker-compose.example.yml docker-compose.yml
# edit docker-compose yml to set db dsn
docker-compose up -d
```

## environment variables

- YDB_ADDRESS (string): db endpoint host:port to connect to
- YDB_SA_META_AUTH (bool): use metadata to authorize requests ([documentation](https://cloud.yandex.com/docs/compute/operations/vm-connect/auth-inside-vm#auth-inside-vm))
- YDB_SA_ID (string): service account id for Yandex.Cloud authorization (doc on service accounts: https://cloud.yandex.com/docs/iam/concepts/users/service-accounts)
- YDB_SA_KEY_ID (string): service account key id for Yandex.Cloud authorization
- YDB_SA_PRIVATE_KEY_FILE (string): path to service account private key for Yandex.Cloud authorization
- YDB_PATH (string): database path
- YDB_FOLDER (string): folder to store data in)
- YDB_CONNECT_TIMEOUT (duration, default: 10s): db connect timeout
- YDB_WRITE_TIMEOUT (duration, default: 1s): write queries timeout
- YDB_READ_TIMEOUT (duration, default: 10s): read queries timeout
- YDB_READ_QUERY_PARALLEL (integer, default: 16): controls number of parallel read subqueries
- YDB_POOL_SIZE (integer, default: 100): db session pool size
- YDB_QUERY_CACHE_SIZE (integer, default: 50): db query cache size
- YDB_WRITER_BUFFER_SIZE (integer, default: 1000): span buffer size for batch writer
- YDB_WRITER_BATCH_SIZE (integer, default: 100): number of spans in batch write calls
- YDB_WRITER_BATCH_WORKERS (integer, default: 10): number of workers processing batch writes
- YDB_INDEXER_BUFFER_SIZE (integer, default: 1000): span buffer size for indexer
- YDB_INDEXER_MAX_TRACES (integer, default: 100): maximum trace_id count in a sinigle index record
- YDB_INDEXER_MAX_TTL (duration, default: 5s): maximum amount of time for indexer to batch trace_idы for index records
- YDB_SCHEMA_NUM_PARTITIONS (integer, default: 10): number of partitioned tables per day. Changing it requires recreating full data set
- YDB_TOKEN (string): auth token for internal purposes
- YDB_FEATURE_SPLIT_BY_LOAD (bool): enable table split by load feature
- YDB_FEATURE_COMPRESSION (bool): enable table compression feature, used for span storage

## schema watcher configuration

- WATCHER_AGE (duration, default: 24h): delete partition tables older than this value
- WATCHER_INTERVAL (duration, default: 5m): check interval

## conference talks

- https://youtu.be/nyt_e4ULrUo?t=660
- https://youtu.be/hXH_tRBxFnA?t=11471
