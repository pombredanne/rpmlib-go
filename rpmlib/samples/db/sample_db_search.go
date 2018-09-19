// Copyrigt 2018 necomeshi, All Rights Reserved.

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/necomeshi/rpmlib"
)

func main() {
	ts, err := rpmlib.NewTransactionSet()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get transaction set: %s\n", err)
		os.Exit(1)
	}

	defer ts.Free()

	iter, err := ts.SequencialIterator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get iterator: %s\n", err)
		os.Exit(1)
	}

	defer iter.Free()

	for {
		h, itr_err := iter.Next()

		if itr_err == io.EOF {
			break
		}

		if itr_err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get next iterator: %s\n", itr_err)
			break
		}

		defer h.Free()

		name, h_err := h.GetString(rpmlib.RPMTAG_NAME)
		if h_err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get name from rpm header: %s\n", h_err)
			continue
		}

	     version,h_err := h.GetString(rpmlib.RPMTAG_VERSION)
	    if h_err != nil {
	      version = string("unknown")
	    }

	    release,h_err := h.GetString(rpmlib.RPMTAG_RELEASE)
	    if h_err != nil {
	      release = string("unknown")
	    }
	    
	    pktname := name + string('-') + version + string('-') + release
	  
	    fmt.Println(pktname)


	}
}
