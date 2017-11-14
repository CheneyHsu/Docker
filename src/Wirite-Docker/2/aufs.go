package main

import (
	"os/exec"
	"fmt"
)


func main() {
cmd := exec.Command("/bin/bash","./exec.sh")
bytes,err :=cmd.Output()
if err != nil{
	fmt.Println("cmd.Output:",err)
	return
}
fmt.Println(string(bytes))
}
