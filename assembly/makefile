all: gao.o cao.o
	gcc -m32 gao.o cao.o -o all -I.

gao.o: gao.c
	gcc -m32 -c gao.c

cao.o: cao.S
	gcc -m32 -c cao.S

clean:
	rm *.o all
