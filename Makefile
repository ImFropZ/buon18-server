db-create:
	migrate create -dir sql/migrations -ext sql $(name)

db-up:
	migrate -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable -path sql/migrations up

db-down:
	migrate -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable -path sql/migrations down

db-drop:
	migrate -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable -path sql/migrations drop

db-force:
	migrate -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable -path sql/migrations force $(version)

db-version:
	migrate -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable -path sql/migrations version
