typedef signed long long int64;
typedef signed int int32;
typedef signed short int int16;
typedef signed char int8;
typedef unsigned long long uint64;
typedef unsigned int uint32;
typedef unsigned short int uint16;
typedef unsigned char uint8;

int
cas(uint32 volatile *ptr, uint32 old, uint32 new)
{
	return __sync_bool_compare_and_swap(ptr, old, new);
}

int
cas64(uint64 volatile *ptr, uint64 old, uint64 new)
{
	return __sync_bool_compare_and_swap(ptr, old, new);
}

uint32
xadd(uint32 volatile *ptr, int32 delta)
{
	return __sync_add_and_fetch(ptr, delta);
}

uint64
xadd64(uint64 volatile *ptr, int64 delta)
{
	return __sync_add_and_fetch(ptr, delta);
}

uint32
xchg(uint32 volatile *ptr, uint32 val)
{
	return __atomic_exchange_n(ptr, val, __ATOMIC_SEQ_CST);
}

uint64
xchg64(uint64 volatile *ptr, uint64 val)
{
	return __atomic_exchange_n(ptr, val, __ATOMIC_SEQ_CST);
}

void
store(uint32 volatile *ptr, uint32 val)
{
	__atomic_store_n(ptr, val, __ATOMIC_SEQ_CST);
}

void
store64(uint64 volatile *ptr, uint64 val)
{
	__atomic_store_n(ptr, val, __ATOMIC_SEQ_CST);
}

uint32
load(uint32 volatile *ptr)
{
	return __atomic_load_n(ptr, __ATOMIC_SEQ_CST);
}

uint64
load64(uint64 volatile *ptr)
{
	return __atomic_load_n(ptr, __ATOMIC_SEQ_CST);
}

void
or8(uint8 volatile *ptr, uint8 val)
{
	__sync_or_and_fetch(ptr, val);
}

void
and8(uint8 volatile *ptr, uint8 val)
{
	__sync_and_and_fetch(ptr, val);
}
