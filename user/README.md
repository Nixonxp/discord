# REST

## Генерация моделей на go-swagger
1. `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`
2. `swagger generate model -f ./api/swagger.json -t ./internal/`
3. `go mod tidy`