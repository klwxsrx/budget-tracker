[![Build docker images](https://github.com/klwxsrx/budget-tracker/actions/workflows/build-docker-images.yml/badge.svg?branch=master&event=push)](https://github.com/klwxsrx/budget-tracker/actions/workflows/build-docker-images.yml)
# Для чего это
Читая различные книги по архитектуре нередко встречался с такими понятиями как `CQRS`, `MaterializedView` и `EventSourcing`, но не встречал работающего (не учебного) решения, где эти три подхода использовали полноценно и вместе.
Лучший способ разобраться _как оно работает_, понять _плюсы и минусы_ подходов это _запилить решение самому_, потому вот :)

_Моя цель_:
* Реализовать работающее приложение использующее подход CQRS и EventSourcing, тем самым получить полезный опыт в проектировании
* Попробовать для себя новые крутые программерские штуки, которые еще не использовал на работе
* Весело провести время решая нетривиальные задачи :)

# Техническое решение
## Схема инфраструктуры
![docs/infra.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/docs/infra.puml)
## Архитектура приложения
// TODO: написать про чистую архитектуру и слои, нарисовать схему со слоями/сервисами/агрегатами/репо и т.п.
## EventSourcing
// TODO: схему как оно работает на примере агрегата
## Домен приложения
![docs/domain.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/budget-tracker/master/docs/domain.puml)

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
- [ ] Реализовать модель приложения и сопутствующие операции с ними (сейчас полностью реализован Account)
- [ ] Добавить endpoint для получения данных авторизации в centrifugo
- [ ] Добавить APIGateway раскидывающий запроосы по соответствующим сервисам