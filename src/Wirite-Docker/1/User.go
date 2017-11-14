package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWUSER,
	}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(1), Gid: uint32(1)}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)

	}
	os.Exit(-1)
}

/*
我们在原来的基础上增加了syscall.CLONE_NEWUSER。首先我们以root来运行这个程序，运行前在宿主机上我们看一下当前用户和用户组

root@iZ254rt8xf1Z:~/gocode/src/book# id
uid=0(root) gid=0(root) groups=0(root)

可以看到我们是root 用户，我们运行一下程序

root@iZ254rt8xf1Z:~/gocode/src/book# go run main.go
$ id
uid=65534(nobody) gid=65534(nogroup) groups=65534(nogroup)


*/
