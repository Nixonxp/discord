# Discord App

## Quick Start

### Собрать контейнеры приложения
```bash
$ make build
```

### Запустить приложение
```bash
$ make up
```

### Остановить приложение
```bash
$ make down
```

### Перезапустить приложение
```bash
$ make restart
```

### Swagger UI
1. `http://localhost:8801/swaggerui/#` - список всех методов приложения

## Tracing
http://localhost:16686/

## Prometheus

1. http://localhost:9091/

## Grafana

1. http://localhost:3000/
2. See more: https://grafana.com/grafana/dashboards/

## Pprof

1. `brew install graphviz`
2. http://localhost:8084/debug/pprof/
