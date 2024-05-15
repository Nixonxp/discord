# syntax=docker/dockerfile:experimental

FROM golang:1.22-alpine3.19 AS build

WORKDIR /app

# Копируем файлы go.mod и go.sum в WORKDIR
COPY ../go.mod ../go.sum ./
# Скачиваем необходимые Go модули (зависимости нашего проетка)
RUN go mod download

# Копируем все исходные go файлы нашего проекта в образ
# https://docs.docker.com/reference/dockerfile/#copy
COPY .. ./
# Собираем бинарный файл нашего приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main cmd/main.go

#######################################
# STAGE 2. FINAL STAGE
#######################################

FROM scratch AS final

WORKDIR /

COPY --from=build /bin/main /main

# Указываем какой порт необходимо слушать
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/main"]

