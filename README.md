[![Build docker images](https://github.com/klwxsrx/budget-tracker/actions/workflows/build-docker-images.yml/badge.svg?branch=master&event=push)](https://github.com/klwxsrx/budget-tracker/actions/workflows/build-docker-images.yml)
# Что это
Читая различные книги по архитектуре нередко встречался с такими понятиями как `CQRS`, `MaterializedView` и `EventSourcing`, но не встречал работающего решения, где эти три подхода использовали полноценно и вместе.
Лучший способ разобраться _как оно работает_, понять _плюсы и минусы_ подходов это _запилить решение самому_ :)

_Цель проекта_:
* Реализовать приложение использующее подход CQRS и EventSourcing
* Попробовать для себя новые крутые программерские штуки, которые еще не использовал на работе
* Весело провести время решая нетривиальные задачи :)

# Техническое решение
## Схема сервисов
![docs/infra.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/docs/infra.puml)

## Архитектура приложения
В проекте используется разделение кода на слои в соотвествии с принципами Clean Architecture Дядюшки Боба, применяется правило независимости внутренних слоев от внешних.

За соблюдением независимости слоев следит утилита [go-cleanarch](https://github.com/roblaszczak/go-cleanarch), её использование см. ниже.

Структура кода приложения:
```
pkg:
    {имя приложения}
        {слой}
        {слой}
```
Например `pkg/budget/domain`.

## Домен приложения
![docs/domain.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/docs/domain.puml)

## EventSourcing
При попытке скрестить Агрегаты DDD с ES получилась следующая схема:

![docs/es.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/docs/es.puml)

Агрегат получаемый из БД строится путем применения его доменных событий к его пустому состоянию в порядке создания этих событий.
Для этого выделен `*State`, хранящий все свойства агрегата, изменяющий их значения при применении событий к нему. Кстати этот же `*State` можно сериализовать в БД.

Агрегат по-прежнему сам следит за соблюдением бизнес-правил при попытке совершения операций над ним, но меняет внутреннее состояние через применение доменного события к `*State`.

Для сохранения/получения агрегата в БД используется паттерн Repository, за которым прячется EventStore.

## CQRS
* Write-операции выполняются в сервисе `Budget` (например создание кошелька). 
В контроллере создается команда и далее обрабатывается соответствующим ApplicationService'ом. 
Из EventStore выгребается агрегат, меняется путем применения доменного события, изменения агрегата обратно отправляются в EventStore (БД).
* Далее доменное событие из EventStore (БД) отправляется в EventBus (используя TransactionalOutbox).
* На события подписан сервис `BudgetView`, который дописывает изменения в `MaterializedView` БД.
* Read-операции (например получение инфы по кошельку) дергаются сервисом `BudgetView` (так-то стоит распилить сервис и его отвественность на два).

# Запуск
## Требования
* `docker` и `docker-compose`
* (опционально) `go` - если хотите собрать docker-образ приложения сами
* (опционально) [golangci-lint](https://github.com/golangci/golangci-lint) - для запуска линтера по проекту командой `make lint`
* (опционально) [go-cleanarch](https://github.com/roblaszczak/go-cleanarch) - для проверки нарушения слоев архитектуры командой `make go-cleanarch`

## Запуск актуальных версий приложений с docker.io
Выполняем в корне проекта:
```shell
docker-compose up -d
```

Два сервиса будут доступны на:
* Budget - `127.0.0.1:8080`, [openAPI-схема](https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/api/budget/api.yml)
* BudgetView - `127.0.0.1:8081`, [openAPI-схема](https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/api/budgetview/api.yml)

## Сборка docker-образов и запуск собранных приложений
* Копируем файл `docker-compose.override.yml.dist` в `docker-compose.override.yml`
* Правим скопированный конфиг под нужды, оставляем секцию `build` для требуемых приложений
* Выполняем в корне проекта:
```shell
make bin/budget bin/budgetview && docker-compose up -d --build
```

## Тестирование, линтер
Выполняем в корне проекта:
* `make test` - запуск unit-тестов
* `make lint` - запуск golangci-lint, см. конфиг [.golangci.yml](https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/.golangci.yml) в корне проекта
* `make check-arch` - запуск проверки нарушения слоев архитектуры

# TODO:
- [ ] Реализовать модель приложения и сопутствующие операции (сейчас полностью реализован Account)
- [ ] Добавить endpoint для получения данных авторизации в centrifugo
- [ ] Добавить APIGateway раскидывающий запроосы по соответствующим сервисам
- [ ] Добавить шаг `make check-arch` в билд Github Actions