
package bucketController

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/gin-gonic/gin"

    "store/config"
    "store/tools"
)

const (
    MaxBucketDepth int = 64
)

type Bucket struct {
    Name string     `json:"name"`
    Size int64      `json:"size"`
}

type Response struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message,omitempty"`
    Result      interface{} `json:"result,omitempty"`
}

type Page struct {
    Total           int         `json:"total"`
    Offset          int         `json:"offset"`
    Limit           int         `json:"limit"`
    Pattern         string      `json:"pattern,omitempty"`
    Buckets         *[]Bucket   `json:"buckets,omitempty"`
}

type Controller struct {
    config *config.Config
}

func sendError(context *gin.Context, err error) {
    if err == nil {
        err = errors.New("undefined")
    }
    log.Printf("%s\n", err)
    response:= Response{
        Error: true,
        Message: fmt.Sprintf("%s", err),
        Result: nil,
    }
    context.JSON(http.StatusBadRequest, response)
}

func sendMessage(context *gin.Context, message string) {
    log.Printf("%s\n", message)
    responce := Response{
        Error: false,
        Message: fmt.Sprintf("%s", message),
        Result: nil,
    }
    context.JSON(http.StatusBadRequest, responce)
}

func sendResult(context *gin.Context, result interface{}) {
    responce := Response{
        Error: false,
        Message: "",
        Result: result,
    }
    context.JSON(http.StatusOK, responce)
}

func (this *Controller) PageList(context *gin.Context) {

    var page Page
    _ = context.Bind(&page)

    storeDir, _ := this.config.GetStoreDir()

    _, err := os.Stat(storeDir)
    if err != nil {
        sendError(context, err)
        return
    }

    directoryNameList, err := tools.PathWalkDir(storeDir, MaxBucketDepth)
    if err != nil {
        sendError(context, err)
        return
    }

    list := []Bucket{}
    for i := range directoryNameList {
        directoryName := directoryNameList[i]
        fi, err := os.Stat(directoryName)
        if err != nil {
            log.Printf("%s\n", err)
            continue
        }
        if fi.Mode().IsRegular() {
            continue
        }
        size, err := tools.BucketSize(directoryName)
        if err != nil {
            continue
        }
        name := strings.TrimPrefix(directoryName, storeDir)
        name = strings.TrimLeft(name, "/")

        pattern := "*" + page.Pattern + "*"
        match, _ := filepath.Match(pattern, name)
        if match {
            list = append(list, Bucket{
                    Name: name,
                    Size: size,
            })

        }

    }

    up := page.Offset + page.Limit
    if (up > len(list)) {
        up = len(list)
    }
    down := page.Offset
    if (down < 0) {
        down = 0
    }
    if (down > len(list)) {
        down = len(list)
    }

    subList := list[down:up]
    page.Buckets = &subList
    page.Total = len(list)
    sendResult(context, &page)
}


func (this *Controller) List(context *gin.Context) {

    storeDir, _ := this.config.GetStoreDir()

    _, err := os.Stat(storeDir)
    if err != nil {
        sendError(context, err)
        return
    }

    directoryNameList, err := tools.PathWalkDir(storeDir, MaxBucketDepth)
    if err != nil {
        sendError(context, err)
        return
    }

    list := []Bucket{}
    for i := range directoryNameList {
        directoryName := directoryNameList[i]
        fi, err := os.Stat(directoryName)
        if err != nil {
            log.Printf("%s\n", err)
            continue
        }
        if fi.Mode().IsRegular() {
            continue
        }
        size, err := tools.BucketSize(directoryName)
        if err != nil {
            continue
        }
        name := strings.TrimLeft(directoryName, storeDir)
        list = append(list, Bucket{
                Name: name,
                Size: size,
        })
    }
    sendResult(context, list)
}

func (this *Controller) Hello(context *gin.Context) {
    sendMessage(context, "hello")
}


func New(config *config.Config) *Controller {
    return &Controller{
        config: config,
    }
}
