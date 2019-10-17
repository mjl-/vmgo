#include "textflag.h"
#include "funcdata.h"

// outl(dx uint32, ax uintptr)
// For making solo5 hypercalls.
TEXT syscall·outl(SB),NOSPLIT,$0-16
	MOVL    dx+0(FP), DX
	MOVQ    ax+8(FP), AX
	OUTL
	RET
