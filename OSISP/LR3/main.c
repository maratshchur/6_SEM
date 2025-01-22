#include <stdio.h>
#include <stdlib.h>
#include "reverse.h"

int main(int argc, char *argv[]) {
    FILE *input = stdin;  // Стандартный ввод по умолчанию
    FILE *output = stdout; // Стандартный вывод по умолчанию

    // Обработка аргументов командной строки
    if (argc > 1) {
        input = fopen(argv[1], "rb");
        if (!input) {
            fprintf(stderr, "Ошибка: не удалось открыть файл %s для чтения.\n", argv[1]);
            return 1;
        }
    }

    if (argc > 2) {
        output = fopen(argv[2], "wb");
        if (!output) {
            fprintf(stderr, "Ошибка: не удалось открыть файл %s для записи.\n", argv[2]);
            fclose(input);
            return 1;
        }
    }

    // Инвертируем поток
    reverse_stream(input, output);

    // Закрываем файлы
    if (input != stdin) fclose(input);
    if (output != stdout) fclose(output);

    return 0;
}