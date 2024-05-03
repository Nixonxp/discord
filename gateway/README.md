# REST

## Агрегация всех моделей на go-swagger
1. `go-swagger-merger -o ./api/swagger.yaml ../auth/api/swagger.json`
2. `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`
3. `swagger generate model -f ./api/swagger.yaml -t ./internal/`
4. `go mod tidy`