// build +main

// Demonstrate that computed args to defer'ed or go'ed functions evalate immediately.
// Run several times to see how nondeterministic goroutine timing is.
package main

import "time"

func g(i int) int { println("g", i); return i + 1 }
func h(i int) int { println("h", i); return i + i }
func fn(a, b int) { println("fn", a, b) }

func main() {
	defer fn(g(10), h(20))
	println("ok")

	for i := 0; i < 30; i++ {
		println("going", i)
		go fn(g(1000*i), h(10000*i))
		println("gone", i)
	}
	time.Sleep(1 * time.Second)
	println("bye")
}
