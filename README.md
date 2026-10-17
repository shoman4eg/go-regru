# go-regru

[![Go Reference](https://pkg.go.dev/badge/github.com/shoman4eg/go-regru.svg)](https://pkg.go.dev/github.com/shoman4eg/go-regru)

Go-клиент для [REG.API 2](https://www.reg.ru/reseller/api2doc) — API регистратора REG.RU. Устроен по образцу [google/go-github](https://github.com/google/go-github), без сторонних зависимостей.

## Установка

```bash
go get github.com/shoman4eg/go-regru
```

Требуется Go 1.26+.

## Быстрый старт

```go
package main

import (
	"context"
	"fmt"
	"log"

	regru "github.com/shoman4eg/go-regru"
)

func main() {
	client := regru.NewClient(nil).WithAuthPassword("test", "test")
	ctx := context.Background()

	answer, _, err := client.Domain.Check(ctx, &regru.DomainCheckOptions{
		Domains: []regru.ServiceIdent{{DName: "example.ru"}, {DName: "example.com"}},
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range answer.Domains {
		if err := d.Err(); err != nil {
			fmt.Println(d.DName, "ошибка:", err)
			continue
		}
		fmt.Println(d.DName, d.Result)
	}
}
```

Для проверки подходит публичный тестовый аккаунт `test` / `test`.

## Перед началом работы

- **Доступ по IP.** По умолчанию API принимает запросы только с адресов, разрешённых в [настройках API](https://www.reg.ru/user/account/#/settings/api/). С других адресов придёт ошибка `ACCESS_DENIED_FROM_IP` (`regru.ErrCodeAccessDeniedFromIP`).
- **Тестовый аккаунт `test`/`test`.** Он проверяет параметры и возвращает ответ, но ничего не выполняет и деньги не списывает. Данные в ответах условные и могут не совпадать с ответами боевого аккаунта. Отличаются даже типы: например, `prio` и `service_id` в `zone/get_resource_records` у тестового и боевого аккаунтов приходят по-разному, то числом, то строкой. Клиент принимает оба варианта.
- **Проверка с боевыми данными.** Для этого предназначены функции `nop` (`client.Nop`, `client.Domain.Nop`, `client.Zone.Nop` и другие): они ничего не меняют и отвечают быстро.
- **Лимит.** Не больше 1200 запросов в час с одного IP и с одного аккаунта.

## Сервисы

Функции API сгруппированы по категориям, как в документации:

| Поле клиента | Категория API | Что делает |
|---|---|---|
| `client.Nop`, `GetUserID`, … | корневые функции | проверка доступа, ID пользователя и услуги |
| `client.User` | `user/*` | учётная запись, баланс, статистика |
| `client.Bill` | `bill/*` | счета |
| `client.Service` | `service/*` | заказ, продление и управление услугами |
| `client.Domain` | `domain/*` | регистрация, перенос, контакты, NS |
| `client.Zone` | `zone/*` | DNS-записи |
| `client.DNSSEC` | `dnssec/*` | DNSSEC |
| `client.Hosting` | `hosting/*` | хостинг |
| `client.Folder` | `folder/*` | папки услуг |
| `client.Shop` | `shop/*` | магазин доменов |

Каждый метод принимает `context.Context` и структуру параметров `XxxOptions` и возвращает ответ, `*regru.Response` и ошибку:

```go
_, _, err := client.Zone.AddAlias(ctx, &regru.ZoneAddAliasOptions{
	ZoneDomainsOptions: regru.ZoneDomainsOptions{ServiceIdent: regru.ServiceIdent{DomainName: "example.ru"}},
	Subdomain:          "www",
	IPAddr:             "192.0.2.1",
})
```

## Аутентификация

**По паролю.** Подходит основной пароль или отдельный пароль для API из настроек партнёра:

```go
client := regru.NewClient(nil).WithAuthPassword("login", "api-password")
```

**По подписи.** Запрос подписывается RSA-SHA512 приватным ключом, сертификат которого загружен в настройках API:

```go
client := regru.NewClient(httpClient).WithSignature("login", privateKey) // *rsa.PrivateKey
```

**SSL-сертификат.** При подписи он обязателен, при пароле — если в настройках API загружен хотя бы один сертификат. Сертификат передаётся через TLS-настройки вашего `http.Client`:

```go
cert, err := tls.LoadX509KeyPair("my.crt", "my.key")
if err != nil {
	log.Fatal(err)
}
httpClient := &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
}}
```

Язык текстов ошибок задаётся полем `client.Lang`: `en` (по умолчанию), `ru` или `th`.

## Ошибки

Если API вернул `"result": "error"`, метод возвращает `*regru.ErrorResponse` с кодом и текстом ошибки:

```go
_, _, err := client.Nop(ctx)
var apiErr *regru.ErrorResponse
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.Code, apiErr.Text) // PASSWORD_AUTH_FAILED Username/password incorrect
}

if regru.ErrorCode(err) == regru.ErrCodeIPRateLimit {
	// превышен лимит запросов (1200 в час с одного IP или аккаунта)
}
```

Функции, которые принимают список доменов или услуг, могут завершиться успешно в целом, но с ошибками по отдельным элементам. Каждый такой элемент нужно проверять через `item.Err()`, как в примере в начале.

## Особенности

- API возвращает числа то числом, то строкой (`12345` и `"12345"`). Для полей ответа используются типы `regru.FlexInt`, `regru.FlexFloat` и `regru.FlexString`, которые принимают оба варианта.
- Флаги в параметрах имеют тип `regru.Bool` и отправляются как `0`/`1`.
- Параметры, набор которых зависит от типа услуги или доменной зоны, передаются как `map[string]any`. Это `service/create`, `service/update` и контакты домена (`regru.DomainContacts`).
- `Domain.DownloadCertificate` возвращает PDF в виде `[]byte`.

## Разработка

```bash
gofmt -l . && go vet ./... && go test ./...
```

Тесты работают с локальным `httptest`-сервером на примерах из документации и не обращаются к настоящему API.

## Лицензия

[MIT](LICENSE)
