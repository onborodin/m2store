
package statusController

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "syscall"

    "github.com/gin-gonic/gin"

    "store/config"
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


func (this *Controller) Hello(context *gin.Context) {
    sendMessage(context, "hello")
}

type Disk struct {
    Free    uint64      `json:"free"`
}


func (this *Controller) Disk(context *gin.Context) {
    var disk Disk

    var stat syscall.Statfs_t
    syscall.Statfs(this.config.StoreDir, &stat)
    disk.Free = uint64(stat.Bavail * int64(stat.Bsize))

    sendResult(context, disk)
}


func New(config *config.Config) *Controller {
    return &Controller{
        config: config,
    }
}
