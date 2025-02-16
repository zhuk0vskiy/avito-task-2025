#!/bin/bash

# Путь к файлу конфигурации
CONFIG_FILE="benchmark.conf"

# Читаем параметры из файла
while IFS='=' read -r key value; do
  case "$key" in
    "URL") URL=$(echo "$value" | xargs) ;;
    "REQUESTS") REQUESTS=$(echo "$value" | xargs) ;;
    "CONCURRENCY") CONCURRENCY=$(echo "$value" | xargs) ;;
    "KEEP_ALIVE") KEEP_ALIVE=$(echo "$value" | xargs) ;;
    "HEADERS") HEADERS=$(echo "$value" | xargs) ;;
    "STATISTICS") STATISTICS=$(echo "$value" | xargs) ;;
  esac
done < "$CONFIG_FILE"

# Проверяем, что REQUESTS и CONCURRENCY не равны нулю
if [ "$REQUESTS" -le 0 ] || [ "$CONCURRENCY" -le 0 ]; then
  echo "Ошибка: Количество запросов и уровень нагрузки должны быть больше нуля."
  exit 1
fi

# Формируем команду для ab
COMMAND="ab -n $REQUESTS -c $CONCURRENCY"

# Добавляем опции
if [ "$KEEP_ALIVE" = "yes" ]; then
  COMMAND="$COMMAND -k"
fi

if [ "$HEADERS" = "yes" ]; then
  COMMAND="$COMMAND -i"
fi

if [ "$STATISTICS" = "yes" ]; then
  COMMAND="$COMMAND -s"
fi

# Добавляем тип контента и файл с телом запроса
COMMAND="$COMMAND -T 'application/json' -p payload.json"

# Запускаем команду
echo "Запуск тестирования с параметрами:"
echo "$COMMAND $URL"
$COMMAND $URL