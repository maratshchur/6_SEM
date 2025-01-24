# Имя исполняемого файла
TARGET = main

# Компилятор и флаги
CC = gcc
CFLAGS = -Wall -Wextra -pthread -O2

# Исходные файлы
SRC = main.c

# Объектные файлы (создаются из SRC)
OBJ = $(SRC:.c=.o)

# Правило сборки исполняемого файла
$(TARGET): $(OBJ)
	$(CC) $(CFLAGS) -o $(TARGET) $(OBJ)

# Правило сборки объектных файлов
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Правило очистки
clean:
	rm -f $(OBJ) $(TARGET)

# Правило для запуска программы
run: $(TARGET)
	./$(TARGET)