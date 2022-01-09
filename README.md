# kagu-miru

## Architecture
![architecture](./docs/kagu_miru_architecture.png)

Prerequisites
- Go v1.17.x
- Node v16.13.x
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

### Run CMS
To work on the CMS used in our Media, we can run local CMS server
```
make run-media-studio
```

## Test
```
make test
```

## Docs
- [GraphQL](./docs/graphql)
