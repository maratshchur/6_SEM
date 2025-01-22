#!/bin/bash

# Запускаем программу в фоне
./self_healing &

# Получаем PID
pid=$!

# Цикл проверки
for i in {1..5}; do
    sleep 5  # Ждём 5 секунд
    kill -TERM $pid # Отправляем SIGTERM
    sleep 1 # Даём время на перезапуск
    counter=$(cat counter.txt)
    echo "Iteration $i: Counter = $counter, PID = $(pgrep self_healing)"
done

# Завершаем программу после тестов (замените <PID> на актуальный PID)
kill -TERM $(pgrep self_healing)