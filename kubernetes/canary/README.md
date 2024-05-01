# Стратегия деплоя blue-green в kubernetes для сервиса auth

## Deployment
1. `minikube image load discord-app-auth:v1` - добавляем образ auth версии 1 в minikube
2. `kubectl --context minikube apply -f ./namespace.yaml` - создаем namespace
3. `kubectl --context minikube --namespace discord-app apply -f ./deploy_v1.yaml` - создаем деплой первой версии, ждем когда запустятся
4. `kubectl --context minikube --namespace discord-app apply -f ./service_load_balancer.yaml` - запускаем через load balancer
5. `minikube tunnel` - открываем туннель, для доступа извне
6. `curl --location "localhost:80"` - проверяем ответ от сервиса
### Выпускаем вторую версию приложения и собираем образ версии 2
1. `minikube image load discord-app-auth:v2` - добавляем образ auth версии 2 в minikube
2. `kubectl --context minikube --namespace discord-app apply -f ./deploy_v2.yaml` - применяем деплой второй версии, ждем когда запустится 1 под
3. `kubectl scale --replicas=9 deploy auth-1.1` - уменьшаем кол-во подов 1 версии, до 9
4. `while sleep 0.1; do curl "localhost:80"; done` - проверяем ответ от сервиса второй версии тоже приходит корректный
5. `kubectl scale --replicas=10 deploy auth-1.2` - увеличиваем кол-во подов 2 версии, до 10
6. `kubectl delete deploy auth-1.1` - тушим поды 1 версии