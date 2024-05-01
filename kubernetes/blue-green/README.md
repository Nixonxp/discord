# Стратегия деплоя blue-green в kubernetes для сервиса auth

## Deployment
1. `minikube image load discord-app-auth:v1` - добавляем образ auth версии 1 в minikube
2. `kubectl --context minikube apply -f .namespace.yaml` - создаем namespace
3. `kubectl --context minikube --namespace discord-app apply -f ./blue-deploy.yaml` - создаем деплой первой версии, ждем когда запустятся
4. `kubectl --context minikube --namespace main apply -f ./service_load_balancer.yaml` - запускаем через load balancer
5. `minikube tunnel` - открываем туннель, для доступа извне
6. `curl --location "localhost:80"` - проверяем ответ от сервиса
### Выпускаем вторую версию приложения и собираем образ версии 2
1. `minikube image load discord-app-auth:v2` - добавляем образ auth версии 2 в minikube
2. `kubectl --context minikube --namespace discord-app apply -f ./green-deploy.yaml` - создаем деплой второй версии, ждем когда запустятся, ждем когда запустятся все поды
3. `kubectl --context minikube --namespace discord-app apply -f ./service_load_balancer_v2.yaml` - применяем параметры для 2 версии load balancer
4. `curl --location "localhost:80"` - проверяем ответ от сервиса второй версии
5. `kubectl delete deploy auth-1.1` - тушим поды 1 версии