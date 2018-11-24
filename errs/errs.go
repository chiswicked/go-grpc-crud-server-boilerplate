package errs

import (
	"fmt"
	"log"
	"os"
)

// FatalIf terminates program execution with error message if err is set
// Deferred functions are not run
func FatalIf(msg string, err error) {
	if err != nil {
		log.Fatalf(msg+": %s\n", err)
	}
}

// PanicIf terminates program execution with error message and stack trace if err is set
// Deferred functions are run
func PanicIf(msg string, err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, msg+"\n")
		panic(err)
	}
}
