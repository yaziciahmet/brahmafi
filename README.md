# UniSwap Pool Tracker

## Setup

```bash
docker-compose -f docker/docker-compose.yml up -d
```

## Run

```bash
go run main.go
```

## Endpoints

`/v1/api/pool/:poolId?block={'latest', '123123'}`: Get snapshot of the pool with the closest block number to `block`. It defaults to `latest`.
`/v1/api/pool/:poolId/historic`: Get all snapshots of the given pool id. This includes token delta information.

### Example Request

`/v1/api/pool/0xCBCdF9626bC03E24f779434178A73a0B4bad62eD?block=latest`

### NOTE

Pool contract address is used as pool id.
