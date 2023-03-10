set -e
rm -f log-ccgo
make clean || true
make distclean || true
go install -v modernc.org/ccgo/v2/ccgo
./configure CC=ccgo \
	CFLAGS='--ccgo-full-paths --ccgo-struct-checks --ccgo-use-import exec.ErrNotFound,os.DevNull,atomic.Value{} -D_GNU_SOURCE' \
	LDFLAGS='--warn-unresolved-libs --ccgo-go --ccgo-import os,os/exec,sync/atomic'
make $MAKEJ binaries
make $MAKEJ test
go version
date

# all.tcl:	Total	25904	Passed	1530	Skipped	943	Failed	23431	# custom match fail unpatched
# all.tcl:	Total	25904	Passed	1530	Skipped	884	Failed	23490	# with -DTCL_MEM_DEBUG
# all.tcl:	Total	13453	Passed	12580	Skipped	831	Failed	42 	# removed -DTCL_MEM_DEBUG, trying musl memory allocator
# all.tcl:	Total	15687	Passed	14883	Skipped	798	Failed	6	# musl memory allocator with -DTCL_MEM_DEBUG
# all.tcl:	Total	25904	Passed	24963	Skipped	884	Failed	57	# Fixed vdso clock_gettime
# all.tcl:	Total	14925	Passed	13922	Skipped	909	Failed	94	# Removed -DTCL_MEM_DEBUG
# all.tcl:	Total	25959	Passed	24919	Skipped	944	Failed	96	# dtto
# all.tcl:	Total	26037	Passed	25028	Skipped	944	Failed	65	# 2018-11-13
# all.tcl:	Total	26092	Passed	25050	Skipped	945	Failed	97
# all.tcl:	Total	26037	Passed	25028	Skipped	944	Failed	65	# 2018-11-14
# all.tcl:	Total	26092	Passed	25050	Skipped	945	Failed	97	# 2018-11-16
# all.tcl:	Total	26037	Passed	25028	Skipped	944	Failed	65	# 2018-11-19 e5-1650 go version go1.11.2 linux/amd64 openSuse Leap 15.0
# all.tcl:	Total	25651	Passed	24771	Skipped	815	Failed	65	# 2018-11-26 e5-1650 go version go1.11.2 linux/amd64 openSuse Leap 15.0
# all.tcl:	Total	26037	Passed	25028	Skipped	944	Failed	65	# 2018-11-27 e5-1650 go version go1.11.2 linux/amd64 openSuse Leap 15.0
# all.tcl:	Total	26037	Passed	25028	Skipped	944	Failed	65	# 2018-11-30 e5-1650 go version go1.11.2 linux/amd64 openSuse Leap 15.0
# all.tcl:	Total	26092	Passed	25050	Skipped	945	Failed	97	# 2018-11-30 4670 go version go1.11.2 linux/amd64 openSuse 42.3
# all.tcl:	Total	26037	Passed	25028	Skipped	944	Failed	65	# 2019-01-31 e5-1650 go version go1.11.4 linux/amd64 openSuse Leap 15.0

# all.tcl:	Total	26092	Passed	25089	Skipped	945	Failed	58	# 2019-02-03 4670 go version go1.11.2 linux/amd64 openSuse 42.3
# all.tcl:	Total	26092	Passed	25089	Skipped	945	Failed	58	# 2019-02-05 e5-1650 1.11.4 linux/amd64 openSuse Leap 15.0
