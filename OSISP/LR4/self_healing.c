// self_healing.c
#include <stdio.h>
#include <unistd.h>
#include <signal.h>
#include <sys/types.h>
#include <stdlib.h>
#include <time.h>

int counter = 0;

void handle_signal(int sig) {
    pid_t pid = fork();

    if (pid == 0) {
        printf("Дочерний процесс продолжает работу. PID: %d\n", getpid());
        signal(sig, handle_signal); 
        return;
    } else if (pid > 0) {
        printf("Родительский процесс завершается. PID: %d\n", getpid());
        exit(0);
    } else {
        perror("Ошибка fork");
        exit(1);
    }
}

int main() {
    signal(SIGINT, handle_signal);
    signal(SIGTERM, handle_signal); 

    while (1) {
        counter++;
        printf("Счетчик: %d, PID: %d\n", counter, getpid());
        sleep(1);
        FILE *f = fopen("counter.txt", "w");
        if (f == NULL) {
            perror("Ошибка открытия файла");
            exit(1);
        }
        fprintf(f, "%d\n", counter);
        fclose(f);
    }

    return 0;
}

// ./self_healing