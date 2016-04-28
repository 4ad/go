#include <unistd.h>

extern void bar(void);

int
foo(void* one, void* two, void* three)
{
	uintptr_t sav[16];
	bar();
}
