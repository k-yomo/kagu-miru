# kagu-miru

## Architecture
![architecture](./docs/kagu_miru_architecture.png)

Prerequisites
- Go v1.17.x
- Node v14.8.x
- Docker

## Setup
```
make setup
```

## Run servers
Run servers locally and middlewares(Elasticsearch, Pub/Sub) on docker
```
make run
```

### Run item fetcher
To populate item data, we need to run item fetcher

1. Add application ids to `.env` in reference to [.env.sample](.env.sample)
2. Run servers
```
make run-item-fetcher
```

## Test
```
make test
```

## Docs
- [GraphQL](./docs/graphql)
