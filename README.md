# Go.Mod Link

## Overview
У нас есть набор сервисов/библиотек, и хочется понимать, насколько сильно они между собой связаны. Берем репозитории с сервисами, достаем из них `go.mod` и ищем в них зависимости между выбраными репозиториями.

Результат: видим какие сервисы/библиотеки от кого зависят.

Идеальный результат: зависимостей между **N** сервисов равно **N - 1**, т.е. по сути у нас пайплайн из сервисов, когда в каждый сервис может обратиться только один сервис. Очевидно, что в жизни так бывает редко, но это идеальный результат который можно принять за эталон. Можно придумать и другие метрики, но эта метрика довольно простая и ее легко считать.

## Configuration
Создаем json- файл в формате:
```
{
  "repo": [
    {
      "url": "ssh://git@gitlab.devops.ourrepo:51111/cp/project-first.git",
      "branch": "develop"
    },
    {
      "url": "ssh://git@gitlab.devops.ourrepo:51111/cp/project-second.git",
      "branch": "master"
    },
    {
      "url": "ssh://git@gitlab.devops.ourrepo:51111/cp/project-third.git",
      "branch": "branch_with_feature"
    }
  ]
}
```
где:
 * **url** - адрес репозитория (обязательно ssh, https не поддерживается, если необходимо, создайте issue - добавлю)
 * **branch** - ветка из которой брать зависимости

## Run
Запуск:
```
gomodlink --from example.json
```
где файл `example.json` - описание ваших репозиториев

## Output
Пример отчета:
```
Repository: gitlab.devops.ourdomain.dev/api-gateway (5)
    gitlab.devops.ourdomain.dev/project-first/pkg/project-first-api
    gitlab.devops.ourdomain.dev/project-second/pkg/project-second-api/v2
    gitlab.devops.ourdomain.dev/project-third/pkg/project-third-api/v2
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/project-fifth/pkg/project-fifth-api

Repository: gitlab.devops.ourdomain.dev/project-first (5)
    gitlab.devops.ourdomain.dev/project-fourth/pkg/project-fourth-api
    gitlab.devops.ourdomain.dev/project-second/pkg/project-second-api/v2
    gitlab.devops.ourdomain.dev/project-third/pkg/project-third-api/v2
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/project-fifth/pkg/project-fifth-api

Repository: gitlab.devops.ourdomain.dev/project-second (5)
    gitlab.devops.ourdomain.dev/project-fourth/pkg/project-fourth-api
    gitlab.devops.ourdomain.dev/project-first/pkg/project-first-api
    gitlab.devops.ourdomain.dev/project-third/pkg/project-third-api/v2
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/project-fifth/pkg/project-fifth-api

Repository: gitlab.devops.ourdomain.dev/project-fourth (2)
    gitlab.devops.ourdomain.dev/customer/pkg/customer-api
    gitlab.devops.ourdomain.dev/project-fifth/pkg/project-fifth-api

Repository: gitlab.devops.ourdomain.dev/customer (2)
    gitlab.devops.ourdomain.dev/project-first/pkg/project-first-api
    gitlab.devops.ourdomain.dev/project-fifth/pkg/project-fifth-api

Repository: gitlab.devops.ourdomain.dev/project-third (0)

Repository: gitlab.devops.ourdomain.dev/project-fifth (0)

Total repository: 7
Total dependencies: 19
AVG dependencies: 3.17
```
 * **Total repository** - кол-во репозиториев
 * **Total dependencies** - кол-во найденых зависимостей
 * **AVG dependencies** - коэффициент, показывающий во сколько раз больше зависимостей от идеального показателя. Вы можете выбрать оптимальное значение для вас и ориентироваться на него.
