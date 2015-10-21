typedef signed long long int64_t;
typedef signed int int32_t;
typedef signed short int int16_t;

int16_t
foo(int64_t x, int64_t y)
{
	return (int16_t)y;
}

int32_t
bar(int64_t x, int64_t y)
{
	return (int32_t)y;
}
