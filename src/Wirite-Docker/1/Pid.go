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
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)

	}
}

/*
root@hsukk-VirtualBox:/home/hsukk/gocode# go run Uts.go
# pstree -pl

           ├─sshd(873)───sshd(2174)───sshd(2254)─┬─bash(2257)───su(6210)───bash(6211)───go(6245)─┬─Uts(6263)─┬─bash(6266)───go(27794)─┬─Uts(27812)─┬─sh(27815)───pstree(27816)
           │                                     │                                               │           │                        │            ├─{Uts}(27813)
           │                                     │                                               │           │                        │            └─{Uts}(27814)
           │                                     │                                               │           │                        ├─{go}(27795)
           │                                     │                                               │           │                        ├─{go}(27796)
           │                                     │                                               │           │                        ├─{go}(27797)
           │                                     │                                               │           │                        └─{go}(27800)
           │                                     │                                               │           ├─{Uts}(6264)
           │                                     │                                               │           └─{Uts}(6265)
           │                                     │                                               ├─{go}(6246)
           │                                     │                                               ├─{go}(6247)
           │                                     │                                               ├─{go}(6248)
           │                                     │                                               └─{go}(6251)
           │                                     └─bash(4223)
# echo $$
27815
# cat /proc/27815/uts
cat: /proc/27815/uts: No such file or directory
# readlink /proc/27815/ns/uts
uts:[4026532332]
# readlink /proc/6211/ns/uts
uts:[4026531838]
# readlink /proc/6254/ns/uts
# readlink /proc/6245/ns/uts
uts:[4026531838]

可以看到uts空间所处不同。
*/

/*
 hostname  证明封闭空间

 T1:  hostname testgo
 T2: sudo hostname testbox
 T3: hostname   (来做验证，看看哪个被更改，哪个是封闭的不影响)

*/
