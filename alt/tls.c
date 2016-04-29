#include <stdio.h>
#include <unistd.h>
#include <pthread.h>

__thread int64_t g;

void *
f(void *p)
{
	printf("f: %p\n", &g);
	return &g;
}

int
main()
{
	pthread_t t0, t1;
	void *p0, *p1;
	if(pthread_create(&t0, NULL, f, NULL)) {
		perror("create 0");
	}
	if(pthread_create(&t1, NULL, f, NULL)) {
		perror("create 1");
	}
	if(pthread_join(t0, &p0)) {
		perror("join 0");
	}
	if(pthread_join(t1, &p1)) {
		perror("join 1");
	}
	printf("g[0]=%p\n", p0);
	printf("g[1]=%p\n", p1);

	return 0;
}
