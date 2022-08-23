// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

var testExecKnownFails = map[string]struct{}{
	// ====================================================================
	// Compiles and builds but fails at execution.

	// --------------------------------------------------------------------
	// Won't fix
	//
	// Needs real long double support.
	`assets/github.com/vnmakarov/mir/c-tests/lacc/long-double-load.c`: {}, // EXEC FAIL
	// --------------------------------------------------------------------

	//TODO?
	`assets/tcc-0.9.27/tests/tests2/92_enum_bitfield.c`: {}, // EXEC FAIL

	// ====================================================================
	// Compiles but does not build.

	// --------------------------------------------------------------------
	// Won't fix
	//
	// Wrong prototype argument type.
	`assets/github.com/vnmakarov/mir/c-tests/new/var-size-in-var-initializer.c`: {}, // BUILD FAIL
	// --------------------------------------------------------------------

	// goto/label
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030909-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040704-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20111208-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-6.c`:                   {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920728-1.c`:                   {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950221-1.c`:                   {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr17078-1.c`:                  {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38051.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43269.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr77766.c`:                    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030909-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040704-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20111208-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-6.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920728-1.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950221-1.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr17078-1.c`:  {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38051.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43269.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr77766.c`:    {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0010-goto1.c`:             {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/goto.c`:                               {}, // BUILD FAIL
	`assets/tcc-0.9.27/tests/tests2/54_goto.c`:                                          {}, // BUILD FAIL

	// Long double constant overflows floa64.
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960405-1.c`:                 {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960405-1.c`: {}, // BUILD FAIL

	//TODO invalid void expr
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040629-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040705-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040705-2.c`:                 {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040629-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040705-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040705-2.c`: {}, // BUILD FAIL

	//TODO missed assignment expression value
	`assets/github.com/vnmakarov/mir/c-tests/lacc/logical-operators-basic.c`: {}, // BUILD FAIL

	// ====================================================================
	// Does not compile (transpile).

	// assembler
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20001009-2.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030222-1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050203-1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20061031-1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20061220-1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071211-1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071220-1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080122-1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071220-2.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960312-1.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990130-1.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990413-2.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-5.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr50310.c`:                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38533.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr40022.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr40657.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr41239.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43560.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43385.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr44852.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr45695.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr46309.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr47925.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49279.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49390.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51877.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52286.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51933.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51581-2.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51581-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56205.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56866.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-3.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-4.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56982.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-2.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58277-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58419.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64242.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr63641.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65053-2.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65053-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65648.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65956.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68328.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr69320-2.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr69691.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78438.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78726.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79354.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79737-2.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr80421.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr81588.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82954.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84524.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85156.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84478.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85756.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr88904.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stkalign.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20001009-2.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030222-1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050203-1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20061031-1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20061220-1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071211-1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071220-1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071220-2.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080122-1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960312-1.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990130-1.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990413-2.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-5.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr50310.c`: {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38533.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr40022.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr40657.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr41239.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43385.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43560.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr44852.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr45695.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr46309.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr47925.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49279.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49390.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51877.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51581-2.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51933.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51581-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52286.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56205.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56866.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56982.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-2.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-3.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-4.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58277-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58419.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr63641.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65053-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65053-2.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65648.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65956.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68328.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr69691.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr69320-2.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78438.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78726.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79354.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79737-2.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr80421.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr81588.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82954.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84524.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85156.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84478.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85756.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr88904.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93945.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94130.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stkalign.c`:     {}, // COMPILE FAIL

	// void func(void) __attribute__((aligned(256))) etc.
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/align-3.c`:                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr23467.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/align-3.c`: {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr23467.c`: {}, // COMPILE FAIL

	// uses signal(2)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20101011-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-1.c`:                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-2.c`:                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-3.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20101011-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-1.c`: {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-2.c`: {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-3.c`: {}, // COMPILE FAIL

	// VLA
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr41935.c`:                       {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr41935.c`:       {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-5.c`:                                {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-2.c`:                                {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-7.c`:                                {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fannkuchredux-5.c`:                        {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fannkuchredux.c`:                          {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-9.c`:                           {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/spectral-norm.c`:                          {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010209-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040308-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040811-1.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920721-2.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920929-1.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921017-1.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/970217-1.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/align-nest.c`:                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-2.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43220.c`:                       {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr77767.c`:                       {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/vla-dealloc-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010209-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040308-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040811-1.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920721-2.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920929-1.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921017-1.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/970217-1.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/align-nest.c`:    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-2.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43220.c`:       {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr77767.c`:       {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/vla-dealloc-1.c`: {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/78_vla_label.c`:                                        {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/79_vla_continue.c`:                                     {}, // COMPILE FAIL

	// vector
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-4.c`:                            {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/nbody-4.c`:                                 {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/nbody-8.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050316-2.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050316-1.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050316-3.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050604-1.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050607-1.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20060420-1.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr72824-2.c`:                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr23135.c`:                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr53645-2.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr53645.c`:                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60960.c`:                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65427.c`:                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70903.c`:                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71626-1.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71626-2.c`:                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85169.c`:                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85331.c`:                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/scal-to-vec1.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/scal-to-vec2.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-4.c`:                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-2.c`:                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-1.c`:                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/scal-to-vec3.c`:                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-5.c`:                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-6.c`:                         {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050316-1.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050316-2.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050316-3.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050607-1.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050604-1.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20060420-1.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr72824-2.c`: {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr23135.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr53645.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr53645-2.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60960.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65427.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70903.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71626-1.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71626-2.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85331.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85169.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr92618.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94412.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94524-1.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94524-2.c`:      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94591.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/scal-to-vec2.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/scal-to-vec3.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/scal-to-vec1.c`:   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-1.c`:         {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-2.c`:         {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-6.c`:         {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-5.c`:         {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-4.c`:         {}, // COMPILE FAIL

	//TODO longjmp/setjmp
	`assets/github.com/vnmakarov/mir/c-benchmarks/except.c`: {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/setjmp.c`:  {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/setjmp2.c`: {}, // COMPILE FAIL

	//TODO
	`assets/CompCert-3.6/test/c/vmach.c`:                                                                {}, // COMPILE FAIL: 24.go:4223:62: TODO U+002B '+' (check.go:822:check: check.go:285:checkExpr: check.go:904:check:)
	`assets/CompCert-3.6/test/c/fannkuch.c`:                                                             {}, // COMPILE FAIL: TODO (expr.go:2182:primaryExpression: expr.go:426:expr0: expr.go:1973:assignmentExpression)
	`assets/CompCert-3.6/test/c/aes.c`:                                                                  {}, // COMPILE FAIL: 1.go:4858:28: TODO gc.PredefinedType (check.go:663:checkFn: check.go:285:checkExpr: check.go:1253:check:)
	`assets/benchmarksgame-team.pages.debian.net/fasta-4.c`:                                             {}, // COMPILE FAIL: fasta-4.c.go:4591:3: undefined reference to 'fwrite_unlocked' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-8.c`:                                        {}, // COMPILE FAIL: TODO vector (decl.go:346:initDeclarator: type.go:17:typedef: type.go:57:typ0)
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-2.c`:                                        {}, // COMPILE FAIL: 20.go:2030:12: undefined: Xatoi (check.go:1338:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:226:checkExprOrType: check.go:1338:check: check.go:205:symbolResolver:)
	`assets/benchmarksgame-team.pages.debian.net/reverse-complement-4.c`:                                {}, // COMPILE FAIL: reverse-complement-4.c.go:4071:38: undefined reference to '__builtin_memmove' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/benchmarksgame-team.pages.debian.net/reverse-complement-6.c`:                                {}, // COMPILE FAIL: 47.go:4977:27: TODO uint64 -> int64 (check.go:1729:checkStatement: check.go:1909:check: check.go:735:isAssignable:)
	`assets/ccgo/bug/csmith2.c`:                                                                         {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:232:typ0)
	`assets/ccgo/bug/bitfield.c`:                                                                        {}, // COMPILE FAIL: TODO bitfield (expr.go:1438:postfixExpression: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/ccgo/bug/union3.c`:                                                                          {}, // COMPILE FAIL: TODO (type.go:17:typedef: type.go:184:typ0: type.go:240:typ0)
	`assets/ccgo/bug/struct.c`:                                                                          {}, // COMPILE FAIL: TODO union {ceiling char; unused char} 1, ft struct outer {magic char; pad1 array of 3 char; union {ceiling char; unused char}; pad2 array of 3 char; owner char} 9 (init.go:211:initializerStruct: init.go:72:initializer: init.go:269:initializerUnion)
	`assets/ccgo/bug/union.c`:                                                                           {}, // COMPILE FAIL: TODO (expr.go:1438:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/ccgo/bug/union2.c`:                                                                          {}, // COMPILE FAIL: TODO (expr.go:1438:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/ccgo/bug/union4.c`:                                                                          {}, // COMPILE FAIL: TODO (type.go:17:typedef: type.go:184:typ0: type.go:240:typ0)
	`assets/ccgo/bug/struct2.c`:                                                                         {}, // COMPILE FAIL: TODO (expr.go:1438:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/ccgo/bug/sqlite.c`:                                                                          {}, // COMPILE FAIL: 33.go:352:82: TODO gc.PredefinedType (check.go:1304:check: check.go:226:checkExprOrType: check.go:1535:check:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000113-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000703-1.c`:                                 {}, // COMPILE FAIL: 20000703-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000801-3.c`:                                 {}, // COMPILE FAIL: 20000801-3.c:5:20: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000822-1.c`:                                 {}, // COMPILE FAIL: 20000822-1.c:12:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000917-1.c`:                                 {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ana)))" pointer to int3 exprUintptr -> pointer to int3 exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20001203-2.c`:                                 {}, // COMPILE FAIL: 20001203-2.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010122-1.c`:                                 {}, // COMPILE FAIL: 20010122-1.c.go:323:37: undefined reference to '__builtin_return_address' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010605-1.c`:                                 {}, // COMPILE FAIL: 20010605-1.c:5:14: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010605-2.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010904-1.c`:                                 {}, // COMPILE FAIL: 20010904-1.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010904-2.c`:                                 {}, // COMPILE FAIL: 20010904-2.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20011113-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020206-2.c`:                                 {}, // COMPILE FAIL: TODO "tnA{fda: (iqlibc.ppUint16FromInt32(((ccv1 << (iqlibc.ppInt32FromInt32(8))) + (ani << (iqlibc.ppInt32FromInt32(4)))))), }" A exprDefault -> A exprSelect (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020320-1.c`:                                 {}, // COMPILE FAIL: TODO "aav3" struct large {x int; y array of 9 int} exprDefault -> struct large {x int; y array of 9 int} exprSelect (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020314-1.c`:                                 {}, // COMPILE FAIL: 20020314-1.c.go:330:34: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020411-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020412-1.c`:                                 {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020418-1.c`:                                 {}, // COMPILE FAIL: 167.go:340:8: undefined: X__builtin_trap (check.go:1338:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:226:checkExprOrType: check.go:1338:check: check.go:205:symbolResolver:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030109-1.c`:                                 {}, // COMPILE FAIL: 20030109-1.c:5:1: incomplete type: array of int (type.go:184:typ0: type.go:40:typ0: type.go:311:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20021113-1.c`:                                 {}, // COMPILE FAIL: 20021113-1.c.go:331:36: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030501-1.c`:                                 {}, // COMPILE FAIL: 20030501-1.c:7:9: nested functions not supported (stmt.go:23:statement: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030330-1.c`:                                 {}, // COMPILE FAIL: 20030330-1.c.go:324:4: undefined reference to 'link_error' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030714-1.c`:                                 {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:448:expr0: expr.go:1311:postfixExpression: expr.go:1595:postfixExpressionSelect)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030323-1.c`:                                 {}, // COMPILE FAIL: 20030323-1.c.go:322:36: undefined reference to '__builtin_return_address' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030811-1.c`:                                 {}, // COMPILE FAIL: 20030811-1.c:9:5: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030910-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040307-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040302-1.c`:                                 {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040223-1.c`:                                 {}, // COMPILE FAIL: 20040223-1.c.go:2442:53: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040411-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionSizeofType (expr.go:67:expr: expr.go:456:expr0: expr.go:1080:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040423-1.c`:                                 {}, // COMPILE FAIL: 20040423-1.c:13:14: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040520-1.c`:                                 {}, // COMPILE FAIL: 20040520-1.c:6:13: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041124-1.c`:                                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex unsigned short _Complex unsigned short (type.go:331:structLiteral: type.go:184:typ0: type.go:122:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041201-1.c`:                                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex char _Complex char (type.go:17:typedef: type.go:184:typ0: type.go:122:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041214-1.c`:                                 {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041218-2.c`:                                 {}, // COMPILE FAIL: 20041218-2.c:7:10: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050121-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-2.c`:                                 {}, // COMPILE FAIL: 20040709-2.c.go:388:5: undefined reference to '__builtin_classify_type' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-1.c`:                                 {}, // COMPILE FAIL: 20040709-1.c.go:389:5: undefined reference to '__builtin_classify_type' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-3.c`:                                 {}, // COMPILE FAIL: 20040709-3.c.go:392:5: undefined reference to '__builtin_classify_type' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050613-1.c`:                                 {}, // COMPILE FAIL: 20050613-1.c:8:1: incomplete type: array of int (type.go:184:typ0: type.go:40:typ0: type.go:311:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20051012-1.c`:                                 {}, // COMPILE FAIL: 341.go:353:5: TODO *gc.Arguments (check.go:285:checkExpr: check.go:623:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070614-1.c`:                                 {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070919-1.c`:                                 {}, // COMPILE FAIL: 20070919-1.c:31:7: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071029-1.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:1438:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070824-1.c`:                                 {}, // COMPILE FAIL: 20070824-1.c.go:334:34: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071120-1.c`:                                 {}, // COMPILE FAIL: 20071120-1.c:13:41: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071210-1.c`:                                 {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080519-1.c`:                                 {}, // COMPILE FAIL: TODO exprIndex (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1044:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080502-1.c`:                                 {}, // COMPILE FAIL: 20080502-1.c.go:320:5: undefined reference to '__builtin_signbit' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20081117-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090113-2.c`:                                 {}, // COMPILE FAIL: 20090113-2.c:1:1: invalid type size: -1 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090113-3.c`:                                 {}, // COMPILE FAIL: 20090113-3.c:1:1: invalid type size: -1 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090219-1.c`:                                 {}, // COMPILE FAIL: 20090219-1.c:12:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20120111-1.c`:                                 {}, // COMPILE FAIL: 432.go:3531:79: undefined: arg (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:285:checkExpr: check.go:1215:check: check.go:205:symbolResolver:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20181120-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20180921-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920302-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920415-1.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920428-2.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (stmt.go:23:statement: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-3.c`:                                   {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-4.c`:                                   {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-5.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-7.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-1.c`:                                   {}, // COMPILE FAIL: 473.go:352:5: TODO *gc.Arguments (check.go:285:checkExpr: check.go:623:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920612-2.c`:                                   {}, // COMPILE FAIL: 920612-2.c:6:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920625-1.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920721-4.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920731-1.c`:                                   {}, // COMPILE FAIL: 920731-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920908-1.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921202-1.c`:                                   {}, // COMPILE FAIL: 524.go:352:3: TODO *gc.Arguments (check.go:1972:check: check.go:623:check: check.go:655:checkFn:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921215-1.c`:                                   {}, // COMPILE FAIL: 921215-1.c:5:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921208-2.c`:                                   {}, // COMPILE FAIL: 529.go:347:3: TODO *gc.Arguments (check.go:1972:check: check.go:623:check: check.go:655:checkFn:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930406-1.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (expr.go:2188:primaryExpression: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931002-1.c`:                                   {}, // COMPILE FAIL: 931002-1.c:10:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-10.c`:                                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930718-1.c`:                                   {}, // COMPILE FAIL: 561.go:360:26: undefined: tsrtx_def (check.go:1348:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:238:checkType: check.go:1348:check: check.go:205:symbolResolver:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-12.c`:                                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-14.c`:                                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-2.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-4.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-6.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-8.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931031-1.c`:                                   {}, // COMPILE FAIL: 590.go:363:86: value 8589934590 overflows uint32 (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/941202-1.c`:                                   {}, // COMPILE FAIL: 941202-1.c.go:328:36: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950628-1.c`:                                   {}, // COMPILE FAIL: TODO (expr.go:67:expr: expr.go:448:expr0: expr.go:1299:postfixExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950906-1.c`:                                   {}, // COMPILE FAIL: 950906-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960416-1.c`:                                   {}, // COMPILE FAIL: 960416-1.c:57:23: internal error: t_be (expr.go:2170:primaryExpression: expr.go:2427:primaryExpressionIntConst: type.go:27:helper)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/980526-1.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/980929-1.c`:                                   {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ann)))" pointer to int exprUintptr -> pointer to int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990208-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990222-1.c`:                                   {}, // COMPILE FAIL: TODO (expr.go:2182:primaryExpression: expr.go:426:expr0: expr.go:1973:assignmentExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990326-1.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991014-1.c`:                                   {}, // COMPILE FAIL: 991014-1.c:10:1: invalid type size: -496 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991118-1.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991112-1.c`:                                   {}, // COMPILE FAIL: 991112-1.c.go:328:5: undefined reference to '__builtin_isprint' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-2.c`:                                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:125:externalDeclaration: decl.go:274:declaration: decl.go:303:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-3.c`:                                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:125:externalDeclaration: decl.go:274:declaration: decl.go:303:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-4.c`:                                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:125:externalDeclaration: decl.go:274:declaration: decl.go:303:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/align-1.c`:                                    {}, // COMPILE FAIL: align-1.c:1:13: unsupported alignment 16 of new_int (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/anon-1.c`:                                     {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:448:expr0: expr.go:1311:postfixExpression: expr.go:1586:postfixExpressionSelect)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alloca-1.c`:                                   {}, // COMPILE FAIL: alloca-1.c.go:329:34: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf-sign-1.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf64-1.c`:                                     {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-4.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-3.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-6.c`:                                   {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-7.c`:                                   {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/builtin-types-compatible-p.c`:                 {}, // COMPILE FAIL: TODO *cc.InvalidType (expr.go:448:expr0: expr.go:1308:postfixExpression: expr.go:1722:postfixExpressionCall)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/built-in-setjmp.c`:                            {}, // COMPILE FAIL: built-in-setjmp.c.go:322:2: undefined reference to '__builtin_longjmp' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/builtin-bitops-1.c`:                           {}, // COMPILE FAIL: builtin-bitops-1.c.go:1581:52: undefined reference to '__builtin_ctzl' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/comp-goto-2.c`:                                {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/call-trap-1.c`:                                {}, // COMPILE FAIL: 925.go:357:2: TODO gc.PredefinedType, uintptr (check.go:1719:checkStatement: check.go:1972:check: check.go:597:check:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-2.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-1.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/comp-goto-1.c`:                                {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-5.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/compndlit-1.c`:                                {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-7.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-6.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ffs-2.c`:                                      {}, // COMPILE FAIL: ffs-2.c.go:331:6: undefined reference to '__builtin_ffs' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ffs-1.c`:                                      {}, // COMPILE FAIL: ffs-1.c.go:319:5: undefined reference to '__builtin_ffs' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/frame-address.c`:                              {}, // COMPILE FAIL: frame-address.c.go:360:39: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/fprintf-2.c`:                                  {}, // COMPILE FAIL: fprintf-2.c.go:4463:15: undefined reference to 'tmpnam' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/compare-fp-1.c`:                          {}, // COMPILE FAIL: compare-fp-1.c.go:332:26: undefined reference to '__builtin_islessgreater' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4.c`:                              {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4f.c`:                             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4l.c`:                             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-5.c`:                              {}, // COMPILE FAIL: fp-cmp-5.c.go:338:9: undefined reference to '__builtin_islessequal' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8l.c`:                             {}, // COMPILE FAIL: fp-cmp-8l.c.go:341:5: undefined reference to '__builtin_isless' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8.c`:                              {}, // COMPILE FAIL: fp-cmp-8.c.go:342:5: undefined reference to '__builtin_isless' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8f.c`:                             {}, // COMPILE FAIL: fp-cmp-8f.c.go:342:5: undefined reference to '__builtin_isless' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/mzero4.c`:                                {}, // COMPILE FAIL: mzero4.c.go:359:18: undefined reference to 'atanf' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr38016.c`:                               {}, // COMPILE FAIL: pr38016.c.go:364:5: undefined reference to '__builtin_islessequal' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/medce-1.c`:                                    {}, // COMPILE FAIL: medce-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/memchr-1.c`:                                   {}, // COMPILE FAIL: 1071.go:377:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nest-align-1.c`:                               {}, // COMPILE FAIL: nest-align-1.c:25:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nest-stdar-1.c`:                               {}, // COMPILE FAIL: nest-stdar-1.c:5:10: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-1.c`:                                 {}, // COMPILE FAIL: nestfunc-1.c:15:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-2.c`:                                 {}, // COMPILE FAIL: nestfunc-2.c:13:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-3.c`:                                 {}, // COMPILE FAIL: nestfunc-3.c:12:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-6.c`:                                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-5.c`:                                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-7.c`:                                 {}, // COMPILE FAIL: nestfunc-7.c:15:12: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr17377.c`:                                    {}, // COMPILE FAIL: pr17377.c.go:331:36: undefined reference to '__builtin_return_address' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr19449.c`:                                    {}, // COMPILE FAIL: pr19449.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-2.c`:                                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-2.c:10:24:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-1.c`:                                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-1.c:10:24:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-1.c`:                                  {}, // COMPILE FAIL: TODO (expr.go:448:expr0: expr.go:1248:postfixExpression: expr.go:1130:postfixExpressionIndex)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-3.c`:                                  {}, // COMPILE FAIL: pr22061-3.c:4:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-4.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionSizeofExpr (expr.go:67:expr: expr.go:456:expr0: expr.go:1067:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-3.c`:                                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: Xf(aatls), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-3.c:12:24:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr23324.c`:                                    {}, // COMPILE FAIL: pr23324.c:4:21: invalid type size: -1 (type.go:325:unionLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr24135.c`:                                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr28289.c`:                                    {}, // COMPILE FAIL: 1152.go:368:45: TODO *gc.Arguments (check.go:1972:check: check.go:623:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr30185.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:67:expr: expr.go:448:expr0: expr.go:1299:postfixExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr34154.c`:                                    {}, // COMPILE FAIL: TODO LabeledStatementRange (stmt.go:182:blockItem: stmt.go:21:statement: stmt.go:71:labeledStatement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr35456.c`:                                    {}, // COMPILE FAIL: pr35456.c.go:336:7: undefined reference to '__builtin_signbit' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr36321.c`:                                    {}, // COMPILE FAIL: pr36321.c.go:323:34: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr37780.c`:                                    {}, // COMPILE FAIL: pr37780.c.go:322:35: undefined reference to '__builtin_ctz' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39339.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:211:initializerStruct: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38969.c`:                                    {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38151.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:331:structLiteral: type.go:184:typ0: type.go:122:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39228.c`:                                    {}, // COMPILE FAIL: pr39228.c.go:331:9: undefined reference to '__builtin_isinf' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39100.c`:                                    {}, // COMPILE FAIL: 1231.go:356:6: undefined: tsC (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:226:checkExprOrType: check.go:1215:check: check.go:205:symbolResolver:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr42570.c`:                                    {}, // COMPILE FAIL: pr42570.c:2:9: invalid type size: -1 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr42614.c`:                                    {}, // COMPILE FAIL: TODO "siinlined_wrong" staticInternal (link.go:1046:print0: link.go:993:print0: link.go:949:name)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43783.c`:                                    {}, // COMPILE FAIL: pr43783.c:6:3: unsupported alignment 16 of UINT192 (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr44575.c`:                                    {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr47237.c`:                                    {}, // COMPILE FAIL: pr47237.c.go:334:68: undefined reference to '__builtin_apply_args' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49218.c`:                                    {}, // COMPILE FAIL: pr49218.c:2:18: unsupported alignment 16 of L (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49644.c`:                                    {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49768.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51447.c`:                                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52979-1.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52979-2.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr54471.c`:                                    {}, // COMPILE FAIL: pr54471.c:15:22: unsupported alignment 16 of unsigned __int128 (type.go:29:helper: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56837.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:35:typ: type.go:258:typ0: type.go:122:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57130.c`:                                    {}, // COMPILE FAIL: 1366.go:336:28: undefined: tsS (check.go:1348:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:238:checkType: check.go:1348:check: check.go:205:symbolResolver:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58385.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((((iqlibc.ppBoolInt32(((((iqlibc.ppInt32FromInt32(0)) != 0)) || ((Xa != 0)))) & iqlibc.ppBoolInt32((Xfoo(aatls) >= (iqlibc.ppInt32FromInt32(0))))) <= (iqlibc.ppInt32FromInt32(1)))) && (((iqlibc.ppInt32FromInt32(1)) != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58431.c`:                                    {}, // COMPILE FAIL: 1390.go:351:6: TODO int32 -> int16 (check.go:1729:checkStatement: check.go:1909:check: check.go:735:isAssignable:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58726.c`:                                    {}, // COMPILE FAIL: TODO bitfield (decl.go:376:initDeclarator: type.go:35:typ: type.go:210:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58984.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60017.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60003.c`:                                    {}, // COMPILE FAIL: pr60003.c.go:335:5: undefined reference to '__builtin_setjmp' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr63302.c`:                                    {}, // COMPILE FAIL: TODO "(-((iqlibc.ppInt32FromInt32(1))))" int exprDefault -> __int128 exprDefault (pr63302.c:16:33:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr61375.c`:                                    {}, // COMPILE FAIL: pr61375.c:17:29: unsupported alignment 16 of __int128 (type.go:29:helper: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64006.c`:                                    {}, // COMPILE FAIL: pr64006.c.go:331:6: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr61725.c`:                                    {}, // COMPILE FAIL: pr61725.c.go:321:9: undefined reference to '__builtin_ffs' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64756.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = (((Xd != 0)) || ((Xd != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65170.c`:                                    {}, // COMPILE FAIL: pr65170.c:4:27: unsupported alignment 16 of V (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65215-3.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65369.c`:                                    {}, // COMPILE FAIL: pr65369.c:42:31: unsupported alignment 16 of array of 97 unsigned char (type.go:35:typ: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr66556.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr67037.c`:                                    {}, // COMPILE FAIL: 1463.go:382:2: TODO *gc.Arguments (check.go:1972:check: check.go:626:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68321.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = (((Xu != 0)) || ((Xn != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68185.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv3) && ((ccv2 != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68249.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:432:expr0: expr.go:637:conditionalExpression: expr.go:51:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68381.c`:                                    {}, // COMPILE FAIL: pr68381.c.go:321:5: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70127.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70460.c`:                                    {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70602.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71494.c`:                                    {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71700.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78170.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71554.c`:                                    {}, // COMPILE FAIL: pr71554.c.go:336:5: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79286.c`:                                    {}, // COMPILE FAIL: pr79286.c:2:21: incomplete type: array of array of 8 int (type.go:35:typ: type.go:40:typ0: type.go:311:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79737-1.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr80692.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Decimal64 _Decimal64 (expr.go:2436:primaryExpressionFloatConst: type.go:29:helper: type.go:122:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82210.c`:                                    {}, // COMPILE FAIL: pr82210.c:18:9: unsupported alignment 16 of struct S {a array of struct T; b array of int} (type.go:35:typ: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82387.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv2) && ((ccv1 != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84169.c`:                                    {}, // COMPILE FAIL: pr84169.c:4:27: unsupported alignment 16 of T (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84748.c`:                                    {}, // COMPILE FAIL: pr84748.c:3:27: unsupported alignment 16 of u128 (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85529-1.c`:                                  {}, // COMPILE FAIL: TODO "pp_ = ((Xs.fda) != iqlibc.ppBoolInt32(((ccv3) && ((ccv1 != 0)))))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85582-2.c`:                                  {}, // COMPILE FAIL: pr85582-2.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85582-3.c`:                                  {}, // COMPILE FAIL: pr85582-3.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85095.c`:                                    {}, // COMPILE FAIL: pr85095.c.go:322:33: undefined reference to '__builtin_add_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr87053.c`:                                    {}, // COMPILE FAIL: TODO (init.go:267:initializerUnion: type.go:35:typ: type.go:217:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr89195.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr86528.c`:                                    {}, // COMPILE FAIL: pr86528.c.go:320:36: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr89434.c`:                                    {}, // COMPILE FAIL: pr89434.c.go:324:2: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr88714.c`:                                    {}, // COMPILE FAIL: 1585.go:350:44: undefined: tsT (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:226:checkExprOrType: check.go:1215:check: check.go:205:symbolResolver:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/printf-2.c`:                                   {}, // COMPILE FAIL: printf-2.c.go:4479:8: undefined reference to 'freopen' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stdarg-3.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strct-stdarg-1.c`:                             {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strct-varg-1.c`:                               {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stdarg-1.c`:                                   {}, // COMPILE FAIL: stdarg-1.c.go:450:2: undefined reference to '__builtin_va_copy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/string-opt-18.c`:                              {}, // COMPILE FAIL: string-opt-18.c.go:332:5: undefined reference to 'mempcpy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-3.c`:                                   {}, // COMPILE FAIL: 1635.go:377:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/struct-ini-2.c`:                               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/struct-ini-3.c`:                               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-6.c`:                                   {}, // COMPILE FAIL: 1638.go:338:21: TODO *gc.BasicLit (check.go:1013:check: check.go:1082:checkArray: check.go:1130:checkArrayElem:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-2.c`:                                   {}, // COMPILE FAIL: 1634.go:1150:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-4.c`:                                   {}, // COMPILE FAIL: 1636.go:377:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-14.c`:                                  {}, // COMPILE FAIL: va-arg-14.c.go:374:2: undefined reference to '__builtin_va_copy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-5.c`:                                   {}, // COMPILE FAIL: TODO (init.go:18:initializerOuter: init.go:72:initializer: init.go:236:initializerUnion)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-10.c`:                                  {}, // COMPILE FAIL: va-arg-10.c.go:384:2: undefined reference to '__builtin_va_copy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/user-printf.c`:                                {}, // COMPILE FAIL: user-printf.c.go:4516:8: undefined reference to 'freopen' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-22.c`:                                  {}, // COMPILE FAIL: va-arg-22.c:24:1: invalid type size: 0 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-2.c`:                                   {}, // COMPILE FAIL: 1666.go:409:5: TODO *gc.Arguments (check.go:285:checkExpr: check.go:623:check: check.go:655:checkFn:)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-pack-1.c`:                              {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zero-struct-1.c`:                              {}, // COMPILE FAIL: zero-struct-1.c:1:1: invalid type size: -1 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zero-struct-2.c`:                              {}, // COMPILE FAIL: zero-struct-2.c:3:19: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zerolen-1.c`:                                  {}, // COMPILE FAIL: 1699.go:351:25: TODO *gc.StructType struct{Fname_len [1]uint8} (check.go:1444:check: check.go:226:checkExprOrType: check.go:1393:check:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/AbsInt/CompCert/test/c/vmach.c`:                                                  {}, // COMPILE FAIL: 24.go:4223:62: TODO U+002B '+' (check.go:822:check: check.go:285:checkExpr: check.go:904:check:)
	`assets/github.com/AbsInt/CompCert/test/c/fannkuch.c`:                                               {}, // COMPILE FAIL: TODO (expr.go:2182:primaryExpression: expr.go:426:expr0: expr.go:1973:assignmentExpression)
	`assets/github.com/AbsInt/CompCert/test/c/aes.c`:                                                    {}, // COMPILE FAIL: 1.go:4858:28: TODO gc.PredefinedType (check.go:663:checkFn: check.go:285:checkExpr: check.go:1253:check:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000113-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000703-1.c`:                 {}, // COMPILE FAIL: gcc.c-torture/execute/20000703-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000801-3.c`:                 {}, // COMPILE FAIL: 20000801-3.c:5:20: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000917-1.c`:                 {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ana)))" pointer to int3 exprUintptr -> pointer to int3 exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000822-1.c`:                 {}, // COMPILE FAIL: 20000822-1.c:12:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20001203-2.c`:                 {}, // COMPILE FAIL: gcc.c-torture/execute/20001203-2.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010122-1.c`:                 {}, // COMPILE FAIL: 20010122-1.c.go:323:37: undefined reference to '__builtin_return_address' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010605-1.c`:                 {}, // COMPILE FAIL: 20010605-1.c:5:14: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010605-2.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010904-1.c`:                 {}, // COMPILE FAIL: 20010904-1.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010904-2.c`:                 {}, // COMPILE FAIL: 20010904-2.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20011113-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020206-2.c`:                 {}, // COMPILE FAIL: TODO "tnA{fda: (iqlibc.ppUint16FromInt32(((ccv1 << (iqlibc.ppInt32FromInt32(8))) + (ani << (iqlibc.ppInt32FromInt32(4)))))), }" A exprDefault -> A exprSelect (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020320-1.c`:                 {}, // COMPILE FAIL: TODO "aav3" struct large {x int; y array of 9 int} exprDefault -> struct large {x int; y array of 9 int} exprSelect (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020314-1.c`:                 {}, // COMPILE FAIL: 20020314-1.c.go:330:34: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020411-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020412-1.c`:                 {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020418-1.c`:                 {}, // COMPILE FAIL: 167.go:340:8: undefined: X__builtin_trap (check.go:1338:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:226:checkExprOrType: check.go:1338:check: check.go:205:symbolResolver:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20021113-1.c`:                 {}, // COMPILE FAIL: 20021113-1.c.go:331:36: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030109-1.c`:                 {}, // COMPILE FAIL: 20030109-1.c:5:1: incomplete type: array of int (type.go:184:typ0: type.go:40:typ0: type.go:311:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030501-1.c`:                 {}, // COMPILE FAIL: 20030501-1.c:7:9: nested functions not supported (stmt.go:23:statement: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030330-1.c`:                 {}, // COMPILE FAIL: 20030330-1.c.go:324:4: undefined reference to 'link_error' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030323-1.c`:                 {}, // COMPILE FAIL: 20030323-1.c.go:322:36: undefined reference to '__builtin_return_address' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030714-1.c`:                 {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:448:expr0: expr.go:1311:postfixExpression: expr.go:1595:postfixExpressionSelect)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030811-1.c`:                 {}, // COMPILE FAIL: 20030811-1.c:9:5: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030910-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040307-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040302-1.c`:                 {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040223-1.c`:                 {}, // COMPILE FAIL: 20040223-1.c.go:2442:53: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040411-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionSizeofType (expr.go:67:expr: expr.go:456:expr0: expr.go:1080:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040520-1.c`:                 {}, // COMPILE FAIL: 20040520-1.c:6:13: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040423-1.c`:                 {}, // COMPILE FAIL: 20040423-1.c:13:14: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041124-1.c`:                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex unsigned short _Complex unsigned short (type.go:331:structLiteral: type.go:184:typ0: type.go:122:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041201-1.c`:                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex char _Complex char (type.go:17:typedef: type.go:184:typ0: type.go:122:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041214-1.c`:                 {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041218-2.c`:                 {}, // COMPILE FAIL: 20041218-2.c:7:10: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050121-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-1.c`:                 {}, // COMPILE FAIL: 20040709-1.c.go:389:5: undefined reference to '__builtin_classify_type' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-3.c`:                 {}, // COMPILE FAIL: 20040709-3.c.go:392:5: undefined reference to '__builtin_classify_type' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-2.c`:                 {}, // COMPILE FAIL: 20040709-2.c.go:388:5: undefined reference to '__builtin_classify_type' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050613-1.c`:                 {}, // COMPILE FAIL: 20050613-1.c:8:1: incomplete type: array of int (type.go:184:typ0: type.go:40:typ0: type.go:311:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20051012-1.c`:                 {}, // COMPILE FAIL: 341.go:353:5: TODO *gc.Arguments (check.go:285:checkExpr: check.go:623:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070614-1.c`:                 {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070919-1.c`:                 {}, // COMPILE FAIL: 20070919-1.c:31:7: invalid type size: 0 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071029-1.c`:                 {}, // COMPILE FAIL: TODO (expr.go:1438:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070824-1.c`:                 {}, // COMPILE FAIL: 20070824-1.c.go:334:34: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071120-1.c`:                 {}, // COMPILE FAIL: 20071120-1.c:13:41: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071210-1.c`:                 {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080519-1.c`:                 {}, // COMPILE FAIL: TODO exprIndex (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1044:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080502-1.c`:                 {}, // COMPILE FAIL: 20080502-1.c.go:320:5: undefined reference to '__builtin_signbit' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20081117-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090113-2.c`:                 {}, // COMPILE FAIL: 20090113-2.c:1:1: invalid type size: -1 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090113-3.c`:                 {}, // COMPILE FAIL: 20090113-3.c:1:1: invalid type size: -1 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090219-1.c`:                 {}, // COMPILE FAIL: 20090219-1.c:12:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20120111-1.c`:                 {}, // COMPILE FAIL: 432.go:3531:79: undefined: arg (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:285:checkExpr: check.go:1215:check: check.go:205:symbolResolver:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20181120-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20180921-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920302-1.c`:                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920415-1.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920428-2.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (stmt.go:23:statement: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-3.c`:                   {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-4.c`:                   {}, // COMPILE FAIL: TODO JumpStatementGotoExpr (stmt.go:182:blockItem: stmt.go:44:statement: stmt.go:365:jumpStatement)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-5.c`:                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-7.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-1.c`:                   {}, // COMPILE FAIL: 476.go:352:5: TODO *gc.Arguments (check.go:285:checkExpr: check.go:623:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920612-2.c`:                   {}, // COMPILE FAIL: 920612-2.c:6:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920625-1.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920721-4.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920731-1.c`:                   {}, // COMPILE FAIL: gcc.c-torture/execute/920731-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920908-1.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921215-1.c`:                   {}, // COMPILE FAIL: 921215-1.c:5:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921202-1.c`:                   {}, // COMPILE FAIL: 527.go:352:3: TODO *gc.Arguments (check.go:1972:check: check.go:623:check: check.go:655:checkFn:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921208-2.c`:                   {}, // COMPILE FAIL: 532.go:347:3: TODO *gc.Arguments (check.go:1972:check: check.go:623:check: check.go:655:checkFn:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930406-1.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (expr.go:2188:primaryExpression: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930718-1.c`:                   {}, // COMPILE FAIL: 564.go:360:26: undefined: tsrtx_def (check.go:1348:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:238:checkType: check.go:1348:check: check.go:205:symbolResolver:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931002-1.c`:                   {}, // COMPILE FAIL: 931002-1.c:10:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-10.c`:                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-14.c`:                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-12.c`:                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-2.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-6.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-4.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-8.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931031-1.c`:                   {}, // COMPILE FAIL: 593.go:363:86: value 8589934590 overflows uint32 (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/941202-1.c`:                   {}, // COMPILE FAIL: 941202-1.c.go:328:36: undefined reference to 'alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950628-1.c`:                   {}, // COMPILE FAIL: TODO (expr.go:67:expr: expr.go:448:expr0: expr.go:1299:postfixExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950906-1.c`:                   {}, // COMPILE FAIL: gcc.c-torture/execute/950906-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960416-1.c`:                   {}, // COMPILE FAIL: gcc.c-torture/execute/960416-1.c:57:23: internal error: t_be (expr.go:2170:primaryExpression: expr.go:2427:primaryExpressionIntConst: type.go:27:helper)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/980526-1.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/980929-1.c`:                   {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ann)))" pointer to int exprUintptr -> pointer to int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990208-1.c`:                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990222-1.c`:                   {}, // COMPILE FAIL: TODO (expr.go:2182:primaryExpression: expr.go:426:expr0: expr.go:1973:assignmentExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990326-1.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991014-1.c`:                   {}, // COMPILE FAIL: 991014-1.c:10:1: invalid type size: -496 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991118-1.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991112-1.c`:                   {}, // COMPILE FAIL: 991112-1.c.go:328:5: undefined reference to '__builtin_isprint' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-2.c`:                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:125:externalDeclaration: decl.go:274:declaration: decl.go:303:initDeclarator)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-3.c`:                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:125:externalDeclaration: decl.go:274:declaration: decl.go:303:initDeclarator)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-4.c`:                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:125:externalDeclaration: decl.go:274:declaration: decl.go:303:initDeclarator)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/align-1.c`:                    {}, // COMPILE FAIL: align-1.c:1:13: unsupported alignment 16 of new_int (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/anon-1.c`:                     {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:448:expr0: expr.go:1311:postfixExpression: expr.go:1586:postfixExpressionSelect)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alloca-1.c`:                   {}, // COMPILE FAIL: alloca-1.c.go:329:34: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf-sign-1.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf64-1.c`:                     {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-4.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-3.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-7.c`:                   {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-6.c`:                   {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/built-in-setjmp.c`:            {}, // COMPILE FAIL: built-in-setjmp.c.go:322:2: undefined reference to '__builtin_longjmp' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/builtin-types-compatible-p.c`: {}, // COMPILE FAIL: TODO *cc.InvalidType (expr.go:448:expr0: expr.go:1308:postfixExpression: expr.go:1722:postfixExpressionCall)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/builtin-bitops-1.c`:           {}, // COMPILE FAIL: builtin-bitops-1.c.go:1560:51: undefined reference to '__builtin_ctz' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/comp-goto-2.c`:                {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/call-trap-1.c`:                {}, // COMPILE FAIL: 797.go:357:2: TODO gc.PredefinedType, uintptr (check.go:1719:checkStatement: check.go:1972:check: check.go:597:check:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/comp-goto-1.c`:                {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-1.c`:                  {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-2.c`:                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-5.c`:                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-7.c`:                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/compndlit-1.c`:                {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-6.c`:                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ffs-1.c`:                      {}, // COMPILE FAIL: ffs-1.c.go:319:5: undefined reference to '__builtin_ffs' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ffs-2.c`:                      {}, // COMPILE FAIL: ffs-2.c.go:331:6: undefined reference to '__builtin_ffs' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/frame-address.c`:              {}, // COMPILE FAIL: frame-address.c.go:335:34: undefined reference to '__builtin_frame_address' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/fprintf-2.c`:                  {}, // COMPILE FAIL: fprintf-2.c.go:4466:15: undefined reference to 'tmpnam' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/compare-fp-1.c`:          {}, // COMPILE FAIL: compare-fp-1.c.go:332:26: undefined reference to '__builtin_islessgreater' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4.c`:              {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4f.c`:             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4l.c`:             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-5.c`:              {}, // COMPILE FAIL: fp-cmp-5.c.go:338:9: undefined reference to '__builtin_islessequal' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8f.c`:             {}, // COMPILE FAIL: fp-cmp-8f.c.go:364:5: undefined reference to '__builtin_islessequal' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8.c`:              {}, // COMPILE FAIL: fp-cmp-8.c.go:386:5: undefined reference to '__builtin_isgreater' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8l.c`:             {}, // COMPILE FAIL: fp-cmp-8l.c.go:429:5: undefined reference to '__builtin_islessgreater' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/mzero4.c`:                {}, // COMPILE FAIL: mzero4.c.go:359:18: undefined reference to 'atanf' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr38016.c`:               {}, // COMPILE FAIL: pr38016.c.go:342:5: undefined reference to '__builtin_isless' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/medce-1.c`:                    {}, // COMPILE FAIL: gcc.c-torture/execute/medce-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nest-align-1.c`:               {}, // COMPILE FAIL: nest-align-1.c:25:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nest-stdar-1.c`:               {}, // COMPILE FAIL: nest-stdar-1.c:5:10: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/memchr-1.c`:                   {}, // COMPILE FAIL: 943.go:377:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-3.c`:                 {}, // COMPILE FAIL: nestfunc-3.c:12:8: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-2.c`:                 {}, // COMPILE FAIL: nestfunc-2.c:13:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-1.c`:                 {}, // COMPILE FAIL: nestfunc-1.c:15:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-5.c`:                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-6.c`:                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-7.c`:                 {}, // COMPILE FAIL: nestfunc-7.c:15:12: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr19449.c`:                    {}, // COMPILE FAIL: gcc.c-torture/execute/pr19449.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr17377.c`:                    {}, // COMPILE FAIL: pr17377.c.go:331:36: undefined reference to '__builtin_return_address' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-1.c`:                  {}, // COMPILE FAIL: TODO (expr.go:448:expr0: expr.go:1248:postfixExpression: expr.go:1130:postfixExpressionIndex)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-1.c`:                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-1.c:10:24:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-4.c`:                  {}, // COMPILE FAIL: TODO UnaryExpressionSizeofExpr (expr.go:67:expr: expr.go:456:expr0: expr.go:1067:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-3.c`:                  {}, // COMPILE FAIL: pr22061-3.c:4:7: nested functions not supported (decl.go:189:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:185:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-3.c`:                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: Xf(aatls), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-3.c:12:24:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-2.c`:                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-2.c:10:24:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr23324.c`:                    {}, // COMPILE FAIL: pr23324.c:4:21: invalid type size: -1 (type.go:325:unionLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr24135.c`:                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr28289.c`:                    {}, // COMPILE FAIL: 1025.go:368:45: TODO *gc.Arguments (check.go:1972:check: check.go:623:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr30185.c`:                    {}, // COMPILE FAIL: TODO (expr.go:67:expr: expr.go:448:expr0: expr.go:1299:postfixExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr34154.c`:                    {}, // COMPILE FAIL: TODO LabeledStatementRange (stmt.go:182:blockItem: stmt.go:21:statement: stmt.go:71:labeledStatement)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr35456.c`:                    {}, // COMPILE FAIL: pr35456.c.go:336:7: undefined reference to '__builtin_signbit' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr36321.c`:                    {}, // COMPILE FAIL: pr36321.c.go:323:34: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr37780.c`:                    {}, // COMPILE FAIL: pr37780.c.go:322:35: undefined reference to '__builtin_ctz' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38151.c`:                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:331:structLiteral: type.go:184:typ0: type.go:122:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38969.c`:                    {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:67:expr: expr.go:456:expr0: expr.go:1115:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39339.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:211:initializerStruct: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39228.c`:                    {}, // COMPILE FAIL: pr39228.c.go:330:9: undefined reference to '__builtin_isinf' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39100.c`:                    {}, // COMPILE FAIL: 1103.go:356:6: undefined: tsC (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:226:checkExprOrType: check.go:1215:check: check.go:205:symbolResolver:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr42570.c`:                    {}, // COMPILE FAIL: pr42570.c:2:9: invalid type size: -1 (type.go:35:typ: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr42614.c`:                    {}, // COMPILE FAIL: TODO "siinlined_wrong" staticInternal (link.go:1046:print0: link.go:993:print0: link.go:949:name)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr44575.c`:                    {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr47237.c`:                    {}, // COMPILE FAIL: pr47237.c.go:334:2: undefined reference to '__builtin_apply' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49218.c`:                    {}, // COMPILE FAIL: pr49218.c:2:18: unsupported alignment 16 of L (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49644.c`:                    {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:450:expr0: expr.go:2172:primaryExpression: expr.go:2440:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49768.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51447.c`:                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:151:functionDefinition0: stmt.go:151:compoundStatement: stmt.go:180:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52979-2.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52979-1.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr54471.c`:                    {}, // COMPILE FAIL: pr54471.c:15:22: unsupported alignment 16 of unsigned __int128 (type.go:29:helper: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56837.c`:                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:35:typ: type.go:258:typ0: type.go:122:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57130.c`:                    {}, // COMPILE FAIL: 1237.go:336:28: undefined: tsS (check.go:1348:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:238:checkType: check.go:1348:check: check.go:205:symbolResolver:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58385.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((((iqlibc.ppBoolInt32(((((iqlibc.ppInt32FromInt32(0)) != 0)) || ((Xa != 0)))) & iqlibc.ppBoolInt32((Xfoo(aatls) >= (iqlibc.ppInt32FromInt32(0))))) <= (iqlibc.ppInt32FromInt32(1)))) && (((iqlibc.ppInt32FromInt32(1)) != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58431.c`:                    {}, // COMPILE FAIL: 1261.go:351:6: TODO int32 -> int16 (check.go:1729:checkStatement: check.go:1909:check: check.go:735:isAssignable:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58726.c`:                    {}, // COMPILE FAIL: TODO bitfield (decl.go:376:initDeclarator: type.go:35:typ: type.go:210:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58984.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60017.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60003.c`:                    {}, // COMPILE FAIL: pr60003.c.go:322:2: undefined reference to '__builtin_longjmp' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr61375.c`:                    {}, // COMPILE FAIL: pr61375.c:17:29: unsupported alignment 16 of __int128 (type.go:29:helper: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr63302.c`:                    {}, // COMPILE FAIL: TODO "(-((iqlibc.ppInt32FromInt32(1))))" int exprDefault -> __int128 exprDefault (pr63302.c:16:33:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr61725.c`:                    {}, // COMPILE FAIL: pr61725.c.go:321:9: undefined reference to '__builtin_ffs' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65170.c`:                    {}, // COMPILE FAIL: pr65170.c:4:27: unsupported alignment 16 of V (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64756.c`:                    {}, // COMPILE FAIL: TODO "pp_ = (((Xd != 0)) || ((Xd != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65215-3.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64242.c`:                    {}, // COMPILE FAIL: pr64242.c.go:325:2: undefined reference to '__builtin_longjmp' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64006.c`:                    {}, // COMPILE FAIL: pr64006.c.go:331:6: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65369.c`:                    {}, // COMPILE FAIL: pr65369.c:42:31: unsupported alignment 16 of array of 97 unsigned char (type.go:35:typ: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr66556.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr67037.c`:                    {}, // COMPILE FAIL: 1334.go:382:2: TODO *gc.Arguments (check.go:1972:check: check.go:626:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68185.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv3) && ((ccv2 != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68321.c`:                    {}, // COMPILE FAIL: TODO "pp_ = (((Xu != 0)) || ((Xn != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68249.c`:                    {}, // COMPILE FAIL: TODO (expr.go:432:expr0: expr.go:637:conditionalExpression: expr.go:51:expr)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68381.c`:                    {}, // COMPILE FAIL: pr68381.c.go:321:5: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70127.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70460.c`:                    {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70602.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71494.c`:                    {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:67:expr: expr.go:456:expr0: expr.go:1085:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71700.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78170.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71554.c`:                    {}, // COMPILE FAIL: pr71554.c.go:336:5: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79286.c`:                    {}, // COMPILE FAIL: pr79286.c:2:21: incomplete type: array of array of 8 int (type.go:35:typ: type.go:40:typ0: type.go:311:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79737-1.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr80692.c`:                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Decimal64 _Decimal64 (expr.go:2436:primaryExpressionFloatConst: type.go:29:helper: type.go:122:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82210.c`:                    {}, // COMPILE FAIL: pr82210.c:18:9: unsupported alignment 16 of struct S {a array of struct T; b array of int} (type.go:35:typ: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82387.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv2) && ((ccv1 != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84169.c`:                    {}, // COMPILE FAIL: pr84169.c:4:27: unsupported alignment 16 of T (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84748.c`:                    {}, // COMPILE FAIL: pr84748.c:3:27: unsupported alignment 16 of u128 (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84521.c`:                    {}, // COMPILE FAIL: pr84521.c.go:320:2: undefined reference to '__builtin_longjmp' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85529-1.c`:                  {}, // COMPILE FAIL: TODO "pp_ = ((Xs.fda) != iqlibc.ppBoolInt32(((ccv3) && ((ccv1 != 0)))))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85582-2.c`:                  {}, // COMPILE FAIL: pr85582-2.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85095.c`:                    {}, // COMPILE FAIL: pr85095.c.go:322:33: undefined reference to '__builtin_add_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85582-3.c`:                  {}, // COMPILE FAIL: pr85582-3.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr87053.c`:                    {}, // COMPILE FAIL: TODO (init.go:267:initializerUnion: type.go:35:typ: type.go:217:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr86528.c`:                    {}, // COMPILE FAIL: pr86528.c.go:321:36: undefined reference to '__builtin_alloca' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr89195.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr88714.c`:                    {}, // COMPILE FAIL: 1459.go:350:44: undefined: tsT (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:226:checkExprOrType: check.go:1215:check: check.go:205:symbolResolver:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr89434.c`:                    {}, // COMPILE FAIL: pr89434.c.go:324:2: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr90311.c`:                    {}, // COMPILE FAIL: pr90311.c.go:327:2: undefined reference to '__builtin_add_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91450-2.c`:                  {}, // COMPILE FAIL: pr91450-2.c.go:321:5: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91450-1.c`:                  {}, // COMPILE FAIL: pr91450-1.c.go:321:7: undefined reference to '__builtin_mul_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93213.c`:                    {}, // COMPILE FAIL: pr93213.c:8:27: unsupported alignment 16 of u128 (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93434.c`:                    {}, // COMPILE FAIL: TODO "(!((ani != 0)))" int exprBool -> double exprDefault (expr.go:753:equalityExpression: expr.go:75:expr: expr.go:145:convert)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91635.c`:                    {}, // COMPILE FAIL: pr91635.c.go:327:66: undefined reference to '__builtin_add_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93249.c`:                    {}, // COMPILE FAIL: pr93249.c.go:327:2: undefined reference to '__builtin_strncpy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93494.c`:                    {}, // COMPILE FAIL: pr93494.c.go:324:8: undefined reference to '__builtin_add_overflow' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94724.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((iqlibc.ppInt32FromInt16(Xb)) != (iqlibc.ppInt32FromInt32(53601)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94809.c`:                    {}, // COMPILE FAIL: TODO "pp_ = (((-((iqlibc.ppUint64FromUint64(1)))) / anone) < (iqlibc.ppUint64FromInt32(ccv1)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr96549.c`:                    {}, // COMPILE FAIL: gcc.c-torture/execute/pr96549.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr92904.c`:                    {}, // COMPILE FAIL: pr92904.c:12:10: unsupported alignment 16 of __int128 (type.go:35:typ: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr97325.c`:                    {}, // COMPILE FAIL: pr97325.c.go:320:10: undefined reference to '__builtin_ffs' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr98366.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr98474.c`:                    {}, // COMPILE FAIL: pr98474.c:4:21: unsupported alignment 16 of T (type.go:17:typedef: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/printf-2.c`:                   {}, // COMPILE FAIL: printf-2.c.go:4482:8: undefined reference to 'freopen' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/return-addr.c`:                {}, // COMPILE FAIL: 1534.go:565:8: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stdarg-3.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stdarg-1.c`:                   {}, // COMPILE FAIL: stdarg-1.c.go:450:2: undefined reference to '__builtin_va_copy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strct-stdarg-1.c`:             {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strct-varg-1.c`:               {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/string-opt-18.c`:              {}, // COMPILE FAIL: string-opt-18.c.go:332:5: undefined reference to 'mempcpy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-3.c`:                   {}, // COMPILE FAIL: 1567.go:377:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-5.c`:                   {}, // COMPILE FAIL: TODO (init.go:18:initializerOuter: init.go:72:initializer: init.go:236:initializerUnion)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-6.c`:                   {}, // COMPILE FAIL: 1570.go:338:21: TODO *gc.BasicLit (check.go:1013:check: check.go:1082:checkArray: check.go:1130:checkArrayElem:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-2.c`:                   {}, // COMPILE FAIL: 1566.go:1150:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-4.c`:                   {}, // COMPILE FAIL: 1568.go:377:9: TODO untyped int -> struct{} (check.go:1708:check: check.go:1729:checkStatement: check.go:1910:check:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/struct-ini-3.c`:               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/struct-ini-2.c`:               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-10.c`:                  {}, // COMPILE FAIL: va-arg-10.c.go:384:2: undefined reference to '__builtin_va_copy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/user-printf.c`:                {}, // COMPILE FAIL: user-printf.c.go:4519:8: undefined reference to 'freopen' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-14.c`:                  {}, // COMPILE FAIL: va-arg-14.c.go:374:2: undefined reference to '__builtin_va_copy' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-2.c`:                   {}, // COMPILE FAIL: 1598.go:409:5: TODO *gc.Arguments (check.go:285:checkExpr: check.go:623:check: check.go:655:checkFn:)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-22.c`:                  {}, // COMPILE FAIL: va-arg-22.c:24:1: invalid type size: 0 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-pack-1.c`:              {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zero-struct-1.c`:              {}, // COMPILE FAIL: zero-struct-1.c:1:1: invalid type size: -1 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zero-struct-2.c`:              {}, // COMPILE FAIL: zero-struct-2.c:3:19: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zerolen-1.c`:                  {}, // COMPILE FAIL: 1631.go:351:25: TODO *gc.StructType struct{Fname_len [1]uint8} (check.go:1444:check: check.go:226:checkExprOrType: check.go:1393:check:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0011-switch1.c`:                           {}, // COMPILE FAIL: mir/c-tests/andrewchambers_c/0011-switch1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0015-calls13.c`:                           {}, // COMPILE FAIL: 0015-calls13.c:1:1: invalid type size: -1 (type.go:331:structLiteral: type.go:40:typ0: type.go:316:checkValidType)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0025-duff.c`:                              {}, // COMPILE FAIL: mir/c-tests/andrewchambers_c/0025-duff.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits12.c`:                           {}, // COMPILE FAIL: TODO (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:217:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits13.c`:                           {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:448:expr0: expr.go:1311:postfixExpression: expr.go:1586:postfixExpressionSelect)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits14.c`:                           {}, // COMPILE FAIL: TODO (type.go:331:structLiteral: type.go:184:typ0: type.go:240:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits15.c`:                           {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:448:expr0: expr.go:1311:postfixExpression: expr.go:1586:postfixExpressionSelect)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/anonymous-members.c`:                                  {}, // COMPILE FAIL: TODO (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:217:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/anonymous-struct.c`:                                   {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:448:expr0: expr.go:1311:postfixExpression: expr.go:1586:postfixExpressionSelect)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-basic.c`:                                     {}, // COMPILE FAIL: bitfield-basic.c:6:3: unsupported alignment 4 of union A {int; b char} (type.go:325:unionLiteral: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-initialize-zero.c`:                           {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-immediate-bitwise.c`:                         {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-load.c`:                                      {}, // COMPILE FAIL: TODO bitfield (init.go:289:initializerUnion: type.go:35:typ: type.go:210:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-trailing-zero.c`:                             {}, // COMPILE FAIL: bitfield-trailing-zero.c:15:1: unsupported alignment 4 of union U1 {int; int; c char; int; int; int; int} (type.go:325:unionLiteral: type.go:40:typ0: type.go:295:checkValidType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-reset-align.c`:                               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-pack-next.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-packing.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-types.c`:                                     {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield.c`:                                           {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-types-init.c`:                                {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/conditional-void.c`:                                   {}, // COMPILE FAIL: 123.go:346:8: TODO () -> struct{} (check.go:1729:checkStatement: check.go:1909:check: check.go:735:isAssignable:)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/constant-expression.c`:                                {}, // COMPILE FAIL: 126.go:340:74: TODO U+002B '+' (check.go:1939:check: check.go:285:checkExpr: check.go:904:check:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/deref-compare-float.c`:                                {}, // COMPILE FAIL: TODO "((iqlibc.ppFloat64FromFloat64(0)) > (iqlibc.ppFloat64FromInt32((*(*ppint32)(iqunsafe.ppPointer(Xi))))))" int exprBool -> float exprDefault (expr.go:1862:assignmentExpression: expr.go:75:expr: expr.go:145:convert)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/declarator-abstract.c`:                                {}, // COMPILE FAIL: 134.go:431:72: undefined: a (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:285:checkExpr: check.go:1215:check: check.go:205:symbolResolver:)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/duffs-device.c`:                                       {}, // COMPILE FAIL: mir/c-tests/lacc/duffs-device.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/enum.c`:                                               {}, // COMPILE FAIL: 149.go:362:13: TODO U+002D '-' (check.go:1939:check: check.go:285:checkExpr: check.go:904:check:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/float-compare-equal.c`:                                {}, // COMPILE FAIL: TODO "((iqlibc.ppFloat32FromInt32(Xi)) == (*(*ppfloat32)(iqunsafe.ppPointer(Xf))))" int exprBool -> float exprDefault (expr.go:1862:assignmentExpression: expr.go:75:expr: expr.go:145:convert)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/function-pointer-call.c`:                              {}, // COMPILE FAIL: 168.go:334:30: TODO *gc.Arguments (check.go:1719:checkStatement: check.go:1972:check: check.go:604:check:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/identifier.c`:                                         {}, // COMPILE FAIL: invalid object file: multiple defintions of c (ccgo.go:252:Main: link.go:199:link: link.go:408:getFileSymbols)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/function-incomplete.c`:                                {}, // COMPILE FAIL: 167.go:332:44: TODO *gc.Arguments (check.go:285:checkExpr: check.go:623:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/initialize-call.c`:                                    {}, // COMPILE FAIL: TODO "sigetobj(aatls, (iqlibc.ppInt32FromInt32(2)), (iqlibc.ppInt32FromInt32(8)), (iqlibc.ppInt32FromInt32(1)))" struct obj {s array of 3 short} exprDefault -> struct {t struct obj} exprDefault (initialize-call.c:25:9:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/initialize-string.c`:                                  {}, // COMPILE FAIL: TODO "\"wat\\x00\"" pointer to char exprDefault -> array of 7 char exprDefault (initialize-string.c:14:14:) (expr.go:75:expr: expr.go:129:convert: expr.go:320:convertType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/return-bitfield.c`:                                    {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/return-point.c`:                                       {}, // COMPILE FAIL: 246.go:337:42: undefined: tspoint (check.go:1348:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:238:checkType: check.go:1348:check: check.go:205:symbolResolver:)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/short-circuit-comma.c`:                                {}, // COMPILE FAIL: TODO "pp_ = ((ccv1) && (((iqlibc.ppInt32FromInt32(1)) != 0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/string-addr.c`:                                        {}, // COMPILE FAIL: TODO (expr.go:450:expr0: expr.go:2178:primaryExpression: expr.go:2295:primaryExpressionStringConst)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/stringify.c`:                                          {}, // COMPILE FAIL: TODO `'\xA'` -> invalid syntax (expr.go:450:expr0: expr.go:2174:primaryExpression: expr.go:2368:primaryExpressionCharConst)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/tail-compare-jump.c`:                                  {}, // COMPILE FAIL: TODO "pp_ = ((iqlibc.ppInt64FromInt32(Xf)) < (iqlibc.ppInt64FromInt64(0)))" int exprBool -> int exprVoid (expr.go:75:expr: expr.go:125:convert: expr.go:271:convertMode)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/union-bitfield.c`:                                     {}, // COMPILE FAIL: TODO bitfield (type.go:346:defineUnion: type.go:325:unionLiteral: type.go:210:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/union-zero-init.c`:                                    {}, // COMPILE FAIL: loading union-zero-init.c.go: undefined type tueither (type.go:502:typeID0: type.go:487:typeID: type.go:537:typeID0) (ccgo.go:252:Main: link.go:220:link: link.go:660:link)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg-complex-1.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg-complex-2.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg.c`:                                             {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/new/issue142.c`:                                            {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2182:primaryExpression: expr.go:456:expr0: expr.go:1014:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/new/issue117.c`:                                            {}, // COMPILE FAIL: issue117.c:2:3: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/76_dollars_in_identifiers.c`:                                        {}, // COMPILE FAIL: 76_dollars_in_identifiers.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/tcc-0.9.27/tests/tests2/80_flexarray.c`:                                                     {}, // COMPILE FAIL: 80_flexarray.c:5:1: incomplete type: array of int (type.go:184:typ0: type.go:40:typ0: type.go:311:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/87_dead_code.c`:                                                     {}, // COMPILE FAIL: 87_dead_code.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/tcc-0.9.27/tests/tests2/73_arm64.c`:                                                         {}, // COMPILE FAIL: TODO (expr.go:67:expr: expr.go:448:expr0: expr.go:1299:postfixExpression)
	`assets/tcc-0.9.27/tests/tests2/85_asm-outside-function.c`:                                          {}, // COMPILE FAIL: 85_asm-outside-function.c.go:317:2: undefined reference to 'vide' (ccgo.go:252:Main: link.go:220:link: link.go:578:link)
	`assets/tcc-0.9.27/tests/tests2/81_types.c`:                                                         {}, // COMPILE FAIL: 69.go:361:9: TODO *gc.Arguments (check.go:285:checkExpr: check.go:626:check: check.go:655:checkFn:) (link.go:220:link: link.go:765:link: link.go:1069:postProcess)
	`assets/tcc-0.9.27/tests/tests2/89_nocode_wanted.c`:                                                 {}, // COMPILE FAIL: 89_nocode_wanted.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/tcc-0.9.27/tests/tests2/88_codeopt.c`:                                                       {}, // COMPILE FAIL: 88_codeopt.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:335:compile: compile.go:335:compile)
	`assets/tcc-0.9.27/tests/tests2/93_integer_promotion.c`:                                             {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/tcc-0.9.27/tests/tests2/90_struct-init.c`:                                                   {}, // COMPILE FAIL: 90_struct-init.c:2:19: invalid type size: -1 (type.go:17:typedef: type.go:40:typ0: type.go:316:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/94_generic.c`:                                                       {}, // COMPILE FAIL: TODO PrimaryExpressionGeneric (expr.go:67:expr: expr.go:450:expr0: expr.go:2197:primaryExpression)
	`assets/tcc-0.9.27/tests/tests2/95_bitfields_ms.c`:                                                  {}, // COMPILE FAIL: 95_bitfields.c:27:5: unsupported alignment 16 of struct __s {x int; y char; z long long; a char; b long long} (type.go:35:typ: type.go:40:typ0: type.go:295:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/95_bitfields.c`:                                                     {}, // COMPILE FAIL: 95_bitfields.c:27:5: unsupported alignment 16 of struct __s {x int; y char; z long long; a char; b long long} (type.go:35:typ: type.go:40:typ0: type.go:295:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/98_al_ax_extend.c`:                                                  {}, // COMPILE FAIL: 87.go:4029:8: undefined: _us (check.go:1215:check: check.go:204:symbolResolver: link.go:1264:SymbolResolver) (check.go:285:checkExpr: check.go:1215:check: check.go:205:symbolResolver:)
}