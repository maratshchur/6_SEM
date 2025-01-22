# Компилятор и флаги
CC = gcc
CFLAGS = -Wall -Wextra -std=c11

# Исполняемый файл
TARGET = reverse_filter

# Список исходников и объектных файлов
SOURCES = main.c reverse.c
OBJECTS = $(SOURCES:.c=.o)

# Сборка исполняемого файла
$(TARGET): $(OBJECTS)
	$(CC) $(CFLAGS) -o $(TARGET) $(OBJECTS)

# Компиляция объектных файлов
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Очистка
clean:
	rm -f $(OBJECTS) $(TARGET)

# Тестирование
test: $(TARGET)
	@echo "Тестирование..."
	echo "123456789\nabcdefg\n987654321" | ./$(TARGET) > output.txt
	@echo "Результат:"
	cat output.txt