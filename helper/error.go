package helper

import "fmt"

func IfError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
