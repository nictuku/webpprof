package ppcommon

import (
	"testing"
)

var src = `heap profile: 1: 288 [1470: 87488] @ heap/1048576
1: 288 [1: 288] @ 0x423e45 0x423928 0x425071 0x4250bb 0x412125 0x41446e 0x40e5a6 0x4115e5 0x410f18 0x42ad64
#       0x423e45        profilealloc+0xb5       /usr/local/go/src/pkg/runtime/malloc.goc:258
#       0x423928        runtime.mallocgc+0x218  /usr/local/go/src/pkg/runtime/malloc.goc:197
#       0x425071        cnew+0xc1               /usr/local/go/src/pkg/runtime/malloc.goc:836
#       0x4250bb        runtime.cnew+0x3b       /usr/local/go/src/pkg/runtime/malloc.goc:843
#       0x412125        allocg+0x55             /usr/local/go/src/pkg/runtime/proc.c:701
#       0x41446e        runtime.malg+0x1e       /usr/local/go/src/pkg/runtime/proc.c:1770
#       0x40e5a6        runtime.mpreinit+0x26   /usr/local/go/src/pkg/runtime/os_linux.c:198
#       0x4115e5        mcommoninit+0xb5        /usr/local/go/src/pkg/runtime/proc.c:361
#       0x410f18        runtime.schedinit+0x48  /usr/local/go/src/pkg/runtime/proc.c:151
#       0x42ad64        _rt0_go+0x114           /usr/local/go/src/pkg/runtime/asm_amd64.s:91

`

var result = `--- symbol
binary=unknown
0x423e45 profilealloc
0x425071 cnew
0x4250bb runtime.cnew
0x412125 allocg
0x40e5a6 runtime.mpreinit
0x4115e5 mcommoninit
0x42ad64 _rt0_go
0x423928 runtime.mallocgc
0x41446e runtime.malg
0x410f18 runtime.schedinit
---
heap profile: 1: 288 [1470: 87488] @ heap/1048576
1: 288 [1: 288] @ 0x423e45 0x423928 0x425071 0x4250bb 0x412125 0x41446e 0x40e5a6 0x4115e5 0x410f18 0x42ad64

`

func TestRawProfile(t *testing.T) {

	_, err := RawProfile(src)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Don't depend on map iteration order so we can compare tehe results.
	// if raw != result {
	//	t.Errorf("Raw profile not in the expected result: ====\n%v\n====\n%v\n", raw, result)
	// }
}
