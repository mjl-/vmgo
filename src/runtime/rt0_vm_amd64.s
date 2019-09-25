// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// HVT_HYPERCALL_PIO_BASE $0x500

#define HVT_HYPERCALL_WALLTIME	0x501
#define HVT_HYPERCALL_PUTS	0x502
#define HVT_HYPERCALL_POLL	0x503
#define HVT_HYPERCALL_BLKWRITE	0x504
#define HVT_HYPERCALL_BLKREAD	0x505
#define HVT_HYPERCALL_NETWRITE	0x506
#define HVT_HYPERCALL_NETREAD	0x507
#define HVT_HYPERCALL_HALT	0x508

TEXT _rt0_amd64_vm(SB),NOSPLIT,$-8
	// Do a hypercall, puts of a string.
	// MOVW	$HVT_HYPERCALL_PUTS, DX
	// MOVQ	$runtime·helloSolo5(SB), AX
	// OUTL

	PUSHQ	DI	// *bootInfo
	CALL	runtime·solo5init(SB)
	SUBQ	$8, SP

	// Solo5 has info in bootInfo, no argc/argv.
	MOVQ	$0, DI	// argc
	MOVQ	$0, SI	// argv
	JMP	runtime·rt0_go(SB)

// outl(dx uint32, ax uintptr)
// For making hypercalls to solo5.
TEXT runtime·outl(SB),NOSPLIT,$0-16
	MOVL	dx+0(FP), DX
	MOVQ	ax+8(FP), AX
	OUTL
	RET
