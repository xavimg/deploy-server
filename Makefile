serve:
	go run main.go
doc:
	swag init
migrate-up:
	migrate -path db/migration -database "postgresql://postgres:@localhost:5432/turingdb?sslmode=disable" -verbose up
