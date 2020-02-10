/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "store/client"
    "fmt"
    "flag"
    "os"
    "path/filepath"
    "strings"
)

func main() {

    optNode := flag.String("node", "localhost:8080", "node set")
    optUserName := flag.String("user", "user1", "username")
    optPassword := flag.String("pass", "12345", "password")

        //node
    listCommands := flag.NewFlagSet("list", flag.ExitOnError)
        optListBucket := listCommands.String("bucket", "", "bucket name")
        optListPattern := listCommands.String("pattern", "*", "name pattern")

    putCommands := flag.NewFlagSet("put", flag.ExitOnError)
        optPutBucket := putCommands.String("bucket", "", "bucket name")
        optPutFileName := putCommands.String("file", "", "file name")

    getCommands := flag.NewFlagSet("get", flag.ExitOnError)
        optGetBucket := getCommands.String("bucket", "", "bucket name")
        optGetFileName := getCommands.String("file", "", "file name")

    dropCommands := flag.NewFlagSet("drop", flag.ExitOnError)
        optDropBucket := dropCommands.String("bucket", "", "bucket name")
        optDropFileName := dropCommands.String("file", "", "file name")

    listBucketsCommands := flag.NewFlagSet("listb", flag.ExitOnError)

    exeName := filepath.Base(os.Args[0])
    flag.Usage = func() {
        fmt.Printf("usage: %s [global option] command [command option]\n", exeName)

        fmt.Println("")
        fmt.Println("commands: list, put, get, drop, listb")
        fmt.Println("")

        fmt.Println("global option:")
        flag.PrintDefaults()
        fmt.Println("")

        fmt.Println("list option:")
        listCommands.PrintDefaults()
        fmt.Println("")

        fmt.Println("put option:")
        putCommands.PrintDefaults()
        fmt.Println("")

        fmt.Println("get option:")
        getCommands.PrintDefaults()
        fmt.Println("")

        fmt.Println("drop option:")
        dropCommands.PrintDefaults()
        fmt.Println("")

        fmt.Println("list option:")
        listCommands.PrintDefaults()
        fmt.Println("")

        fmt.Println("listb option:")
        listBucketsCommands.PrintDefaults()
        fmt.Println("")
    }

    flag.Parse()
    fmt.Println("node:", *optNode)

    localArgs := flag.Args()
    if len(localArgs) == 0 {
        flag.Usage()
        os.Exit(1)
    }

    command := localArgs[0]
    fmt.Println("command: ", command)

    if len(localArgs) < 1 {
        flag.Usage()
        os.Exit(1)
    }

    localArgs = localArgs[1:]

    if strings.HasPrefix(command, "list") {

        listCommands.Parse(localArgs)
        client := client.New()
        res, err := client.List(*optNode, *optUserName, *optPassword, *optListBucket, *optListPattern)
        if err != nil {
            fmt.Println("error:", err)
            os.Exit(1)
        }
        fmt.Println(res)

    } else if strings.HasPrefix(command, "put") {

        putCommands.Parse(localArgs)
        client := client.New()
        res, err := client.Put(*optNode, *optUserName, *optPassword, *optPutBucket, *optPutFileName)
        if err != nil {
            fmt.Println("error:", err)
            os.Exit(1)
        }
        fmt.Println(res)

    } else if strings.HasPrefix(command, "get") {

        getCommands.Parse(localArgs)
        client := client.New()
        res, err := client.Get(*optNode, *optUserName, *optPassword, *optGetBucket, *optGetFileName)
        if err != nil {
            fmt.Println("error:", err)
            os.Exit(1)
        }
        fmt.Println(res)
    } else if strings.HasPrefix(command, "drop") {

        dropCommands.Parse(localArgs)
        client := client.New()
        res, err := client.Drop(*optNode, *optUserName, *optPassword, *optDropBucket, *optDropFileName)
        if err != nil {
            fmt.Println("error:", err)
            os.Exit(1)
        }
        fmt.Println(res)
    }


}
