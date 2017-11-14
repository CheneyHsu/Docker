package main
import (
log "github.com/Sirupsen/logrus"
"github.com/urfave/cli"
"os"
	"fmt"
	"os/exec"
	"syscall"
)

const usage = `mydocker is a simple container runtime implementation.`

func main() {
	app := cli.NewApp()
	app.Name ="mydocker"
	app.Usage=usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error{
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)

		return nil
	}

	if err := app.Run(os.Args); err != nil{
		log.Fatal(err)
	}
}

//这里定义了 runCommand的Flags，类似于我们运行命令的时候--指定的参数
var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit
            mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "ti",
			Usage:       "enable tty",
		},
	},
	/*
	这里是run命令执行的真正函数。
	1. 判断参数是否包含command
	2. 获取用户指定的command
	3. 调用Run function 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

//这里定义了initCommand的具体操作，此操作为内部方法，禁止外部调用
var initCommand = cli.Command{
	Name:   "init",
	Usage:  "Init container process run user's process in container. Do not call it outside",
	/*
	1. 获取传递过来的command参数
	2. 执行容器初始化操作
	*/
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

/*
这里是父进程也就是我们当前进程执行的内容，根据我们上一章介绍的内容，应该比较容易明白
1.这里的/proc/self/exe 调用，其中/proc/self指的是当前运行进程自己的环境，exec其实就是自己调用了自己，我们使用这种方式实现对创建出来的进程进行初始化
2.后面args是参数，其中 init 是传递给本进程的第一个参数，这在本例子中，其实就是会去调用我们的initCommand去初始化进程的一些环境和资源
3. 下面的clone 参数就是去 fork 出来的一个新进程，并且使用了namespace隔离新创建的进程和外部的环境。
4. 如果用户指定了-ti 参数，我们就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

/*
这里的Start方法是真正开始前面创建好的command的调用，他会首先 clone
出来一个 namespace 隔离的进程，然后在子进程中，调用/proc/self/exe 也就是自己，发送 init 参数，调用我们写的init方法，去初始化容器的一些资源
*/
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}

/*
这里的init函数执行是在容器内部的，也就是说，代码执行到这里后，其实容器所在的进程已经创建出来了，我们是本容器执行的第一个进程。
1.使用mount 先去挂载proc 文件系统，方便我们通过ps等系统命令去查看当前进程资源情况
*/
func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}