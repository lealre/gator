# gator

Run migrations
```shell
goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator" up 
```
```shell
goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator" down 
```
