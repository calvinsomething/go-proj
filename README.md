# go-proj

A website with a Go backend.

### Installation

Install dependencies for the server

```bash
cd ./server && go mod vendor
```

and client.

```bash
cd .. && npm i
```

Set environment variables for:

- SERVER_PORT

- CLIENT_PORT

- DB_PORT
- DB_ROOT_PW
- DB_USER
- DB_PASSWORD

_If you need to change the DB environment variables at any point, make sure to delete `/data` before running the container.
Otherwise you can update them manually inside the container._

### How To Run

```bash
docker compose up
```
