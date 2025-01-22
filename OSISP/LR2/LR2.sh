#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Использование: $0 <input_file>"
  exit 1
fi

input_file="$1"

if [ ! -f "$input_file" ]; then
  echo "Ошибка: файл '$input_file' не найден."
  exit 1
fi

sed -E 's/(^|[.!?] *)([a-zа-яё])/\1\u\2/g' "$input_file"