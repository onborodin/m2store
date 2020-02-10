
package main

import (
    "path/filepath"
    "fmt"
    "os"
    "flag"
    "strings"
    "encoding/json"
    "errors"

    "store/server/user-model"
    "store/config"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

)

func main() {

    config := config.New()
    config.Read()

    optFileName := flag.String("file", config.PasswordPath, "password file")

    updateCommands := flag.NewFlagSet("update", flag.ExitOnError)

    optUpdateName := updateCommands.String("name", "", "user name")
    optUpdatePass := updateCommands.String("pass", "", "user key")
    optUpdateAdmin := updateCommands.Bool("admin", false, "user is admin")

    authCommands := flag.NewFlagSet("auth", flag.ExitOnError)
    optAuthName := authCommands.String("name", "", "user name")
    optAuthPass := authCommands.String("pass", "", "user key")

    deleteCommands := flag.NewFlagSet("delete", flag.ExitOnError)
    optDeleteName := deleteCommands.String("name", "", "user name")

    exeName := filepath.Base(os.Args[0])

    flag.Usage = func() {
        fmt.Printf("usage: %s [global option] command [command option]\n", exeName)
        fmt.Println("")
        fmt.Println("commands: migrate, create, update, auth, delete, list")
        fmt.Println("")

        fmt.Println("global option:")
        fmt.Println("")
        flag.PrintDefaults()

        fmt.Println("")
        fmt.Println("create|update option:")
        updateCommands.PrintDefaults()

        fmt.Println("")
        fmt.Println("auth option:")
        authCommands.PrintDefaults()

        fmt.Println("")
        fmt.Println("delete option:")
        deleteCommands.PrintDefaults()

        os.Exit(0)
    }

    flag.Parse()

    localArgs := flag.Args()
    var command string

    if len(localArgs) == 0 {
        flag.Usage()
    }

    command = localArgs[0]
    localArgs = localArgs[1:]
    var err error


    dbUrl := fmt.Sprintf("%s", *optFileName)

    db, err := sqlx.Open("sqlite3", dbUrl)
    if err != nil {
        fmt.Printf("error: %s\n", err)
        os.Exit(1)
    }

    /* Check DB connection */
    err = db.Ping()
    if err != nil {
        fmt.Printf("error: %s\n", err)
        os.Exit(1)
    }

    user := userModel.New(db)

    if strings.HasPrefix(command, "mig") {
        err = user.Migrate()

    } else if strings.HasPrefix(command, "cre") {

        updateCommands.Parse(localArgs)

        if len(*optUpdateName) < 5 || len(*optUpdatePass) < 5 {
            err = errors.New("username or passwotrd is less 5 chars")
        } else {
            _user := userModel.User{
                Username: *optUpdateName,
                Password: *optUpdatePass,
                IsAdmin:  *optUpdateAdmin,
            }
            err = user.Create(_user)
        }

    } else if strings.HasPrefix(command, "upd") {

        updateCommands.Parse(localArgs)

        _user := userModel.User{
            Username: *optUpdateName,
        }

        res, err := user.Find(_user)
        if err == nil {
            _user.Id = res.Id
            _user.Username = *optUpdateName
            _user.Password = *optUpdatePass
            _user.IsAdmin = *optUpdateAdmin

            err = user.Update(_user)
        }

    /* Delete user */
    } else if strings.HasPrefix(command, "del") {

        deleteCommands.Parse(localArgs)
        _user := userModel.User{
            Username: *optDeleteName,
        }

        res, err := user.Find(_user)
        if err == nil {
            _user.Id = res.Id
            err = user.Delete(_user)
        }

    } else if strings.HasPrefix(command, "aut") {

        authCommands.Parse(localArgs)

        _user := userModel.User{
            Username: *optAuthName,
            Password: *optAuthPass,
        }
        err = user.Check(&_user)

    } else if strings.HasPrefix(command, "lis") {
        _page := userModel.Page{
            Limit: 1000,
            Offset: 0,
        }
        _ = user.List(&_page)
        for _, item := range *_page.Users {
            printUser(item)
        }
    } else {
        flag.Usage()
    }

    if err != nil {
        printError(err)
        os.Exit(1)
    } else {
        printOk()
        os.Exit(0)
    }

}

func printUser(user userModel.User) {
    data, _ := json.Marshal(user)
    fmt.Println(string(data))
}

func printError(err error) {
    fmt.Printf(`{"error":true, "message":"%s"}` + "\n", err)
}

func printOk() {
    fmt.Printf(`{"error":false}` + "\n")
}
