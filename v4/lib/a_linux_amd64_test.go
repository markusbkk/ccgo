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
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-6.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920728-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950221-1.c`:                 {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-6.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920728-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950221-1.c`: {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0010-goto1.c`:           {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/goto.c`:                             {}, // BUILD FAIL
	`assets/tcc-0.9.27/tests/tests2/54_goto.c`:                                        {}, // BUILD FAIL

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

	//TODO
	`assets/CompCert-3.6/test/c/aes.c`:                                                   {}, // BUILD FAIL
	`assets/benchmarksgame-team.pages.debian.net/reverse-complement-6.c`:                 {}, // BUILD FAIL
	`assets/ccgo/bug/sqlite.c`:                                                           {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20051012-1.c`:                  {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20120111-1.c`:                  {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-1.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921202-1.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921208-2.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930718-1.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/call-trap-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/memchr-1.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39100.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57130.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58431.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64006.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr67037.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68381.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71554.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85095.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr88714.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr89434.c`:                     {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-3.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-6.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-2.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-4.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zerolen-1.c`:                   {}, // BUILD FAIL
	`assets/github.com/AbsInt/CompCert/test/c/aes.c`:                                     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20051012-1.c`:  {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20120111-1.c`:  {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-1.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921202-1.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921208-2.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930718-1.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/call-trap-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/memchr-1.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39100.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57130.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58431.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64006.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr67037.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68381.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71554.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85095.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr88714.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr90311.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr89434.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91450-1.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91450-2.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93494.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91635.c`:     {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/return-addr.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-3.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-2.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-6.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-4.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zerolen-1.c`:   {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/conditional-void.c`:                    {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/declarator-abstract.c`:                 {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/function-pointer-call.c`:               {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/function-incomplete.c`:                 {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/return-point.c`:                        {}, // BUILD FAIL
	`assets/tcc-0.9.27/tests/tests2/81_types.c`:                                          {}, // BUILD FAIL
	`assets/tcc-0.9.27/tests/tests2/98_al_ax_extend.c`:                                   {}, // BUILD FAIL

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
	`assets/CompCert-3.6/test/c/fannkuch.c`:                                                             {}, // COMPILE FAIL: TODO (expr.go:2354:primaryExpression: expr.go:441:expr0: expr.go:2139:assignmentExpression)
	`assets/benchmarksgame-team.pages.debian.net/fasta-4.c`:                                             {}, // COMPILE FAIL: fasta-4.c.go:4653:4: undefined reference to 'errx' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-8.c`:                                        {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/benchmarksgame-team.pages.debian.net/reverse-complement-4.c`:                                {}, // COMPILE FAIL: reverse-complement-4.c.go:4069:38: undefined reference to '__builtin_memmove' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/ccgo/bug/csmith2.c`:                                                                         {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:239:typ0)
	`assets/ccgo/bug/bitfield.c`:                                                                        {}, // COMPILE FAIL: TODO bitfield (expr.go:1458:postfixExpression: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/ccgo/bug/union2.c`:                                                                          {}, // COMPILE FAIL: TODO (expr.go:1458:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/ccgo/bug/union.c`:                                                                           {}, // COMPILE FAIL: TODO (expr.go:1458:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/ccgo/bug/struct.c`:                                                                          {}, // COMPILE FAIL: TODO union {ceiling char; unused char} 1, ft struct outer {magic char; pad1 array of 3 char; union {ceiling char; unused char}; pad2 array of 3 char; owner char} 9 (init.go:211:initializerStruct: init.go:72:initializer: init.go:269:initializerUnion)
	`assets/ccgo/bug/union3.c`:                                                                          {}, // COMPILE FAIL: TODO (type.go:17:typedef: type.go:190:typ0: type.go:247:typ0)
	`assets/ccgo/bug/union4.c`:                                                                          {}, // COMPILE FAIL: TODO (type.go:17:typedef: type.go:190:typ0: type.go:247:typ0)
	`assets/ccgo/bug/struct2.c`:                                                                         {}, // COMPILE FAIL: TODO (expr.go:1458:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000224-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:72:initializer: init.go:267:initializerUnion: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000326-2.c`:                                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000405-3.c`:                                 {}, // COMPILE FAIL: 20000405-3.c:1:1: unsupported alignment 32 of struct foo {entry array of 40 pointer to void} (type.go:338:structLiteral: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000211-1.c`:                                 {}, // COMPILE FAIL: 20000211-1.c:32:28: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000827-1.c`:                                 {}, // COMPILE FAIL: 20000827-1.c:13:7: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000923-1.c`:                                 {}, // COMPILE FAIL: 20000923-1.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000518-1.c`:                                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000804-1.c`:                                 {}, // COMPILE FAIL: 20000804-1.c:17:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20010202-1.c`:                                 {}, // COMPILE FAIL: 20010202-1.c:1:22: variable length arrays are not supported (decl.go:258:signature: type.go:294:checkValidParamType: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20001121-1.c`:                                 {}, // COMPILE FAIL: 20001121-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20001205-1.c`:                                 {}, // COMPILE FAIL: 20001205-1.c:6:9: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20001222-1.c`:                                 {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20001221-1.c`:                                 {}, // COMPILE FAIL: 20001221-1.c:11:10: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20010118-1.c`:                                 {}, // COMPILE FAIL: 20010118-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20010226-1.c`:                                 {}, // COMPILE FAIL: 20010226-1.c:16:12: nested functions not supported (stmt.go:24:statement: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20010605-1.c`:                                 {}, // COMPILE FAIL: 20010605-1.c:9:9: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20010903-1.c`:                                 {}, // COMPILE FAIL: 20010903-1.c:7:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20010903-2.c`:                                 {}, // COMPILE FAIL: 20010903-2.c:9:14: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20011023-1.c`:                                 {}, // COMPILE FAIL: 20011023-1.c:8:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20011205-1.c`:                                 {}, // COMPILE FAIL: 20011205-1.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20011217-2.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:1458:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20020210-1.c`:                                 {}, // COMPILE FAIL: 20020210-1.c:2:34: invalid type size: 0 (decl.go:258:signature: type.go:294:checkValidParamType: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20020309-1.c`:                                 {}, // COMPILE FAIL: 20020309-1.c:8:5: nested functions not supported (stmt.go:24:statement: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20020312-1.c`:                                 {}, // COMPILE FAIL: 20020312-1.c:12:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20020318-1.c`:                                 {}, // COMPILE FAIL: TODO "ppuintptr(iqunsafe.ppPointer(&(anc)))" pointer to char pointer, exprUintptr -> void void, exprDefault (expr.go:79:expr: expr.go:137:convert: expr.go:368:convertFromPointer)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20020807-1.c`:                                 {}, // COMPILE FAIL: 20020807-1.c:20:8: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20021110.c`:                                   {}, // COMPILE FAIL: 20021110.c:4:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20021204-1.c`:                                 {}, // COMPILE FAIL: 20021204-1.c:8:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20021108-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030224-1.c`:                                 {}, // COMPILE FAIL: 20030224-1.c:6:25: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030503-1.c`:                                 {}, // COMPILE FAIL: TODO (stmt.go:363:iterationStatement: stmt.go:530:iterationStatementFlat: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030418-1.c`:                                 {}, // COMPILE FAIL: 20030418-1.c:13:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030410-1.c`:                                 {}, // COMPILE FAIL: 20030410-1.c:11:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030530-3.c`:                                 {}, // COMPILE FAIL: 20030530-3.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030716-1.c`:                                 {}, // COMPILE FAIL: 20030716-1.c:3:21: variable length arrays are not supported (decl.go:258:signature: type.go:294:checkValidParamType: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030910-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20030903-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20031011-1.c`:                                 {}, // COMPILE FAIL: 20031011-1.c:15:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20040101-1.c`:                                 {}, // COMPILE FAIL: 20040101-1.c:22:9: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20040310-1.c`:                                 {}, // COMPILE FAIL: 20040310-1.c:4:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20040317-1.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:463:expr0: expr.go:1270:postfixExpression: expr.go:1152:postfixExpressionIndex)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20040323-1.c`:                                 {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20040611-1.c`:                                 {}, // COMPILE FAIL: 20040611-1.c:7:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20040317-3.c`:                                 {}, // COMPILE FAIL: 20040317-3.c:4:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20040614-1.c`:                                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20050113-1.c`:                                 {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20050119-1.c`:                                 {}, // COMPILE FAIL: 20050119-1.c:7:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20050122-2.c`:                                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20050510-1.c`:                                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20050801-2.c`:                                 {}, // COMPILE FAIL: 20050801-2.c:6:5: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20051216-1.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:2354:primaryExpression: expr.go:441:expr0: expr.go:2139:assignmentExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20060217-1.c`:                                 {}, // COMPILE FAIL: 20060217-1.c:21:7: assembler statements not supported (stmt.go:351:unbracedStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20070603-2.c`:                                 {}, // COMPILE FAIL: TODO exprIndex (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1066:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20070603-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20071108-1.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:463:expr0: expr.go:1270:postfixExpression: expr.go:1152:postfixExpressionIndex)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20080114-1.c`:                                 {}, // COMPILE FAIL: 20080114-1.c:9:19: assembler statements not supported (stmt.go:235:selectionStatement: stmt.go:342:bracedStatement: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20070919-1.c`:                                 {}, // COMPILE FAIL: TODO exprUintptr (expr.go:71:expr: expr.go:463:expr0: expr.go:1396:postfixExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20080613-1.c`:                                 {}, // COMPILE FAIL: 20080613-1.c:6:1: incomplete type: array of unsigned char (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20071107-1.c`:                                 {}, // COMPILE FAIL: 20071107-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20090907-1.c`:                                 {}, // COMPILE FAIL: 20090907-1.c:19:1: invalid type size: -1 (type.go:190:typ0: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20110131-1.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20110902.c`:                                   {}, // COMPILE FAIL: 20110902.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20121107-1.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20161124-1.c`:                                 {}, // COMPILE FAIL: TODO "pp_ = (((anf != 0)) || ((Xe != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920301-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920428-3.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920428-4.c`:                                   {}, // COMPILE FAIL: TODO exprVoid (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1775:postfixExpressionSelect)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920415-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920501-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920501-22.c`:                                  {}, // COMPILE FAIL: 920501-22.c:1:13: incomplete type: array of int (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920502-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920501-7.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920624-1.c`:                                   {}, // COMPILE FAIL: 920624-1.c:1:5: incomplete type: array of int (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920826-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920831-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/920928-5.c`:                                   {}, // COMPILE FAIL: 920928-5.c:2:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/921012-1.c`:                                   {}, // COMPILE FAIL: 921012-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/930118-1.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/921203-1.c`:                                   {}, // COMPILE FAIL: 921203-1.c:1:6: incomplete type: array of char (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/930506-2.c`:                                   {}, // COMPILE FAIL: 930506-2.c:5:9: nested functions not supported (stmt.go:24:statement: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/931003-1.c`:                                   {}, // COMPILE FAIL: 931003-1.c:4:1: incomplete type: array of double (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/941019-1.c`:                                   {}, // COMPILE FAIL: 941019-1.c:1:54: unsupported alignment 16 of _Complex long double (decl.go:258:signature: type.go:294:checkValidParamType: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/941014-4.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/931018-1.c`:                                   {}, // COMPILE FAIL: TODO exprVoid (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1741:postfixExpressionSelect)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/950919-1.c`:                                   {}, // COMPILE FAIL: 950919-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:291:compile: compile.go:291:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/951116-1.c`:                                   {}, // COMPILE FAIL: 951116-1.c:7:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/950613-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/950610-1.c`:                                   {}, // COMPILE FAIL: 950610-1.c:1:15: variable length arrays are not supported (decl.go:258:signature: type.go:294:checkValidParamType: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/951128-1.c`:                                   {}, // COMPILE FAIL: 951128-1.c:1:6: incomplete type: array of char (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/980511-1.c`:                                   {}, // COMPILE FAIL: 980511-1.c:3:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/980506-1.c`:                                   {}, // COMPILE FAIL: TODO (stmt.go:370:iterationStatement: stmt.go:530:iterationStatementFlat: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/981001-2.c`:                                   {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/981001-4.c`:                                   {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/990517-1.c`:                                   {}, // COMPILE FAIL: 990517-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/981006-1.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/981223-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/991213-3.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/991213-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionImag (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1135:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/asmgoto-1.c`:                                  {}, // COMPILE FAIL: asmgoto-1.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/bbb.c`:                                        {}, // COMPILE FAIL: bbb.c:6:17: incomplete type: array of struct looksets (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/complex-1.c`:                                  {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (expr.go:326:convertType: type.go:29:helper: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/complex-3.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionImag (expr.go:71:expr: expr.go:471:expr0: expr.go:1135:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/complex-4.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/complex-2.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionImag (expr.go:71:expr: expr.go:471:expr0: expr.go:1135:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/complex-6.c`:                                  {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (expr.go:2599:primaryExpressionIntConst: type.go:29:helper: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/complex-5.c`:                                  {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (expr.go:2599:primaryExpressionIntConst: type.go:29:helper: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/debugvlafunction-1.c`:                         {}, // COMPILE FAIL: debugvlafunction-1.c:6:17: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/ex.c`:                                         {}, // COMPILE FAIL: ex.c:12:19: too few arguments to function 'foo', type 'function(int, int) returning int' in 'foo ()' (expr.go:463:expr0: expr.go:1340:postfixExpression: expr.go:1903:postfixExpressionCall)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/icfmatch.c`:                                   {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/init-3.c`:                                     {}, // COMPILE FAIL: init-3.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/labels-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/labels-3.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/labels-2.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/mangle-1.c`:                                   {}, // COMPILE FAIL: TODO *cc.UnknownValue (expr.go:465:expr0: expr.go:2342:primaryExpression: expr.go:2579:primaryExpressionIntConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/nested-3.c`:                                   {}, // COMPILE FAIL: nested-3.c:13:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/nested-1.c`:                                   {}, // COMPILE FAIL: nested-1.c:15:18: variable length arrays are not supported (type.go:190:typ0: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/nested-2.c`:                                   {}, // COMPILE FAIL: nested-2.c:9:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pass.c`:                                       {}, // COMPILE FAIL: pass.c:13:10: too many arguments to function 'foo', type 'function(int, int, int) returning int' in 'foo ((int) & q, q, w, e, q, (int) &w)' (expr.go:463:expr0: expr.go:1340:postfixExpression: expr.go:1908:postfixExpressionCall)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pc44485.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr16566-1.c`:                                  {}, // COMPILE FAIL: pr16566-1.c:6:1: incomplete type: array of pointer to struct S (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr16566-2.c`:                                  {}, // COMPILE FAIL: pr16566-2.c:5:1: incomplete type: array of int (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr16566-3.c`:                                  {}, // COMPILE FAIL: pr16566-3.c:5:1: incomplete type: array of pointer to struct S (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr17407.c`:                                    {}, // COMPILE FAIL: TODO exprIndex (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1066:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr17397.c`:                                    {}, // COMPILE FAIL: pr17397.c:10:8: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr17529.c`:                                    {}, // COMPILE FAIL: pr17529.c:5:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr17913.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr17906.c`:                                    {}, // COMPILE FAIL: pr17906.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr18903.c`:                                    {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr21356.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr21728.c`:                                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr21839.c`:                                    {}, // COMPILE FAIL: pr21839.c:1:21: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr22013-1.c`:                                  {}, // COMPILE FAIL: TODO "[1]tnW{0: (iqlibc.ppUint16FromUint8('R')), }" array of 1 W exprDefault -> P exprDefault (pr22013-1.c:9:15:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr22422.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = (((and)+4)+ppuintptr(((*tsD)(iqunsafe.ppPointer(and)).fdn))*4)" pointer to int exprUintptr -> pointer to int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr23237.c`:                                    {}, // COMPILE FAIL: pr23237.c:8:78: incomplete type: array of pointer to function(void) returning int (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr25224.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr26425.c`:                                    {}, // COMPILE FAIL: pr26425.c:1:1: invalid type size: 0 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27282.c`:                                    {}, // COMPILE FAIL: pr27282.c:4:1: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27341-1.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionImag (expr.go:71:expr: expr.go:471:expr0: expr.go:1135:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27889.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (decl.go:297:signature: type.go:41:typ2: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27373.c`:                                    {}, // COMPILE FAIL: TODO exprIndex (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1066:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27528.c`:                                    {}, // COMPILE FAIL: pr27528.c:34:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27907.c`:                                    {}, // COMPILE FAIL: pr27907.c:2:20: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27863.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr28489.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr30311.c`:                                    {}, // COMPILE FAIL: pr30311.c:6:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr30984.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr29128.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr29201.c`:                                    {}, // COMPILE FAIL: pr29201.c:35:35: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr32355.c`:                                    {}, // COMPILE FAIL: pr32355.c:4:1: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr32584.c`:                                    {}, // COMPILE FAIL: pr32584.c:8:3: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr32919.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr33173.c`:                                    {}, // COMPILE FAIL: pr33173.c:5:1: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr33382.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr33614.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr32571.c`:                                    {}, // COMPILE FAIL: pr32571.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr33617.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr33855.c`:                                    {}, // COMPILE FAIL: pr33855.c:20:7: unsupported alignment 16 of _Complex long double (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr34334.c`:                                    {}, // COMPILE FAIL: pr34334.c:8:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr34091.c`:                                    {}, // COMPILE FAIL: pr34091.c:13:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr34856.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr34688.c`:                                    {}, // COMPILE FAIL: pr34688.c:4:10: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr34966.c`:                                    {}, // COMPILE FAIL: pr34966.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr35431.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (expr.go:2599:primaryExpressionIntConst: type.go:29:helper: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr35318.c`:                                    {}, // COMPILE FAIL: pr35318.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr35432.c`:                                    {}, // COMPILE FAIL: pr35432.c:3:1: invalid type size: 0 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr34885.c`:                                    {}, // COMPILE FAIL: TODO "(iqlibc.ppUintptrFromInt32(0))" pointer to void exprDefault -> __CONST_SOCKADDR_ARG exprDefault (pr34885.c:8:17:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr35006.c`:                                    {}, // COMPILE FAIL: pr35006.c:8:8: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr37056.c`:                                    {}, // COMPILE FAIL: pr37056.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr37395.c`:                                    {}, // COMPILE FAIL: pr37395.c:7:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr37380.c`:                                    {}, // COMPILE FAIL: pr37380.c:17:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr38123.c`:                                    {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr38360.c`:                                    {}, // COMPILE FAIL: pr38360.c:7:3: too few arguments to function 'fputs', type 'function(pointer to char, pointer to void) returning int' in 'fputs ("")' (expr.go:463:expr0: expr.go:1340:postfixExpression: expr.go:1903:postfixExpressionCall)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr38771.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr38789.c`:                                    {}, // COMPILE FAIL: pr38789.c:10:5: assembler statements not supported (stmt.go:238:selectionStatement: stmt.go:354:unbracedStatement: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr39394.c`:                                    {}, // COMPILE FAIL: pr39394.c:9:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr39845.c`:                                    {}, // COMPILE FAIL: TODO exprUintptr (expr.go:463:expr0: expr.go:1345:postfixExpression: expr.go:1655:postfixExpressionPSelect)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr39928-1.c`:                                  {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr40233.c`:                                    {}, // COMPILE FAIL: pr40233.c:1:13: unsupported alignment 64 of aligned (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr40556.c`:                                    {}, // COMPILE FAIL: pr40556.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr40351.c`:                                    {}, // COMPILE FAIL: TODO (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:247:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr40321.c`:                                    {}, // COMPILE FAIL: pr40321.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr41469.c`:                                    {}, // COMPILE FAIL: pr41469.c:13:8: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42196-1.c`:                                  {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:332:unionLiteral: type.go:231:typ0: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42196-2.c`:                                  {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:332:unionLiteral: type.go:253:typ0: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42196-3.c`:                                  {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:332:unionLiteral: type.go:253:typ0: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr41987.c`:                                    {}, // COMPILE FAIL: pr41987.c:15:3: unsupported alignment 16 of _Complex long double (type.go:29:helper: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42164.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:71:expr: expr.go:463:expr0: expr.go:1331:postfixExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42398.c`:                                    {}, // COMPILE FAIL: pr42398.c:4:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42559.c`:                                    {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42716.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42717.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr42956.c`:                                    {}, // COMPILE FAIL: pr42956.c:4:3: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr43066.c`:                                    {}, // COMPILE FAIL: pr43066.c:1:1: invalid type size: -1 (type.go:190:typ0: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr43164.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:211:initializerStruct: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr43191.c`:                                    {}, // COMPILE FAIL: pr43191.c:3:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr43661.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr43679.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((*(*ppint32)(iqunsafe.ppPointer(Xg_29))) != (*(*ppint32)(iqunsafe.ppPointer(anl_39))))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr44038.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:465:expr0: expr.go:2350:primaryExpression: expr.go:2467:primaryExpressionStringConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr44197.c`:                                    {}, // COMPILE FAIL: pr44197.c:27:55: incomplete type: array of __ctype_mask_t (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr44043.c`:                                    {}, // COMPILE FAIL: pr44043.c:14:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr44119.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr44707.c`:                                    {}, // COMPILE FAIL: pr44707.c:14:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr44941.c`:                                    {}, // COMPILE FAIL: pr44941.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr46107.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr46547-1.c`:                                  {}, // COMPILE FAIL: pr46547-1.c:5:7: unsupported alignment 16 of _Complex long double (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr46637.c`:                                    {}, // COMPILE FAIL: pr46637.c:9:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr47428.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr48596.c`:                                    {}, // COMPILE FAIL: pr48596.c:5:20: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr50009.c`:                                    {}, // COMPILE FAIL: pr50009.c:3:1: incomplete type: array of short (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr48641.c`:                                    {}, // COMPILE FAIL: TODO "(X__builtin_modfl(aatls, (iqlibc.ppFloat64FromFloat64(1.5)), aabp) != (*(*ppfloat64)(iqunsafe.ppPointer(aabp))))" int exprBool -> long double exprDefault (expr.go:778:equalityExpression: expr.go:79:expr: expr.go:149:convert)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr51495.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr51354.c`:                                    {}, // COMPILE FAIL: pr51354.c:6:44: unsupported alignment 32 of ai (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr51694.c`:                                    {}, // COMPILE FAIL: pr51694.c:14:3: too few arguments to function 'foo', type 'function(int, pointer to function())' in 'foo (x)' (expr.go:463:expr0: expr.go:1340:postfixExpression: expr.go:1903:postfixExpressionCall)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr51767.c`:                                    {}, // COMPILE FAIL: pr51767.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr52750.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr53409.c`:                                    {}, // COMPILE FAIL: pr53409.c:4:5: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr53410-2.c`:                                  {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr53748.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54103-3.c`:                                  {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54103-1.c`:                                  {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54103-4.c`:                                  {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54103-2.c`:                                  {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54103-6.c`:                                  {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54103-5.c`:                                  {}, // COMPILE FAIL: TODO "pp_ = (!((((iqlibc.ppInt32FromInt32(0)) / (iqlibc.ppInt32FromInt32(0))) != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54559.c`:                                    {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54713-1.c`:                                  {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54713-2.c`:                                  {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr55851.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:463:expr0: expr.go:1270:postfixExpression: expr.go:1152:postfixExpressionIndex)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr55921.c`:                                    {}, // COMPILE FAIL: pr55921.c:18:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr56405.c`:                                    {}, // COMPILE FAIL: pr56405.c:6:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr54713-3.c`:                                  {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr56571.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((((Xa % Xb) != 0)) && ((Xf(aatls) != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr58164.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr58344.c`:                                    {}, // COMPILE FAIL: pr58344.c:4:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr60502.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr63282.c`:                                    {}, // COMPILE FAIL: pr63282.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr65241.c`:                                    {}, // COMPILE FAIL: pr65241.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr69989-2.c`:                                  {}, // COMPILE FAIL: TODO (stmt.go:363:iterationStatement: stmt.go:530:iterationStatementFlat: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70061.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70355.c`:                                    {}, // COMPILE FAIL: pr70355.c:4:27: unsupported alignment 16 of v2ti (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70199.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70190.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70263-1.c`:                                  {}, // COMPILE FAIL: TODO exprVoid (expr.go:463:expr0: expr.go:1270:postfixExpression: expr.go:1222:postfixExpressionIndex)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70240.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70633.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr70916.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr77754-1.c`:                                  {}, // COMPILE FAIL: pr77754-1.c:9:7: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr77754-2.c`:                                  {}, // COMPILE FAIL: pr77754-2.c:5:18: variable length arrays are not supported (decl.go:258:signature: type.go:294:checkValidParamType: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr77754-3.c`:                                  {}, // COMPILE FAIL: pr77754-3.c:5:25: variable length arrays are not supported (decl.go:258:signature: type.go:294:checkValidParamType: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr77754-4.c`:                                  {}, // COMPILE FAIL: pr77754-4.c:5:25: variable length arrays are not supported (decl.go:258:signature: type.go:294:checkValidParamType: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr77754-5.c`:                                  {}, // COMPILE FAIL: pr77754-5.c:6:14: variable length arrays are not supported (decl.go:258:signature: type.go:294:checkValidParamType: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr72802.c`:                                    {}, // COMPILE FAIL: pr72802.c:1:8: incomplete type: array of int (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr77754-6.c`:                                  {}, // COMPILE FAIL: pr77754-6.c:8:7: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr78439.c`:                                    {}, // COMPILE FAIL: pr78439.c:52:29: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr79089.c`:                                    {}, // COMPILE FAIL: pr79089.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr79413.c`:                                    {}, // COMPILE FAIL: pr79413.c:7:7: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr79621.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv2) && ((Xb5 != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr81262.c`:                                    {}, // COMPILE FAIL: pr81262.c:8:3: assembler statements not supported (stmt.go:176:compoundStatement: stmt.go:207:blockItem: stmt.go:51:statement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr82337.c`:                                    {}, // COMPILE FAIL: pr82337.c:7:1: incomplete type: array of char (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr82564.c`:                                    {}, // COMPILE FAIL: pr82564.c:9:5: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr82875.c`:                                    {}, // COMPILE FAIL: pr82875.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr83575.c`:                                    {}, // COMPILE FAIL: TODO (stmt.go:363:iterationStatement: stmt.go:530:iterationStatementFlat: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr83487.c`:                                    {}, // COMPILE FAIL: pr83487.c:3:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr84136.c`:                                    {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr82052.c`:                                    {}, // COMPILE FAIL: pr82052.c:19:10: incomplete type: array of array of 7 array of 2 uint16_t (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr84305.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:463:expr0: expr.go:1270:postfixExpression: expr.go:1152:postfixExpressionIndex)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr85704.c`:                                    {}, // COMPILE FAIL: pr85704.c:3:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr85945.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr84960.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr86122.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (expr.go:2599:primaryExpressionIntConst: type.go:29:helper: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr86123.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex unsigned _Complex unsigned (expr.go:2599:primaryExpressionIntConst: type.go:29:helper: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr87110.c`:                                    {}, // COMPILE FAIL: pr87110.c:8:12: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr87647.c`:                                    {}, // COMPILE FAIL: pr87647.c:3:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr87468.c`:                                    {}, // COMPILE FAIL: pr87468.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr89655.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr90139.c`:                                    {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/simd-1.c`:                                     {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/simd-4.c`:                                     {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/simd-2.c`:                                     {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/simd-6.c`:                                     {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/simd-5.c`:                                     {}, // COMPILE FAIL: TODO vector (expr.go:2599:primaryExpressionIntConst: type.go:29:helper: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/simd-3.c`:                                     {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/sra-1.c`:                                      {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/statement-expression-1.c`:                     {}, // COMPILE FAIL: statement-expression-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/struct-non-lval-3.c`:                          {}, // COMPILE FAIL: TODO exprLvalue (expr.go:2354:primaryExpression: expr.go:441:expr0: expr.go:2030:assignmentExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/uuarg.c`:                                      {}, // COMPILE FAIL: uuarg.c:4:10: too few arguments to function 'foo', type 'function(int, int, int, int, int, int, int, int, int) returning int' in 'foo ()' (expr.go:463:expr0: expr.go:1340:postfixExpression: expr.go:1903:postfixExpressionCall)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vector-2.c`:                                   {}, // COMPILE FAIL: TODO vector (type.go:338:structLiteral: type.go:190:typ0: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vector-1.c`:                                   {}, // COMPILE FAIL: TODO vector (type.go:338:structLiteral: type.go:190:typ0: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vector-6.c`:                                   {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vector-5.c`:                                   {}, // COMPILE FAIL: TODO vector (decl.go:402:initDeclarator: type.go:17:typedef: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vector-4.c`:                                   {}, // COMPILE FAIL: TODO vector (decl.go:432:initDeclarator: type.go:35:typ: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vector-align-1.c`:                             {}, // COMPILE FAIL: vector-align-1.c:11:6: unsupported alignment 128 of char (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vla-const-1.c`:                                {}, // COMPILE FAIL: vla-const-1.c:6:20: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vla-const-2.c`:                                {}, // COMPILE FAIL: vla-const-2.c:5:20: variable length arrays are not supported (type.go:35:typ: type.go:46:typ0: type.go:312:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/vector-3.c`:                                   {}, // COMPILE FAIL: TODO vector (expr.go:2610:primaryExpressionFloatConst: type.go:29:helper: type.go:63:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/zero-strct-1.c`:                               {}, // COMPILE FAIL: zero-strct-1.c:1:20: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/zero-strct-2.c`:                               {}, // COMPILE FAIL: zero-strct-2.c:1:18: invalid type size: -1 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/zero-strct-4.c`:                               {}, // COMPILE FAIL: zero-strct-4.c:1:19: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/zero-strct-3.c`:                               {}, // COMPILE FAIL: zero-strct-3.c:1:19: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/zero-strct-5.c`:                               {}, // COMPILE FAIL: zero-strct-5.c:4:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000113-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000703-1.c`:                                 {}, // COMPILE FAIL: 20000703-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000801-3.c`:                                 {}, // COMPILE FAIL: 20000801-3.c:5:20: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000822-1.c`:                                 {}, // COMPILE FAIL: 20000822-1.c:12:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000917-1.c`:                                 {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ana)))" pointer to int3 exprUintptr -> pointer to int3 exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20001203-2.c`:                                 {}, // COMPILE FAIL: 20001203-2.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010122-1.c`:                                 {}, // COMPILE FAIL: 20010122-1.c.go:321:37: undefined reference to '__builtin_return_address' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010605-1.c`:                                 {}, // COMPILE FAIL: 20010605-1.c:5:14: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010605-2.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010904-2.c`:                                 {}, // COMPILE FAIL: 20010904-2.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010904-1.c`:                                 {}, // COMPILE FAIL: 20010904-1.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20011113-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020206-2.c`:                                 {}, // COMPILE FAIL: TODO "tnA{fda: (iqlibc.ppUint16FromInt32(((ccv1 << (iqlibc.ppInt32FromInt32(8))) + (ani << (iqlibc.ppInt32FromInt32(4)))))), }" A exprDefault -> A exprSelect (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020320-1.c`:                                 {}, // COMPILE FAIL: TODO "aav3" struct large {x int; y array of 9 int} exprDefault -> struct large {x int; y array of 9 int} exprSelect (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020314-1.c`:                                 {}, // COMPILE FAIL: 20020314-1.c.go:328:34: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020411-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020412-1.c`:                                 {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20021113-1.c`:                                 {}, // COMPILE FAIL: 20021113-1.c.go:329:36: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030109-1.c`:                                 {}, // COMPILE FAIL: 20030109-1.c:5:1: incomplete type: array of int (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030501-1.c`:                                 {}, // COMPILE FAIL: 20030501-1.c:7:9: nested functions not supported (stmt.go:24:statement: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030323-1.c`:                                 {}, // COMPILE FAIL: 20030323-1.c.go:320:36: undefined reference to '__builtin_return_address' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030330-1.c`:                                 {}, // COMPILE FAIL: 20030330-1.c.go:322:4: undefined reference to 'link_error' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030714-1.c`:                                 {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1761:postfixExpressionSelect)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030811-1.c`:                                 {}, // COMPILE FAIL: 20030811-1.c:9:5: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030910-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040302-1.c`:                                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040307-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040223-1.c`:                                 {}, // COMPILE FAIL: 20040223-1.c.go:2440:53: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040520-1.c`:                                 {}, // COMPILE FAIL: 20040520-1.c:6:13: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040423-1.c`:                                 {}, // COMPILE FAIL: 20040423-1.c:13:14: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040411-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionSizeofType (expr.go:71:expr: expr.go:471:expr0: expr.go:1102:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041124-1.c`:                                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex unsigned short _Complex unsigned short (type.go:338:structLiteral: type.go:190:typ0: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041201-1.c`:                                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex char _Complex char (type.go:17:typedef: type.go:190:typ0: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041214-1.c`:                                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041218-2.c`:                                 {}, // COMPILE FAIL: 20041218-2.c:7:10: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050121-1.c`:                                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-1.c`:                                 {}, // COMPILE FAIL: 20040709-1.c.go:387:5: undefined reference to '__builtin_classify_type' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-3.c`:                                 {}, // COMPILE FAIL: 20040709-3.c.go:390:5: undefined reference to '__builtin_classify_type' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-2.c`:                                 {}, // COMPILE FAIL: 20040709-2.c.go:386:5: undefined reference to '__builtin_classify_type' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050613-1.c`:                                 {}, // COMPILE FAIL: 20050613-1.c:8:1: incomplete type: array of int (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070614-1.c`:                                 {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070919-1.c`:                                 {}, // COMPILE FAIL: 20070919-1.c:31:7: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071029-1.c`:                                 {}, // COMPILE FAIL: TODO (expr.go:1458:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070824-1.c`:                                 {}, // COMPILE FAIL: 20070824-1.c.go:332:34: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071120-1.c`:                                 {}, // COMPILE FAIL: 20071120-1.c:13:41: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071210-1.c`:                                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080519-1.c`:                                 {}, // COMPILE FAIL: TODO exprIndex (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1066:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080502-1.c`:                                 {}, // COMPILE FAIL: 20080502-1.c.go:318:5: undefined reference to '__builtin_signbit' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20081117-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090219-1.c`:                                 {}, // COMPILE FAIL: 20090219-1.c:12:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090113-3.c`:                                 {}, // COMPILE FAIL: 20090113-3.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090113-2.c`:                                 {}, // COMPILE FAIL: 20090113-2.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20181120-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20180921-1.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920302-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920415-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920428-2.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (stmt.go:24:statement: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-3.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-4.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-7.c`:                                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-5.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920612-2.c`:                                   {}, // COMPILE FAIL: 920612-2.c:6:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920625-1.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920721-4.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920908-1.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921215-1.c`:                                   {}, // COMPILE FAIL: 921215-1.c:5:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930406-1.c`:                                   {}, // COMPILE FAIL: TODO *cc.LabelDeclaration (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931002-1.c`:                                   {}, // COMPILE FAIL: 931002-1.c:10:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-10.c`:                                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-12.c`:                                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-14.c`:                                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-4.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-2.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-6.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-8.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/941202-1.c`:                                   {}, // COMPILE FAIL: 941202-1.c.go:326:36: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950628-1.c`:                                   {}, // COMPILE FAIL: TODO (expr.go:71:expr: expr.go:463:expr0: expr.go:1331:postfixExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950906-1.c`:                                   {}, // COMPILE FAIL: 950906-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960416-1.c`:                                   {}, // COMPILE FAIL: 960416-1.c:57:23: internal error: t_be (expr.go:2342:primaryExpression: expr.go:2599:primaryExpressionIntConst: type.go:27:helper)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/980526-1.c`:                                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/980929-1.c`:                                   {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ann)))" pointer to int exprUintptr -> pointer to int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990208-1.c`:                                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990222-1.c`:                                   {}, // COMPILE FAIL: TODO (expr.go:2354:primaryExpression: expr.go:441:expr0: expr.go:2139:assignmentExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990326-1.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991014-1.c`:                                   {}, // COMPILE FAIL: 991014-1.c:10:1: invalid type size: -496 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991118-1.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991112-1.c`:                                   {}, // COMPILE FAIL: 991112-1.c.go:326:5: undefined reference to '__builtin_isprint' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-2.c`:                                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-3.c`:                                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-4.c`:                                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/align-1.c`:                                    {}, // COMPILE FAIL: align-1.c:1:13: unsupported alignment 16 of new_int (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/anon-1.c`:                                     {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1752:postfixExpressionSelect)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf-sign-1.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alloca-1.c`:                                   {}, // COMPILE FAIL: alloca-1.c.go:327:34: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf64-1.c`:                                     {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-3.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-4.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-7.c`:                                   {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-6.c`:                                   {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/builtin-types-compatible-p.c`:                 {}, // COMPILE FAIL: TODO *cc.InvalidType (expr.go:463:expr0: expr.go:1340:postfixExpression: expr.go:1888:postfixExpressionCall)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/built-in-setjmp.c`:                            {}, // COMPILE FAIL: built-in-setjmp.c.go:332:5: undefined reference to '__builtin_setjmp' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/builtin-bitops-1.c`:                           {}, // COMPILE FAIL: builtin-bitops-1.c.go:1558:6: undefined reference to '__builtin_clrsb' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/comp-goto-2.c`:                                {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-1.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-2.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/compndlit-1.c`:                                {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/comp-goto-1.c`:                                {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-5.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-7.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-6.c`:                                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ffs-1.c`:                                      {}, // COMPILE FAIL: ffs-1.c.go:317:5: undefined reference to '__builtin_ffs' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ffs-2.c`:                                      {}, // COMPILE FAIL: ffs-2.c.go:329:6: undefined reference to '__builtin_ffs' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/frame-address.c`:                              {}, // COMPILE FAIL: frame-address.c.go:333:34: undefined reference to '__builtin_frame_address' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/fprintf-2.c`:                                  {}, // COMPILE FAIL: fprintf-2.c.go:4457:15: undefined reference to 'tmpnam' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/compare-fp-1.c`:                          {}, // COMPILE FAIL: compare-fp-1.c.go:330:26: undefined reference to '__builtin_islessgreater' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4.c`:                              {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4l.c`:                             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4f.c`:                             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-5.c`:                              {}, // COMPILE FAIL: fp-cmp-5.c.go:346:9: undefined reference to '__builtin_isgreater' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8l.c`:                             {}, // COMPILE FAIL: fp-cmp-8l.c.go:405:5: undefined reference to '__builtin_isgreaterequal' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8f.c`:                             {}, // COMPILE FAIL: fp-cmp-8f.c.go:340:5: undefined reference to '__builtin_isless' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8.c`:                              {}, // COMPILE FAIL: fp-cmp-8.c.go:428:5: undefined reference to '__builtin_islessgreater' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/mzero4.c`:                                {}, // COMPILE FAIL: mzero4.c.go:356:18: undefined reference to 'tanf' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr38016.c`:                               {}, // COMPILE FAIL: pr38016.c.go:362:5: undefined reference to '__builtin_islessequal' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/medce-1.c`:                                    {}, // COMPILE FAIL: medce-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nest-align-1.c`:                               {}, // COMPILE FAIL: nest-align-1.c:25:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nest-stdar-1.c`:                               {}, // COMPILE FAIL: nest-stdar-1.c:5:10: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-1.c`:                                 {}, // COMPILE FAIL: nestfunc-1.c:15:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-2.c`:                                 {}, // COMPILE FAIL: nestfunc-2.c:13:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-5.c`:                                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-3.c`:                                 {}, // COMPILE FAIL: nestfunc-3.c:12:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-7.c`:                                 {}, // COMPILE FAIL: nestfunc-7.c:15:12: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-6.c`:                                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr17377.c`:                                    {}, // COMPILE FAIL: pr17377.c.go:329:36: undefined reference to '__builtin_return_address' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr19449.c`:                                    {}, // COMPILE FAIL: pr19449.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-4.c`:                                  {}, // COMPILE FAIL: TODO UnaryExpressionSizeofExpr (expr.go:71:expr: expr.go:471:expr0: expr.go:1089:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-1.c`:                                  {}, // COMPILE FAIL: TODO (expr.go:463:expr0: expr.go:1270:postfixExpression: expr.go:1152:postfixExpressionIndex)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-3.c`:                                  {}, // COMPILE FAIL: pr22061-3.c:4:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-2.c`:                                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-2.c:10:24:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-1.c`:                                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-1.c:10:24:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-3.c`:                                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: Xf(aatls), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-3.c:12:24:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr23324.c`:                                    {}, // COMPILE FAIL: pr23324.c:4:21: invalid type size: -1 (type.go:332:unionLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr24135.c`:                                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr30185.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:71:expr: expr.go:463:expr0: expr.go:1331:postfixExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr34154.c`:                                    {}, // COMPILE FAIL: TODO LabeledStatementRange (stmt.go:207:blockItem: stmt.go:22:statement: stmt.go:78:labeledStatement)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr35456.c`:                                    {}, // COMPILE FAIL: pr35456.c.go:334:7: undefined reference to '__builtin_signbit' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr36321.c`:                                    {}, // COMPILE FAIL: pr36321.c.go:321:34: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38151.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:338:structLiteral: type.go:190:typ0: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr37780.c`:                                    {}, // COMPILE FAIL: pr37780.c.go:320:35: undefined reference to '__builtin_ctz' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38969.c`:                                    {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39339.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:211:initializerStruct: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39228.c`:                                    {}, // COMPILE FAIL: pr39228.c.go:321:9: undefined reference to '__builtin_isinff' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr42570.c`:                                    {}, // COMPILE FAIL: pr42570.c:2:9: invalid type size: -1 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr42614.c`:                                    {}, // COMPILE FAIL: TODO "siinlined_wrong" staticInternal (link.go:1140:print0: link.go:1087:print0: link.go:1043:name)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43783.c`:                                    {}, // COMPILE FAIL: pr43783.c:6:3: unsupported alignment 16 of UINT192 (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr44575.c`:                                    {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr47237.c`:                                    {}, // COMPILE FAIL: pr47237.c.go:328:2: undefined reference to '__builtin_apply' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49218.c`:                                    {}, // COMPILE FAIL: pr49218.c:2:18: unsupported alignment 16 of L (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49644.c`:                                    {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49768.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51447.c`:                                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52979-1.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52979-2.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr54471.c`:                                    {}, // COMPILE FAIL: pr54471.c:15:22: unsupported alignment 16 of unsigned __int128 (type.go:29:helper: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56837.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:35:typ: type.go:265:typ0: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58385.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((((iqlibc.ppBoolInt32(((((iqlibc.ppInt32FromInt32(0)) != 0)) || ((Xa != 0)))) & iqlibc.ppBoolInt32((Xfoo(aatls) >= (iqlibc.ppInt32FromInt32(0))))) <= (iqlibc.ppInt32FromInt32(1)))) && (((iqlibc.ppInt32FromInt32(1)) != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58726.c`:                                    {}, // COMPILE FAIL: TODO bitfield (decl.go:432:initDeclarator: type.go:35:typ: type.go:217:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58984.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60017.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60003.c`:                                    {}, // COMPILE FAIL: pr60003.c.go:333:5: undefined reference to '__builtin_setjmp' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr61375.c`:                                    {}, // COMPILE FAIL: pr61375.c:17:29: unsupported alignment 16 of __int128 (type.go:29:helper: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr61725.c`:                                    {}, // COMPILE FAIL: pr61725.c.go:319:9: undefined reference to '__builtin_ffs' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr63302.c`:                                    {}, // COMPILE FAIL: TODO "(-((iqlibc.ppInt32FromInt32(1))))" int exprDefault -> __int128 exprDefault (pr63302.c:16:33:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64756.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = (((Xd != 0)) || ((Xd != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65170.c`:                                    {}, // COMPILE FAIL: pr65170.c:4:27: unsupported alignment 16 of V (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65215-3.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65369.c`:                                    {}, // COMPILE FAIL: pr65369.c:42:31: unsupported alignment 16 of array of 97 unsigned char (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr66556.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68185.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv3) && ((ccv2 != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68249.c`:                                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68321.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = (((Xu != 0)) || ((Xn != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70127.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70460.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70602.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71494.c`:                                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71700.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78170.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79286.c`:                                    {}, // COMPILE FAIL: pr79286.c:2:21: incomplete type: array of array of 8 int (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79737-1.c`:                                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr80692.c`:                                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Decimal64 _Decimal64 (expr.go:2608:primaryExpressionFloatConst: type.go:29:helper: type.go:128:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82387.c`:                                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv2) && ((ccv1 != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82210.c`:                                    {}, // COMPILE FAIL: pr82210.c:18:9: unsupported alignment 16 of struct S {a array of struct T; b array of int} (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84748.c`:                                    {}, // COMPILE FAIL: pr84748.c:3:27: unsupported alignment 16 of u128 (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84169.c`:                                    {}, // COMPILE FAIL: pr84169.c:4:27: unsupported alignment 16 of T (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85529-1.c`:                                  {}, // COMPILE FAIL: TODO "pp_ = ((Xs.fda) != iqlibc.ppBoolInt32(((ccv3) && ((ccv1 != 0)))))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85582-2.c`:                                  {}, // COMPILE FAIL: pr85582-2.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85582-3.c`:                                  {}, // COMPILE FAIL: pr85582-3.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr87053.c`:                                    {}, // COMPILE FAIL: TODO (init.go:267:initializerUnion: type.go:35:typ: type.go:224:typ0)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr86528.c`:                                    {}, // COMPILE FAIL: pr86528.c.go:318:36: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr89195.c`:                                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/printf-2.c`:                                   {}, // COMPILE FAIL: printf-2.c.go:4473:8: undefined reference to 'freopen' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stdarg-3.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strct-varg-1.c`:                               {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stdarg-1.c`:                                   {}, // COMPILE FAIL: stdarg-1.c.go:448:2: undefined reference to '__builtin_va_copy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strct-stdarg-1.c`:                             {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/string-opt-18.c`:                              {}, // COMPILE FAIL: string-opt-18.c.go:330:5: undefined reference to 'mempcpy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/struct-ini-2.c`:                               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/struct-ini-3.c`:                               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-5.c`:                                   {}, // COMPILE FAIL: TODO (init.go:18:initializerOuter: init.go:72:initializer: init.go:236:initializerUnion)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-10.c`:                                  {}, // COMPILE FAIL: va-arg-10.c.go:384:2: undefined reference to '__builtin_va_copy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/user-printf.c`:                                {}, // COMPILE FAIL: user-printf.c.go:4509:15: undefined reference to 'tmpnam' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-14.c`:                                  {}, // COMPILE FAIL: va-arg-14.c.go:372:2: undefined reference to '__builtin_va_copy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-22.c`:                                  {}, // COMPILE FAIL: va-arg-22.c:24:1: invalid type size: 0 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-pack-1.c`:                              {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zero-struct-1.c`:                              {}, // COMPILE FAIL: zero-struct-1.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zero-struct-2.c`:                              {}, // COMPILE FAIL: zero-struct-2.c:3:19: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/AbsInt/CompCert/test/c/fannkuch.c`:                                               {}, // COMPILE FAIL: TODO (expr.go:2354:primaryExpression: expr.go:441:expr0: expr.go:2139:assignmentExpression)
	`assets/github.com/cxgo/empty array decl.c`:                                                         {}, // COMPILE FAIL: empty array decl.c:3:6: invalid type size: -1 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/cxgo/inet.c`:                                                                     {}, // COMPILE FAIL: socket.h:275:1: incomplete type: array of unsigned char (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/github.com/cxgo/rename decl struct.c`:                                                       {}, // COMPILE FAIL: rename decl struct.c:2:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/cxgo/struct and var.c`:                                                           {}, // COMPILE FAIL: struct and var.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000113-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000703-1.c`:                 {}, // COMPILE FAIL: gcc.c-torture/execute/20000703-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000801-3.c`:                 {}, // COMPILE FAIL: 20000801-3.c:5:20: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000822-1.c`:                 {}, // COMPILE FAIL: 20000822-1.c:12:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000917-1.c`:                 {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ana)))" pointer to int3 exprUintptr -> pointer to int3 exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20001203-2.c`:                 {}, // COMPILE FAIL: gcc.c-torture/execute/20001203-2.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010122-1.c`:                 {}, // COMPILE FAIL: 20010122-1.c.go:344:11: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010605-1.c`:                 {}, // COMPILE FAIL: 20010605-1.c:5:14: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010605-2.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010904-1.c`:                 {}, // COMPILE FAIL: 20010904-1.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010904-2.c`:                 {}, // COMPILE FAIL: 20010904-2.c:9:72: unsupported alignment 32 of X (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20011113-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020206-2.c`:                 {}, // COMPILE FAIL: TODO "tnA{fda: (iqlibc.ppUint16FromInt32(((ccv1 << (iqlibc.ppInt32FromInt32(8))) + (ani << (iqlibc.ppInt32FromInt32(4)))))), }" A exprDefault -> A exprSelect (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020320-1.c`:                 {}, // COMPILE FAIL: TODO "aav3" struct large {x int; y array of 9 int} exprDefault -> struct large {x int; y array of 9 int} exprSelect (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020314-1.c`:                 {}, // COMPILE FAIL: 20020314-1.c.go:328:34: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020411-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020412-1.c`:                 {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20021113-1.c`:                 {}, // COMPILE FAIL: 20021113-1.c.go:329:36: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030109-1.c`:                 {}, // COMPILE FAIL: 20030109-1.c:5:1: incomplete type: array of int (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030330-1.c`:                 {}, // COMPILE FAIL: 20030330-1.c.go:322:4: undefined reference to 'link_error' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030323-1.c`:                 {}, // COMPILE FAIL: 20030323-1.c.go:320:36: undefined reference to '__builtin_return_address' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030501-1.c`:                 {}, // COMPILE FAIL: 20030501-1.c:7:9: nested functions not supported (stmt.go:24:statement: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030811-1.c`:                 {}, // COMPILE FAIL: 20030811-1.c:9:5: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030714-1.c`:                 {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1761:postfixExpressionSelect)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030910-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040302-1.c`:                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040307-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040223-1.c`:                 {}, // COMPILE FAIL: 20040223-1.c.go:2440:53: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040411-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionSizeofType (expr.go:71:expr: expr.go:471:expr0: expr.go:1102:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040423-1.c`:                 {}, // COMPILE FAIL: 20040423-1.c:13:14: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040520-1.c`:                 {}, // COMPILE FAIL: 20040520-1.c:6:13: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041124-1.c`:                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex unsigned short _Complex unsigned short (type.go:338:structLiteral: type.go:190:typ0: type.go:128:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041201-1.c`:                 {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex char _Complex char (type.go:17:typedef: type.go:190:typ0: type.go:128:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041214-1.c`:                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041218-2.c`:                 {}, // COMPILE FAIL: 20041218-2.c:7:10: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050121-1.c`:                 {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-3.c`:                 {}, // COMPILE FAIL: 20040709-3.c.go:390:5: undefined reference to '__builtin_classify_type' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-2.c`:                 {}, // COMPILE FAIL: 20040709-2.c.go:386:5: undefined reference to '__builtin_classify_type' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-1.c`:                 {}, // COMPILE FAIL: 20040709-1.c.go:387:5: undefined reference to '__builtin_classify_type' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050613-1.c`:                 {}, // COMPILE FAIL: 20050613-1.c:8:1: incomplete type: array of int (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070614-1.c`:                 {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070919-1.c`:                 {}, // COMPILE FAIL: 20070919-1.c:31:7: invalid type size: 0 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071029-1.c`:                 {}, // COMPILE FAIL: TODO (expr.go:1458:postfixExpression: init.go:72:initializer: init.go:242:initializerUnion)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071120-1.c`:                 {}, // COMPILE FAIL: 20071120-1.c:13:41: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070824-1.c`:                 {}, // COMPILE FAIL: 20070824-1.c.go:332:34: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071210-1.c`:                 {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080519-1.c`:                 {}, // COMPILE FAIL: TODO exprIndex (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1066:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080502-1.c`:                 {}, // COMPILE FAIL: 20080502-1.c.go:318:5: undefined reference to '__builtin_signbit' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20081117-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090113-2.c`:                 {}, // COMPILE FAIL: 20090113-2.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090113-3.c`:                 {}, // COMPILE FAIL: 20090113-3.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090219-1.c`:                 {}, // COMPILE FAIL: 20090219-1.c:12:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20181120-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20180921-1.c`:                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920302-1.c`:                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920415-1.c`:                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920428-2.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (stmt.go:24:statement: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-3.c`:                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-7.c`:                   {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-5.c`:                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-4.c`:                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920612-2.c`:                   {}, // COMPILE FAIL: 920612-2.c:6:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920625-1.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920721-4.c`:                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920908-1.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921215-1.c`:                   {}, // COMPILE FAIL: 921215-1.c:5:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930406-1.c`:                   {}, // COMPILE FAIL: TODO *cc.LabelDeclaration (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931002-1.c`:                   {}, // COMPILE FAIL: 931002-1.c:10:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-10.c`:                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-12.c`:                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-14.c`:                  {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-2.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-4.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-6.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-8.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/941202-1.c`:                   {}, // COMPILE FAIL: 941202-1.c.go:326:36: undefined reference to 'alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950628-1.c`:                   {}, // COMPILE FAIL: TODO (expr.go:71:expr: expr.go:463:expr0: expr.go:1331:postfixExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950906-1.c`:                   {}, // COMPILE FAIL: gcc.c-torture/execute/950906-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960416-1.c`:                   {}, // COMPILE FAIL: gcc.c-torture/execute/960416-1.c:57:23: internal error: t_be (expr.go:2342:primaryExpression: expr.go:2599:primaryExpressionIntConst: type.go:27:helper)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/980526-1.c`:                   {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/980929-1.c`:                   {}, // COMPILE FAIL: TODO "pp_ = ppuintptr(iqunsafe.ppPointer(&(ann)))" pointer to int exprUintptr -> pointer to int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990208-1.c`:                   {}, // COMPILE FAIL: TODO UnaryExpressionLabelAddr (expr.go:71:expr: expr.go:471:expr0: expr.go:1107:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990222-1.c`:                   {}, // COMPILE FAIL: TODO (expr.go:2354:primaryExpression: expr.go:441:expr0: expr.go:2139:assignmentExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990326-1.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991014-1.c`:                   {}, // COMPILE FAIL: 991014-1.c:10:1: invalid type size: -496 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991118-1.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991112-1.c`:                   {}, // COMPILE FAIL: 991112-1.c.go:326:5: undefined reference to '__builtin_isprint' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-2.c`:                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-4.c`:                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-3.c`:                    {}, // COMPILE FAIL: TODO unsupported attribute(s) (decl.go:181:externalDeclaration: decl.go:330:declaration: decl.go:359:initDeclarator)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/align-1.c`:                    {}, // COMPILE FAIL: align-1.c:1:13: unsupported alignment 16 of new_int (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/anon-1.c`:                     {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1752:postfixExpressionSelect)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf-sign-1.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alloca-1.c`:                   {}, // COMPILE FAIL: alloca-1.c.go:327:34: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf64-1.c`:                     {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-3.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-4.c`:                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-6.c`:                   {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-7.c`:                   {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/built-in-setjmp.c`:            {}, // COMPILE FAIL: built-in-setjmp.c.go:328:34: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/builtin-types-compatible-p.c`: {}, // COMPILE FAIL: TODO *cc.InvalidType (expr.go:463:expr0: expr.go:1340:postfixExpression: expr.go:1888:postfixExpressionCall)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/builtin-bitops-1.c`:           {}, // COMPILE FAIL: builtin-bitops-1.c.go:1579:6: undefined reference to '__builtin_clrsbl' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/comp-goto-2.c`:                {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/comp-goto-1.c`:                {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-2.c`:                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-1.c`:                  {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-5.c`:                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/compndlit-1.c`:                {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-7.c`:                  {}, // COMPILE FAIL: TODO cc.Complex64Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-6.c`:                  {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ffs-1.c`:                      {}, // COMPILE FAIL: ffs-1.c.go:317:5: undefined reference to '__builtin_ffs' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ffs-2.c`:                      {}, // COMPILE FAIL: ffs-2.c.go:329:6: undefined reference to '__builtin_ffs' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/frame-address.c`:              {}, // COMPILE FAIL: frame-address.c.go:333:34: undefined reference to '__builtin_frame_address' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/fprintf-2.c`:                  {}, // COMPILE FAIL: fprintf-2.c.go:4460:15: undefined reference to 'tmpnam' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/compare-fp-1.c`:          {}, // COMPILE FAIL: compare-fp-1.c.go:330:26: undefined reference to '__builtin_islessgreater' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4.c`:              {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4l.c`:             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4f.c`:             {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-5.c`:              {}, // COMPILE FAIL: fp-cmp-5.c.go:356:9: undefined reference to '__builtin_isgreaterequal' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8.c`:              {}, // COMPILE FAIL: fp-cmp-8.c.go:340:5: undefined reference to '__builtin_isless' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8f.c`:             {}, // COMPILE FAIL: fp-cmp-8f.c.go:428:5: undefined reference to '__builtin_islessgreater' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8l.c`:             {}, // COMPILE FAIL: fp-cmp-8l.c.go:405:5: undefined reference to '__builtin_isgreaterequal' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/mzero4.c`:                {}, // COMPILE FAIL: mzero4.c.go:357:18: undefined reference to 'atanf' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr38016.c`:               {}, // COMPILE FAIL: pr38016.c.go:384:5: undefined reference to '__builtin_isgreater' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/medce-1.c`:                    {}, // COMPILE FAIL: gcc.c-torture/execute/medce-1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nest-align-1.c`:               {}, // COMPILE FAIL: nest-align-1.c:25:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-1.c`:                 {}, // COMPILE FAIL: nestfunc-1.c:15:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-2.c`:                 {}, // COMPILE FAIL: nestfunc-2.c:13:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nest-stdar-1.c`:               {}, // COMPILE FAIL: nest-stdar-1.c:5:10: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-6.c`:                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-7.c`:                 {}, // COMPILE FAIL: nestfunc-7.c:15:12: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-5.c`:                 {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-3.c`:                 {}, // COMPILE FAIL: nestfunc-3.c:12:8: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr19449.c`:                    {}, // COMPILE FAIL: gcc.c-torture/execute/pr19449.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr17377.c`:                    {}, // COMPILE FAIL: pr17377.c.go:329:36: undefined reference to '__builtin_return_address' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-1.c`:                  {}, // COMPILE FAIL: TODO (expr.go:463:expr0: expr.go:1270:postfixExpression: expr.go:1152:postfixExpressionIndex)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-3.c`:                  {}, // COMPILE FAIL: pr22061-3.c:4:7: nested functions not supported (decl.go:245:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:210:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-4.c`:                  {}, // COMPILE FAIL: TODO UnaryExpressionSizeofExpr (expr.go:71:expr: expr.go:471:expr0: expr.go:1089:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-2.c`:                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-2.c:10:24:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-1.c`:                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: (iqlibc.ppInt32FromInt32(1)), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-1.c:10:24:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-3.c`:                  {}, // COMPILE FAIL: TODO "[3]ppint32{0: (iqlibc.ppInt32FromInt32(0)), 1: Xf(aatls), 2: (iqlibc.ppInt32FromInt32(2)), }" array of 3 int exprDefault -> pointer to int exprDefault (pr22098-3.c:12:24:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr23324.c`:                    {}, // COMPILE FAIL: pr23324.c:4:21: invalid type size: -1 (type.go:332:unionLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr24135.c`:                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr30185.c`:                    {}, // COMPILE FAIL: TODO (expr.go:71:expr: expr.go:463:expr0: expr.go:1331:postfixExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr34154.c`:                    {}, // COMPILE FAIL: TODO LabeledStatementRange (stmt.go:207:blockItem: stmt.go:22:statement: stmt.go:78:labeledStatement)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr35456.c`:                    {}, // COMPILE FAIL: pr35456.c.go:334:7: undefined reference to '__builtin_signbit' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr36321.c`:                    {}, // COMPILE FAIL: pr36321.c.go:321:34: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38151.c`:                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:338:structLiteral: type.go:190:typ0: type.go:128:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr37780.c`:                    {}, // COMPILE FAIL: pr37780.c.go:320:35: undefined reference to '__builtin_ctz' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38969.c`:                    {}, // COMPILE FAIL: TODO UnaryExpressionReal (expr.go:71:expr: expr.go:471:expr0: expr.go:1137:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39339.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:211:initializerStruct: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39228.c`:                    {}, // COMPILE FAIL: pr39228.c.go:336:9: undefined reference to '__builtin_isinfl' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr42570.c`:                    {}, // COMPILE FAIL: pr42570.c:2:9: invalid type size: -1 (type.go:35:typ: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr42614.c`:                    {}, // COMPILE FAIL: TODO "siinlined_wrong" staticInternal (link.go:1140:print0: link.go:1087:print0: link.go:1043:name)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr44575.c`:                    {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr47237.c`:                    {}, // COMPILE FAIL: pr47237.c.go:328:46: undefined reference to '__builtin_apply_args' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49218.c`:                    {}, // COMPILE FAIL: pr49218.c:2:18: unsupported alignment 16 of L (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49644.c`:                    {}, // COMPILE FAIL: TODO cc.Complex128Value (expr.go:465:expr0: expr.go:2344:primaryExpression: expr.go:2612:primaryExpressionFloatConst)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51447.c`:                    {}, // COMPILE FAIL: TODO BlockItemLabel (decl.go:207:functionDefinition0: stmt.go:176:compoundStatement: stmt.go:205:blockItem)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49768.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52979-2.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52979-1.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr54471.c`:                    {}, // COMPILE FAIL: pr54471.c:15:22: unsupported alignment 16 of unsigned __int128 (type.go:29:helper: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56837.c`:                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Complex int _Complex int (type.go:35:typ: type.go:265:typ0: type.go:128:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58385.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((((iqlibc.ppBoolInt32(((((iqlibc.ppInt32FromInt32(0)) != 0)) || ((Xa != 0)))) & iqlibc.ppBoolInt32((Xfoo(aatls) >= (iqlibc.ppInt32FromInt32(0))))) <= (iqlibc.ppInt32FromInt32(1)))) && (((iqlibc.ppInt32FromInt32(1)) != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58726.c`:                    {}, // COMPILE FAIL: TODO bitfield (decl.go:432:initDeclarator: type.go:35:typ: type.go:217:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58984.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60017.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60003.c`:                    {}, // COMPILE FAIL: pr60003.c.go:320:2: undefined reference to '__builtin_longjmp' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr61375.c`:                    {}, // COMPILE FAIL: pr61375.c:17:29: unsupported alignment 16 of __int128 (type.go:29:helper: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr63302.c`:                    {}, // COMPILE FAIL: TODO "(-((iqlibc.ppInt32FromInt32(1))))" int exprDefault -> __int128 exprDefault (pr63302.c:16:33:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr61725.c`:                    {}, // COMPILE FAIL: pr61725.c.go:319:9: undefined reference to '__builtin_ffs' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64242.c`:                    {}, // COMPILE FAIL: pr64242.c.go:323:2: undefined reference to '__builtin_longjmp' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64756.c`:                    {}, // COMPILE FAIL: TODO "pp_ = (((Xd != 0)) || ((Xd != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65170.c`:                    {}, // COMPILE FAIL: pr65170.c:4:27: unsupported alignment 16 of V (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65215-3.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65369.c`:                    {}, // COMPILE FAIL: pr65369.c:42:31: unsupported alignment 16 of array of 97 unsigned char (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr66556.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68249.c`:                    {}, // COMPILE FAIL: TODO (expr.go:447:expr0: expr.go:662:conditionalExpression: expr.go:55:expr)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68185.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv3) && ((ccv2 != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68321.c`:                    {}, // COMPILE FAIL: TODO "pp_ = (((Xu != 0)) || ((Xn != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70127.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70460.c`:                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70602.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71494.c`:                    {}, // COMPILE FAIL: TODO <nil> (decl.go:192:functionDefinition: decl.go:204:functionDefinition0: decl.go:73:newFnCtx)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71700.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78170.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79286.c`:                    {}, // COMPILE FAIL: pr79286.c:2:21: incomplete type: array of array of 8 int (type.go:35:typ: type.go:46:typ0: type.go:318:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79737-1.c`:                  {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82387.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((ccv2) && ((ccv1 != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82210.c`:                    {}, // COMPILE FAIL: pr82210.c:18:9: unsupported alignment 16 of struct S {a array of struct T; b array of int} (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr80692.c`:                    {}, // COMPILE FAIL: TODO *cc.PredefinedType _Decimal64 _Decimal64 (expr.go:2608:primaryExpressionFloatConst: type.go:29:helper: type.go:128:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84169.c`:                    {}, // COMPILE FAIL: pr84169.c:4:27: unsupported alignment 16 of T (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84521.c`:                    {}, // COMPILE FAIL: pr84521.c.go:318:2: undefined reference to '__builtin_longjmp' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84748.c`:                    {}, // COMPILE FAIL: pr84748.c:3:27: unsupported alignment 16 of u128 (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85529-1.c`:                  {}, // COMPILE FAIL: TODO "pp_ = ((Xs.fda) != iqlibc.ppBoolInt32(((ccv3) && ((ccv1 != 0)))))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr87053.c`:                    {}, // COMPILE FAIL: TODO (init.go:267:initializerUnion: type.go:35:typ: type.go:224:typ0)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85582-3.c`:                  {}, // COMPILE FAIL: pr85582-3.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85582-2.c`:                  {}, // COMPILE FAIL: pr85582-2.c:4:18: unsupported alignment 16 of S (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr86528.c`:                    {}, // COMPILE FAIL: pr86528.c.go:319:36: undefined reference to '__builtin_alloca' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr89195.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93434.c`:                    {}, // COMPILE FAIL: TODO "(!((ani != 0)))" int exprBool -> double exprDefault (expr.go:778:equalityExpression: expr.go:79:expr: expr.go:149:convert)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93213.c`:                    {}, // COMPILE FAIL: pr93213.c:8:27: unsupported alignment 16 of u128 (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93249.c`:                    {}, // COMPILE FAIL: pr93249.c.go:325:2: undefined reference to '__builtin_strncpy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94724.c`:                    {}, // COMPILE FAIL: TODO "pp_ = ((iqlibc.ppInt32FromInt16(Xb)) != (iqlibc.ppInt32FromInt32(53601)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr96549.c`:                    {}, // COMPILE FAIL: gcc.c-torture/execute/pr96549.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94809.c`:                    {}, // COMPILE FAIL: TODO "pp_ = (((-((iqlibc.ppUint64FromUint64(1)))) / anone) < (iqlibc.ppUint64FromInt32(ccv1)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr92904.c`:                    {}, // COMPILE FAIL: pr92904.c:12:10: unsupported alignment 16 of __int128 (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr97325.c`:                    {}, // COMPILE FAIL: pr97325.c.go:318:10: undefined reference to '__builtin_ffs' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr98366.c`:                    {}, // COMPILE FAIL: TODO bitfield (init.go:130:initializerArray: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr98474.c`:                    {}, // COMPILE FAIL: pr98474.c:4:21: unsupported alignment 16 of T (type.go:17:typedef: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/printf-2.c`:                   {}, // COMPILE FAIL: printf-2.c.go:4476:8: undefined reference to 'freopen' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stdarg-3.c`:                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stdarg-1.c`:                   {}, // COMPILE FAIL: stdarg-1.c.go:448:2: undefined reference to '__builtin_va_copy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strct-stdarg-1.c`:             {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strct-varg-1.c`:               {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/string-opt-18.c`:              {}, // COMPILE FAIL: string-opt-18.c.go:330:5: undefined reference to 'mempcpy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/struct-ini-2.c`:               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/struct-ini-3.c`:               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-10.c`:                  {}, // COMPILE FAIL: va-arg-10.c.go:384:2: undefined reference to '__builtin_va_copy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-14.c`:                  {}, // COMPILE FAIL: va-arg-14.c.go:372:2: undefined reference to '__builtin_va_copy' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-5.c`:                   {}, // COMPILE FAIL: TODO (init.go:18:initializerOuter: init.go:72:initializer: init.go:236:initializerUnion)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/user-printf.c`:                {}, // COMPILE FAIL: user-printf.c.go:4512:15: undefined reference to 'tmpnam' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-22.c`:                  {}, // COMPILE FAIL: va-arg-22.c:24:1: invalid type size: 0 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-pack-1.c`:              {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zero-struct-1.c`:              {}, // COMPILE FAIL: zero-struct-1.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zero-struct-2.c`:              {}, // COMPILE FAIL: zero-struct-2.c:3:19: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0011-switch1.c`:                           {}, // COMPILE FAIL: mir/c-tests/andrewchambers_c/0011-switch1.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0015-calls13.c`:                           {}, // COMPILE FAIL: 0015-calls13.c:1:1: invalid type size: -1 (type.go:338:structLiteral: type.go:46:typ0: type.go:323:checkValidType)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0025-duff.c`:                              {}, // COMPILE FAIL: mir/c-tests/andrewchambers_c/0025-duff.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits14.c`:                           {}, // COMPILE FAIL: TODO (type.go:338:structLiteral: type.go:190:typ0: type.go:247:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits12.c`:                           {}, // COMPILE FAIL: TODO (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:224:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits13.c`:                           {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1752:postfixExpressionSelect)
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits15.c`:                           {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1752:postfixExpressionSelect)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/anonymous-members.c`:                                  {}, // COMPILE FAIL: TODO (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:224:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/anonymous-struct.c`:                                   {}, // COMPILE FAIL: TODO PostfixExpressionSelect (expr.go:463:expr0: expr.go:1343:postfixExpression: expr.go:1752:postfixExpressionSelect)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-basic.c`:                                     {}, // COMPILE FAIL: bitfield-basic.c:6:3: unsupported alignment 4 of union A {int; b char} (type.go:332:unionLiteral: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-immediate-bitwise.c`:                         {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-initialize-zero.c`:                           {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-load.c`:                                      {}, // COMPILE FAIL: TODO bitfield (init.go:289:initializerUnion: type.go:35:typ: type.go:217:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-pack-next.c`:                                 {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-trailing-zero.c`:                             {}, // COMPILE FAIL: bitfield-trailing-zero.c:15:1: unsupported alignment 4 of union U1 {int; int; c char; int; int; int; int} (type.go:332:unionLiteral: type.go:46:typ0: type.go:302:checkValidType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-packing.c`:                                   {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-reset-align.c`:                               {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-types.c`:                                     {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-types-init.c`:                                {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield.c`:                                           {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/deref-compare-float.c`:                                {}, // COMPILE FAIL: TODO "((iqlibc.ppFloat64FromFloat64(0)) > (iqlibc.ppFloat64FromInt32((*(*ppint32)(iqunsafe.ppPointer(Xi))))))" int exprBool -> float exprDefault (expr.go:2027:assignmentExpression: expr.go:79:expr: expr.go:149:convert)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/duffs-device.c`:                                       {}, // COMPILE FAIL: mir/c-tests/lacc/duffs-device.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/float-compare-equal.c`:                                {}, // COMPILE FAIL: TODO "((iqlibc.ppFloat32FromInt32(Xi)) == (*(*ppfloat32)(iqunsafe.ppPointer(Xf))))" int exprBool -> float exprDefault (expr.go:2027:assignmentExpression: expr.go:79:expr: expr.go:149:convert)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/identifier.c`:                                         {}, // COMPILE FAIL: invalid object file: multiple defintions of c (ccgo.go:255:Main: link.go:225:link: link.go:436:getFileSymbols)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/initialize-call.c`:                                    {}, // COMPILE FAIL: TODO "sigetobj(aatls, (iqlibc.ppInt32FromInt32(2)), (iqlibc.ppInt32FromInt32(8)), (iqlibc.ppInt32FromInt32(1)))" struct obj {s array of 3 short} exprDefault -> struct {t struct obj} exprDefault (initialize-call.c:25:9:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/initialize-string.c`:                                  {}, // COMPILE FAIL: TODO "\"wat\\x00\"" pointer to char exprDefault -> array of 7 char exprDefault (initialize-string.c:14:14:) (expr.go:79:expr: expr.go:133:convert: expr.go:335:convertType)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/return-bitfield.c`:                                    {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/short-circuit-comma.c`:                                {}, // COMPILE FAIL: TODO "pp_ = ((ccv1) && (((iqlibc.ppInt32FromInt32(1)) != 0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/string-addr.c`:                                        {}, // COMPILE FAIL: TODO (expr.go:465:expr0: expr.go:2350:primaryExpression: expr.go:2467:primaryExpressionStringConst)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/stringify.c`:                                          {}, // COMPILE FAIL: TODO `'\xA'` -> invalid syntax (expr.go:465:expr0: expr.go:2346:primaryExpression: expr.go:2540:primaryExpressionCharConst)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/tail-compare-jump.c`:                                  {}, // COMPILE FAIL: TODO "pp_ = ((iqlibc.ppInt64FromInt32(Xf)) < (iqlibc.ppInt64FromInt64(0)))" int exprBool -> int exprVoid (expr.go:79:expr: expr.go:129:convert: expr.go:286:convertMode)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/union-bitfield.c`:                                     {}, // COMPILE FAIL: TODO bitfield (type.go:419:defineUnion: type.go:332:unionLiteral: type.go:217:typ0)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg-complex-1.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg-complex-2.c`:                                   {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg.c`:                                             {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/new/issue142.c`:                                            {}, // COMPILE FAIL: unsupported va_arg type: struct (expr.go:2354:primaryExpression: expr.go:471:expr0: expr.go:1036:unaryExpression)
	`assets/github.com/vnmakarov/mir/c-tests/new/issue117.c`:                                            {}, // COMPILE FAIL: issue117.c:2:3: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/76_dollars_in_identifiers.c`:                                        {}, // COMPILE FAIL: 76_dollars_in_identifiers.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/tcc-0.9.27/tests/tests2/80_flexarray.c`:                                                     {}, // COMPILE FAIL: 80_flexarray.c:5:1: incomplete type: array of int (type.go:190:typ0: type.go:46:typ0: type.go:318:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/87_dead_code.c`:                                                     {}, // COMPILE FAIL: 87_dead_code.c:17:9: TODO false (expr.go:465:expr0: expr.go:2360:primaryExpression: stmt.go:107:compoundStatement)
	`assets/tcc-0.9.27/tests/tests2/73_arm64.c`:                                                         {}, // COMPILE FAIL: TODO (expr.go:71:expr: expr.go:463:expr0: expr.go:1331:postfixExpression)
	`assets/tcc-0.9.27/tests/tests2/88_codeopt.c`:                                                       {}, // COMPILE FAIL: 88_codeopt.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/tcc-0.9.27/tests/tests2/85_asm-outside-function.c`:                                          {}, // COMPILE FAIL: 85_asm-outside-function.c.go:315:2: undefined reference to 'vide' (ccgo.go:255:Main: link.go:246:link: link.go:614:link)
	`assets/tcc-0.9.27/tests/tests2/89_nocode_wanted.c`:                                                 {}, // COMPILE FAIL: 89_nocode_wanted.c: gofmt: exit status 2 (asm_amd64.s:1594:goexit: compile.go:350:compile: compile.go:350:compile)
	`assets/tcc-0.9.27/tests/tests2/93_integer_promotion.c`:                                             {}, // COMPILE FAIL: TODO bitfield (init.go:18:initializerOuter: init.go:66:initializer: init.go:151:initializerStruct)
	`assets/tcc-0.9.27/tests/tests2/94_generic.c`:                                                       {}, // COMPILE FAIL: TODO PrimaryExpressionGeneric (expr.go:71:expr: expr.go:465:expr0: expr.go:2369:primaryExpression)
	`assets/tcc-0.9.27/tests/tests2/90_struct-init.c`:                                                   {}, // COMPILE FAIL: 90_struct-init.c:2:19: invalid type size: -1 (type.go:17:typedef: type.go:46:typ0: type.go:323:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/95_bitfields.c`:                                                     {}, // COMPILE FAIL: 95_bitfields.c:27:5: unsupported alignment 16 of struct __s {x int; y char; z long long; a char; b long long} (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
	`assets/tcc-0.9.27/tests/tests2/95_bitfields_ms.c`:                                                  {}, // COMPILE FAIL: 95_bitfields.c:27:5: unsupported alignment 16 of struct __s {x int; y char; z long long; a char; b long long} (type.go:35:typ: type.go:46:typ0: type.go:302:checkValidType)
}
