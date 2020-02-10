/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package daemon

import (
    "log"
    "os"
    "io"
    "syscall"
    "strconv"
    "os/signal"
)

func SaveProcessID(filename string) error {
    pid := os.Getpid()
    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0640)
    if err != nil {
         return err
    }
    defer file.Close()
    if _, err := file.WriteString(strconv.Itoa(pid)); err != nil {
        return err
    }
    file.Sync()
    return nil
}

func RedirectLog(filename string, debug bool) (*os.File, error) {
    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
    if err != nil {
        return nil, err
    }
    wrt := io.MultiWriter(os.Stdout, file)
    if debug {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    } else {
        log.SetFlags(log.LstdFlags)
    }
    log.SetOutput(wrt)
    return file, nil
}

func ForkProcess() error {

    if _, exists := os.LookupEnv("GOGOFORK"); !exists {
        os.Setenv("GOGOFORK", "yes")

        cwd, err := os.Getwd()
        if err != nil {
            return err
        }

        procAttr := syscall.ProcAttr{}
        procAttr.Files = []uintptr{ uintptr(syscall.Stdin), uintptr(syscall.Stdout), uintptr(syscall.Stderr) }
        procAttr.Env = os.Environ()
        procAttr.Dir = cwd
        syscall.ForkExec(os.Args[0], os.Args, &procAttr)
        os.Exit(0)
    }
    _, err := syscall.Setsid()
    if err != nil {
        return err
    }
    os.Chdir("/")
    return nil
}

func RedirectIO() (*os.File, error) {
    file, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
    if err != nil {
        return nil, err
    }
    syscall.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
    syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
    syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
    return file, nil
}

func SetSignalHandler() {

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    go func() {
        for {
            log.Printf("signal handler start")
            sig := <- sigs
            log.Printf("receive signal %s", sig.String())

            switch sig {
                case syscall.SIGINT, syscall.SIGTERM:
                    log.Printf("exit process by signal %s", sig.String())
                    os.Exit(0)

                case syscall.SIGHUP:
                    log.Printf("restart program")
                    ForkProcess()
            }
        }
    }()
}
