# webook
A Full-Workflow Backend Practice Demo

The Go module path is `github.com/johsqn/webook`. For a fresh MySQL volume,
`docker compose up -d` creates the `webook` database and user. Set `MYSQL_DSN`
to use that database before starting the Go service.

Existing `mysql_data` volumes keep their previous database and users because
MySQL initialization only runs on an empty volume. If you have existing data,
migrate the database and credentials before switching your `MYSQL_DSN`.
