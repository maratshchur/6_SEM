#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <time.h>

typedef struct {
    int *array;
    int start;
    int end;
} ThreadData;

void *sort_part(void *arg) {
    ThreadData *data = (ThreadData *)arg;
    int *array = data->array;
    int start = data->start;
    int end = data->end;

    for (int i = start; i < end; i++) {
        for (int j = start; j < end - (i - start) - 1; j++) {
            if (array[j] > array[j + 1]) {
                int temp = array[j];
                array[j] = array[j + 1];
                array[j + 1] = temp;
            }
        }
    }

    pthread_exit(NULL);
}

void merge(int *array, int start, int mid, int end) {
    int size1 = mid - start + 1;
    int size2 = end - mid;

    int *left = (int *)malloc(size1 * sizeof(int));
    int *right = (int *)malloc(size2 * sizeof(int));

    for (int i = 0; i < size1; i++) left[i] = array[start + i];
    for (int i = 0; i < size2; i++) right[i] = array[mid + 1 + i];

    int i = 0, j = 0, k = start;
    while (i < size1 && j < size2) {
        if (left[i] <= right[j]) {
            array[k++] = left[i++];
        } else {
            array[k++] = right[j++];
        }
    }

    while (i < size1) array[k++] = left[i++];
    while (j < size2) array[k++] = right[j++];

    free(left);
    free(right);
}

int main() {
    int size, num_threads;
    clock_t start_time, end_time;

    // Read array size and number of threads
    printf("Enter array size: ");
    scanf("%d", &size);
    printf("Enter number of threads: ");
    scanf("%d", &num_threads);

    if (num_threads > size) {
        printf("Number of threads cannot exceed array size.\n");
        return -1;
    }

    int *array = (int *)malloc(size * sizeof(int));
    srand(time(NULL));

    // Fill the array with random numbers
    for (int i = 0; i < size; i++) {
        array[i] = rand() % 100000;
    }

    pthread_t threads[num_threads];
    ThreadData thread_data[num_threads];

    int part_size = size / num_threads;

    start_time = clock(); // Start measuring time

    // Create threads for sorting parts of the array
    for (int i = 0; i < num_threads; i++) {
        thread_data[i].array = array;
        thread_data[i].start = i * part_size;
        thread_data[i].end = (i == num_threads - 1) ? size : thread_data[i].start + part_size;

        pthread_create(&threads[i], NULL, sort_part, &thread_data[i]);
    }

    // Wait for all threads to complete
    for (int i = 0; i < num_threads; i++) {
        pthread_join(threads[i], NULL);
    }

    // Merge sorted parts
    for (int i = 1; i < num_threads; i++) {
        merge(array, 0, thread_data[i - 1].end - 1, thread_data[i].end - 1);
    }

    end_time = clock(); // Stop measuring time

    // Write sorted array to file
    FILE *file = fopen("sorted_array.txt", "w");
    if (file == NULL) {
        perror("Error opening file");
        free(array);
        return -1;
    }

    for (int i = 0; i < size; i++) {
        fprintf(file, "%d\n", array[i]);
    }

    fclose(file); // Close the file

    double elapsed_time = (double)(end_time - start_time) / CLOCKS_PER_SEC;
    printf("Sorting time: %f seconds\n", elapsed_time);

    free(array);
    return 0;
}

//   ./main