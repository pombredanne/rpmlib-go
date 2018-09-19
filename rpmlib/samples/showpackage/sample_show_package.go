// Copyrigt 2018 necomeshi, All Rights Reserved.

package main

import (
	"fmt"
	"os"

	"github.com/necomeshi/rpmlib"
)

func main() {

	if len(os.Args) == 0 {
		fmt.Fprintf(os.Stderr, "No package name is specified")
		os.Exit(1)
	}

	ts, err := rpmlib.NewTransactionSet()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get transaction set: %s\n", err)
		os.Exit(1)
	}

	defer ts.Free()

	h, err := ts.ReadPackageFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get package rpm header: %s\n", err)
		os.Exit(1)
	}

	defer h.Free()

	name, err := h.GetStringArray(rpmlib.RPMTAG_NAME)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get package name: %s\n", err)
	}

	fmt.Println(name)

}
