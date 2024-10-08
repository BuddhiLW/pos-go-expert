# Notes on SQLC

## Commands
### Create migrations (go migrate)

``` sh
migrate create -ext=sql -dir=sql/migrations -seq init
```

### Apply `up` changes

``` sh
migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3333)/courses" -verbose up
```

``` 
2024/10/01 14:58:44 Start buffering 1/u init
2024/10/01 14:58:44 Read and execute 1/u init
2024/10/01 14:58:44 Finished 1/u init (read 2.322348ms, ran 22.328467ms)
2024/10/01 14:58:44 Finished after 30.727339ms
2024/10/01 14:58:44 Closing source and database
```

#### Case Dirty Database

If database is dirty...

``` sh
migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3333)/courses" -verbose up
```

```
2024/09/16 12:20:45 error: Dirty database version 1. Fix and force version.
```

You can force changes

``` sh
migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3333)/courses" -verbose force 1
```


