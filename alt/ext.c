#include <unistd.h>

char hellostr[] = "Hello, world!\n";

void
foo(ssize_t f(int, const void *, size_t))
{
	f(2, hellostr, sizeof hellostr);
}

int
main()
{
	foo(write);
	return 0;
}
