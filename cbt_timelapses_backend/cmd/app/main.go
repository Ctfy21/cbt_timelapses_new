package main

import (
	"fmt"
	"os/exec"
)

func main() {

	ch := make(chan string)

	go getList(ch)

	fmt.Println("After goroutine!")

	fmt.Println(<-ch)

}

func getList(ch chan string) {
	out, err := exec.Command("echo", "Hello World!").Output()

	if err != nil {
		panic(err)
	}

	ch <- string(out)
}
