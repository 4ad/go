// Put this data in a read-only section.
#define RODATA	8
// This data contains no pointers.
#define NOPTR	16

// Linker has a bug, and we need non-zero length symbols in
// these sections.

DATA type·runtime·moduledata(SB)/8, $224 // must match module size
GLOBL type·runtime·moduledata(SB), 0, $224

DATA data(SB)/4, $2
GLOBL data(SB), 0, $4

DATA rodata(SB)/4, $1
GLOBL rodata(SB), RODATA, $4

// .noptrdata
DATA noptrdata(SB)/4, $3
GLOBL noptrdata(SB), NOPTR, $4

// .bss
GLOBL bss(SB), 0, $4

// .noptrbss
GLOBL noptrbss(SB), NOPTR, $4
