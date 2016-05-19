#include <unistd.h>

double __attribute__ ((noinline))
int8_to_float64(int8_t v)
{
	return (double)v;
}

int
main()
{
	int8_t i8 = 42;
	volatile double f64;

	f64 = int8_to_float64(i8);

	return (int)f64;
}
