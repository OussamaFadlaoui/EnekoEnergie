package helpers

import (
	"fmt"
	"os"
)

func Check(err error, msg ... string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(msg)
		os.Exit(1)
	}
}
