##
# SQLC
#
# @file
# @version 0.1

createmigration:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3333)/courses" -verbose up
migratedown:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3333)/courses" -verbose down


.PHONY: createmigration migrate migratedown
# end
