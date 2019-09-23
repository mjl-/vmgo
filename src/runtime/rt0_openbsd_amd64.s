// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

#define HVT_HYPERCALL_PIO_BASE $0x500

#define HVT_HYPERCALL_WALLTIME	$1
#define HVT_HYPERCALL_PUTS	$2
#define HVT_HYPERCALL_POLL	$3
#define HVT_HYPERCALL_BLKINFO	$4
#define HVT_HYPERCALL_BLKWRITE	$5
#define HVT_HYPERCALL_BLKREAD	$6
#define HVT_HYPERCALL_NETINFO	$7
#define HVT_HYPERCALL_NETWRITE	$8
#define HVT_HYPERCALL_NETREAD	$9
#define HVT_HYPERCALL_HALT	$10

TEXT _rt0_amd64_openbsd(SB),NOSPLIT,$-8
	// read parameter from %RDI
	// PUSHQ 	RDI

	// do a hypercall, puts of a string.
	MOVW	HVT_HYPERCALL_PIO_BASE, DX
	ADDW	HVT_HYPERCALL_PUTS, DX
	LEAL	runtime·helloSolo5(SB), AX
	OUTL

	HLT

	JMP	_rt0_amd64(SB)

TEXT _rt0_amd64_openbsd_lib(SB),NOSPLIT,$0
	JMP	_rt0_amd64_lib(SB)
