# Go.Mod Link

## Overview
У нас есть набор сервисов/библиотек, и хочется понимать, насколько сильно они между собой связаны. Берем репозитории с сервисами, достаем из них `go.mod` и ищем в них зависимости между выбраными репозиториями.

Результат: видим какие сервисы/библиотеки от кого зависят.

Идеальный результат: зависимостей между *N* сервисов равно *N - 1*, т.е. по сути у нас пайплайн из сервисов, когда в каждый сервис может обратиться только один сервис. Очевидно, что в жизни так бывает редко, но это идеальный результат который можно принять за эталон. Можно придумать и другие метрики, но эта метрика довольно простая и ее легко считать.

## Configuration
Создаем json- файл в формате:
```
{
  "repo": [
    {
      "url": "ssh://git@gitlab.devops.ourrepo:51111/cp/projectfirst.git",
      "branch": "develop"
    },
    {
      "url": "ssh://git@gitlab.devops.ourrepo:51111/cp/projectsecond.git",
      "branch": "master"
    },
    {
      "url": "ssh://git@gitlab.devops.ourrepo:51111/cp/projectthird.git",
      "branch": "branch_with_feature"
    }
  ]
}
```
где:
 * url - адрес репозитория (обязательно ssh, https не поддерживается, если необходимо, создайте issue - добавлю)
 * branch - ветка из которой брать зависимости

## Run
Запуск:
```
gomodlink --from example.json
```
где файл `example.json` - описание ваших репозиториев

## Output
Пример отчета, где:
 * Total repository - кол-во репозиториев
 * Total dependencies - кол-во найденых зависимостей
 * AVG dependencies - коэффициент, показывающий насколько больше зависимостей от идеального показателя. Вы можете выбрать оптимальное значение для вас и ориентироваться на него.
```
Repository: gitlab.devops.ourdomain.dev/api-gateway (5)
    gitlab.devops.ourdomain.dev/event-manager/pkg/event-manager-api
    gitlab.devops.ourdomain.dev/order-manager/pkg/order-manager-api/v2
    gitlab.devops.ourdomain.dev/item-catalog/pkg/item-catalog-api/v2
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/product-item/pkg/product-item-api

Repository: gitlab.devops.ourdomain.dev/event-manager (5)
    gitlab.devops.ourdomain.dev/notification/pkg/notification-api
    gitlab.devops.ourdomain.dev/order-manager/pkg/order-manager-api/v2
    gitlab.devops.ourdomain.dev/item-catalog/pkg/item-catalog-api/v2
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/product-item/pkg/product-item-api

Repository: gitlab.devops.ourdomain.dev/order-manager (5)
    gitlab.devops.ourdomain.dev/notification/pkg/notification-api
    gitlab.devops.ourdomain.dev/event-manager/pkg/event-manager-api
    gitlab.devops.ourdomain.dev/item-catalog/pkg/item-catalog-api/v2
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/product-item/pkg/product-item-api

Repository: gitlab.devops.ourdomain.dev/notification (2)
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/product-item/pkg/product-item-api

Repository: gitlab.devops.ourdomain.dev/customer (2)
    gitlab.devops.ourdomain.dev/event-manager/pkg/event-manager-api
    gitlab.devops.ourdomain.dev/product-item/pkg/product-item-api

Repository: gitlab.devops.ourdomain.dev/item-catalog (0)

Repository: gitlab.devops.ourdomain.dev/product-item (0)

Total repository: 7
Total dependencies: 19
AVG dependencies: 3.17
```
