// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// System calls and other sys.stuff for AMD64, OpenBSD
// /usr/src/sys/kern/syscalls.master for syscall numbers.
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// outl(dx uint32, ax uintptr)
// For making solo5 hypercalls.
TEXT runtime·outl(SB),NOSPLIT,$0-16
	MOVL	dx+0(FP), DX
	MOVQ	ax+8(FP), AX
	OUTL
	RET

// set tls base to DI
TEXT runtime·settls(SB),NOSPLIT,$0
	// adjust for ELF: wants to use -8(FS) for g
	ADDQ	$8, DI
	MOVQ	DI, AX

	// from gvisor/pkg/sentry/platform/ring0/lib_amd64.s
	MOVQ AX, DX
	SHRQ $32, DX
	MOVQ $0xc0000100, CX // MSR_FS_BASE
	BYTE $0x0f; BYTE $0x30;
	RET

// wrfsbase, not working in solo5 for me.
//	BYTE $0xf3; BYTE $0x48; BYTE $0x0f; BYTE $0xae; BYTE $0xd0;
//	RET
