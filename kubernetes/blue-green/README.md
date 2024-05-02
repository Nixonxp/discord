# Стратегия деплоя blue-green в kubernetes для сервиса auth

## Deployment
1. `minikube image load discord-app-auth:v1` - добавляем образ auth версии 1 в minikube
2. `kubectl --context minikube apply -f ./namespace.yaml` - создаем namespace
3. `kubectl --context minikube --namespace discord-app apply -f ./blue-deploy.yaml` - создаем деплой первой версии, ждем когда запустятся
4. `kubectl --context minikube --namespace discord-app get pods --show-labels` - проверяем что поды поднялись
5. `kubectl --context minikube --namespace discord-app apply -f ./service_load_balancer.yaml` - запускаем через load balancer и роутим на 1 версию
6. `minikube tunnel` - открываем туннель, для доступа извне
7. `curl --location "localhost:80"` - проверяем ответ от сервиса
### Выпускаем вторую версию приложения и собираем образ версии 2
1. `minikube image load discord-app-auth:v2` - добавляем образ auth версии 2 в minikube
2. `kubectl --context minikube --namespace discord-app apply -f ./green-deploy.yaml` - создаем деплой второй версии, ждем когда запустятся, ждем когда запустятся все поды
3. `kubectl --context minikube --namespace discord-app get pods --show-labels` - проверяем что поды поднялись
4. `kubectl --context minikube --namespace discord-app apply -f ./service_load_balancer_v2.yaml` - применяем параметры для 2 версии load balancer и роутим на них
5. `curl --location "localhost:80"` - проверяем ответ от сервиса второй версии
6. `kubectl delete deploy auth-blue` - тушим поды 1 версии