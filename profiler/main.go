package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"cloud.google.com/go/profiler"
)

func main() {
	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(profiler.Config{
		Service:        "myservice",
		ServiceVersion: "1.0.1",
		ProjectID:      "herman-codam-tmp-nov-2022",
	}); err != nil {
		log.Fatal(err)
	}

	rapport := time.NewTicker(5 * time.Second)
	t := time.NewTicker(10 * time.Millisecond)
	_ = t
	iterations := 0
	for {
		select {
		case <-rapport.C:
			fmt.Println("Iterations", iterations)
			iterations = 0
		// case <-t.C:
		default:
			busyWork()
			iterations++
		}
	}
}

func busyWork() {
	for i := 0; i < 1000; i++ {
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

// TODO uncomment & remove line 63 to improve performance 7%
// var re = regexp.MustCompile("(a*)*c")

func busyWorkOther() {
	re := regexp.MustCompile("(a*)*c")
	if re.MatchString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") {
		_ = 1
	}
}

// But remember:
// We should forget about small efficiencies, say about 97% of the time:
// premature optimization is the root of all evil. – Donald Knuth
