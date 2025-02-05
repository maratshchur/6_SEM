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
	echo "1111111\nabcdefg\n222222222" | ./$(TARGET) > output.txt
	@echo "Результат:"
	cat output.txt