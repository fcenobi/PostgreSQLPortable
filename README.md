# PostgreSQL Portable Launcher for Windows

## Description
It allows you to control the portable version PostgreSQL from the tray icon.

## Building
```
go generate
go build -ldflags "-H=windowsgui"
```

## Current status
[x] Config file
[x] Download PostgreSQL distributive
[x] Extracting downloaded archive
[x] Show settings dialog on first start
[ ] Show downloading/extracting progress on the tray icon tooltip
[ ] Auto-install on startup
[ ] Check for updates

---

# Портабельный лаунчер PostgreSQL

## Описание
Позволяет управлять портабельной версией PostgreSQL из трея.

## Сборка
```
go generate
go build -ldflags "-H=windowsgui"
```

## Текущий статус
[x] Файл конфигурации
[x] Загрузка выбраного дистрибутива
[x] Распаковка загруженного архива
[x] Показ диалога настроек при первом запуске
[ ] Показ прогресса загрузки/распаковки в подсказке иконки в трее
[ ] Автоустановка при запуске
[ ] Проверка обновлений

## Исправление кодировки в psql консоли
Шрифты в консоли могут не корректно отображаться.
Для исправления этой проблемы рекомендую:
1. Изменить шрифт консоли по умолчанию на Lucida Console или любой поддерживающий Unicode
2. Cоздать/добавить в файл `%APPDATA%\postgresql\psqlrc.conf` следующие строки:
```
\! chcp 1251
SET client_encoding = 'UTF8'
```
