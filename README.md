# config-audit
 
Go-утилита для анализа JSON/YAML конфигураций веб-приложений и поиска потенциально опасных настроек.

Проект поддерживает CLI, HTTP REST API и gRPC API. Вся бизнес-логика анализа вынесена в общее ядро, поэтому CLI и API используют одинаковые правила и возвращают одинаковые результаты.

## Что реализовано

- Парсинг JSON и YAML конфигураций.
- CLI с позиционным путем к файлу.
- Чтение конфигурации из `stdin` через `--stdin`.
- Флаги `-s` и `--silent`, чтобы не возвращать ошибочный exit code при найденных проблемах.
- Вывод проблем с уровнями `LOW`, `MEDIUM`, `HIGH`.
- Краткое описание проблемы и рекомендация по исправлению.
- Exit code `1`, если найдены проблемы и не включен silent-режим.
- HTTP REST API для анализа конфигурации.
- gRPC API для анализа конфигурации.
- Рекурсивный анализ директории с конфигами.
- Проверка прав доступа к файлу через `os.Stat`.
- Unit и integration тесты.
- GitHub Actions CI.

## Правила анализа

- `debug-logging` - включенный debug/trace logging или `debug: true`.
- `plain-password` - пароль, token, secret или API key прямо в конфиге.
- `bind-all-interfaces` - bind/listen на `0.0.0.0`.
- `tls-disabled` - отключенный TLS или отключенная проверка сертификата.
- `weak-algorithm` - слабые алгоритмы вроде MD5, SHA-1, DES, 3DES, RC4.
- `file-permissions` - слишком широкие права доступа к файлу, если включен `--check-permissions`.

## Требования

- Go 1.24+

## Проверка проекта

Из корня проекта:

```powershell
go test ./...
go vet ./...
```

Ожидаемый результат: тесты проходят, `go vet` не выводит ошибок.

## CLI

Проверить JSON-конфиг:

```powershell
go run ./cmd/config-audit ./testdata/configs/debug.json
```

Ожидаемый результат: будет найден `debug-logging`, программа завершится с exit code `1`.

Проверить безопасный конфиг:

```powershell
go run ./cmd/config-audit ./testdata/configs/safe.json
```

Ожидаемый результат:

```text
No issues found.
```

Проверить YAML-конфиг:

```powershell
go run ./cmd/config-audit ./testdata/configs/weak_algorithm.yaml
```

Проверить `stdin`:

```powershell
'{"log":{"level":"debug"}}' | go run ./cmd/config-audit --stdin
```

Проверить silent-режим:

```powershell
go run ./cmd/config-audit --silent ./testdata/configs/debug.json
```

Проверить рекурсивный анализ директории:

```powershell
go run ./cmd/config-audit --silent --recursive ./testdata/configs
```

Проверить права доступа через `os.Stat`:

```powershell
go run ./cmd/config-audit --silent --check-permissions ./testdata/configs/debug.json
```

## HTTP API

В первом терминале запустить сервер:

```powershell
go run ./cmd/config-audit --http --addr :8080
```

Во втором терминале проверить healthcheck:

```powershell
curl.exe http://localhost:8080/healthz
```

Ожидаемый результат:

```json
{"status":"ok"}
```

Проверить анализ конфига:

```powershell
curl.exe -X POST "http://localhost:8080/v1/analyze?filename=config.json" --data-binary "@testdata/configs/debug.json"
```

Ответом будет JSON-массив найденных проблем.

## gRPC API

Контракт находится в:

```text
api/proto/audit/v1/audit.proto
```

Запуск сервера:

```powershell
go run ./cmd/config-audit --grpc --grpc-addr :9090
```

gRPC endpoint:

```proto
rpc Analyze(AnalyzeRequest) returns (AnalyzeResponse);
```

Go bindings находятся в `api/gen/audit/v1`, поэтому проект собирается без дополнительной генерации protobuf-кода.

## Архитектура

- `cmd/config-audit` - точка входа.
- `internal/cli` - CLI-флаги, stdin/stdout/stderr, exit codes.
- `internal/app` - основной сценарий анализа.
- `internal/config` - чтение, определение формата, JSON/YAML decoding.
- `internal/rules` - независимые правила безопасности.
- `internal/checker` - запуск набора правил.
- `internal/output` - текстовый и JSON-вывод.
- `internal/httpapi` - REST API.
- `internal/grpcapi` - gRPC API.
- `internal/dirscan` - рекурсивный поиск конфигов.
- `internal/filemode` - проверка прав через `os.Stat`.
