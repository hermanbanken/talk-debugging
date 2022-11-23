package util

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func DoBusyWork(iterations int) {
	for i := 0; i < iterations; i++ {
		if i%2 == 0 {
			busyWorkNested()
		} else {
			busyWorkOther()
		}
	}
}

func busyWorkNested() {
	var dst interface{}
	_ = json.Unmarshal([]byte(`[[[[[[[[[["nested string"]]]]]]]]]]`), &dst)
	// do something with dst to avoid optimization
	if len(fmt.Sprint(dst)) > 0 {
		return
	}
}

func busyWorkOther() {
	if regexp.MustCompile("(a*)*c").MatchString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") {
		_ = 1
	}
}

// But remember:
// We should forget about small efficiencies, say about 97% of the time:
// premature optimization is the root of all evil. – Donald Knuth
