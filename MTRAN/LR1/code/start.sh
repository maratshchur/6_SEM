#!/bin/bash

# # Создаем директорию build, если она не существует
# mkdir -p build

# # Копируем все .f90 файлы в директорию build
# find . -name "*.f90" -exec cp {} build \;

# # Переходим в директорию build
# cd build

# # Создаем Makefile
# cat > Makefile << EOF
# all: main

# main: main.o $(patsubst %.f90,%.o,$(wildcard *.f90))
# 	gfortran -o main main.o $(patsubst %.f90,%.o,$(wildcard *.f90))

# %.o: %.f90
# 	gfortran -c $< -o $@

# clean:
# 	rm -f *.o main
# EOF

# # Выполняем make
# make

# # Проверка на ошибки
# if [ $? -eq 0 ]; then
#   echo "Компиляция и сборка успешны!"
#   ./main
# else
#   echo "Ошибка компиляции или сборки!"
#   exit 1
# fi

# # Возвращаемся в исходную директорию
# cd ..

# echo "Скрипт завершен."