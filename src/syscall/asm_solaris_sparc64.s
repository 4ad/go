// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

//
// System calls for solaris/amd64 are implemented in ../runtime/syscall_solaris.go
//

TEXT ·sysvicall6(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_sysvicall6(SB)

TEXT ·rawSysvicall6(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_rawsysvicall6(SB)

TEXT ·chdir(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_chdir(SB)

TEXT ·chroot1(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_chroot(SB)

TEXT ·close(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_close(SB)

TEXT ·execve(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_execve(SB)

TEXT ·exit(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_exit(SB)

TEXT ·fcntl1(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_fcntl(SB)

TEXT ·forkx(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_forkx(SB)

TEXT ·gethostname(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_gethostname(SB)

TEXT ·getpid(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_getpid(SB)

TEXT ·ioctl(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_ioctl(SB)

TEXT ·pipe(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_pipe(SB)

TEXT ·RawSyscall(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_rawsyscall(SB)

TEXT ·setgid(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_setgid(SB)

TEXT ·setgroups1(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_setgroups(SB)

TEXT ·setsid(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_setsid(SB)

TEXT ·setuid(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_setuid(SB)

TEXT ·setpgid(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_setpgid(SB)

TEXT ·Syscall(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_syscall(SB)

TEXT ·wait4(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_wait4(SB)

TEXT ·write1(SB),NOSPLIT|NOFRAME,$0
	JMP	runtime·syscall_write(SB)
