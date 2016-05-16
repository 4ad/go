package runtime

func rtinit() {
	getg().m.curg = getg()
}
