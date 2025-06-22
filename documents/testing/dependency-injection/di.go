package dependencyinjection

import (
	"bytes"
	"fmt"
)

// clearing up my misunderstanding about func param injection vs interface injection
// The Key Insight
// Function injection is about injecting behavior/logic, while interface injection is more about injecting capabilities/services.

// Function: "Here's HOW to do something"
// Interface: "Here's SOMETHING that can do a job"

// Dependency injection main idea
// The pattern is: Create the flexible version first, then add convenient wrappers.

// type Writer interface {
// 	Write([]byte) (int, error)
// }

func Greet(writer *bytes.Buffer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}
