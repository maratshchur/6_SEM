CC = gcc
CFLAGS = -Wall

TARGET = self_healing

all: $(TARGET)

$(TARGET): self_healing.c
	$(CC) $(CFLAGS) -o $@ $<

clean:
	rm -f $(TARGET) counter.txt

.PHONY: all clean