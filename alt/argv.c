#include <stdio.h>

char __attribute__ ((noinline))
foo(int argc, char **argv)
{
	return argv[argc-1][0];
}

int
main(int argc, char **argv)
{
	printf("%c\n", foo(argc, argv));
	return 0;
}
