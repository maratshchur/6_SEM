#include "reverse.h"
#include <stdlib.h>

void reverse_stream(FILE *input, FILE *output) {
    const size_t BUFFER_SIZE = 1024 * 1024;
    char *buffer = (char *)malloc(BUFFER_SIZE);

    if (!buffer) {
        fprintf(stderr, "Ошибка: недостаточно памяти.\n");
        return;
    }

    // Читаем весь поток
    size_t bytes_read = fread(buffer, 1, BUFFER_SIZE, input);
    if (ferror(input)) {
        fprintf(stderr, "Ошибка чтения входного потока.\n");
        free(buffer);
        return;
    }

    // Инвертируем порядок байт
    for (size_t i = 0; i < bytes_read / 2; i++) {
        char temp = buffer[i];
        buffer[i] = buffer[bytes_read - i - 1];
        buffer[bytes_read - i - 1] = temp;
    }

    // Пишем инвертированный поток в выходной файл
    fwrite(buffer, 1, bytes_read, output);
    if (ferror(output)) {
        fprintf(stderr, "Ошибка записи в выходной поток.\n");
    }

    free(buffer);
}