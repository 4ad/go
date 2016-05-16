// +build arm64

package runtime

func rtinit() {
	getg().m.curg = getg()
}
