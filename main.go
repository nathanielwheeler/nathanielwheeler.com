package main

import (
	"fmt"
	"os"

	"nathanielwheeler.com/server"
)

/* TODO
- Decouple User middleware from routes (like with RequireUser)
*/

func main() {
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "server error:\n\t%s\n", err)
		os.Exit(1)
	}
}
