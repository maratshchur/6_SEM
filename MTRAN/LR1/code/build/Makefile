all: main

main: main.o 
	gfortran -o main main.o 

%.o: %.f90
	gfortran -c $< -o 

clean:
	rm -f *.o main
