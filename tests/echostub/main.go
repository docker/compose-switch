package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func main() {
	if reflect.DeepEqual(os.Args, []string{"docker", "compose", "version", "--short"}) {
		fmt.Println("2.0.0-rc.3")
	} else {
		fmt.Println(strings.Join(os.Args, " "))
	}
}
