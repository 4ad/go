void bar(unsigned);

void
f1(void)
{
	unsigned u;

	for(u = 0; u <= 0; u--) {
		bar(u);
	}
}
