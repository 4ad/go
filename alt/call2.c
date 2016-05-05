#include <unistd.h>
#include <stdio.h>

int __attribute__ ((noinline))
foo3(void)
{
	int64_t buf[16];
	return buf[15];
}

int __attribute__ ((noinline))
foo2(int a0, int a1, int a2, int a3, int a4, int a5, int b6, int b7)
{
	return 2;
}

int __attribute__ ((noinline))
foo1(int a0, int a1, int a2, int a3, int a4, int a5, int b6, int b7)
{
	return b6+b7;
}

int __attribute__ ((noinline))
foo0(int a0, int a1, int a2, int a3, int a4, int a5, int b6, int b7)
{
	return foo1(0, 1, 2, 3, 4, 5, 6, 7) + foo2(0, 1, 2, 3, 4, 5, 6, 7) + foo3();
}

int
main()
{
	return foo0(0, 1, 2, 3, 4, 5, 6, 7);
}
