
package fileController

import (
    "errors"
    "fmt"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"

    "store/config"
    "store/tools"
)

const (
    MaxBucketDepth int = 64
)

type File struct {
    Name string     `json:"name"`
    Size int64      `json:"size"`
    ModTime string  `json:"modtime"`
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
    Bucket          string      `json:"bucket"`
    Pattern         string      `json:"pattern,omitempty"`
    Files           *[]File     `json:"files,omitempty"`
}

type Controller struct {
    config *config.Config
}

func sendError(context *gin.Context, err error) {
    if err == nil {
        err = errors.New("undefined")
    }
    log.Printf("%s\n", err)
    response := Response{
        Error: true,
        Message: fmt.Sprintf("%s", err),
        Result: nil,
    }
    context.JSON(http.StatusBadRequest, response)
}

func sendMessage(context *gin.Context, message string) {
    log.Printf("%s\n", message)
    response := Response{
        Error: false,
        Message: fmt.Sprintf("%s", message),
        Result: nil,
    }
    context.JSON(http.StatusBadRequest, response)
}

func sendResult(context *gin.Context, result interface{}) {
    response := Response{
        Error: false,
        Message: "",
        Result: result,
    }
    context.JSON(http.StatusOK, response)
}

func (this *Controller) ValidateFilePath(bucketName, fileName string) (string, error) {
    storeDir, _ := this.config.GetStoreDir()

    directoryPath := filepath.Clean(filepath.Join(storeDir, bucketName))
    if !strings.HasPrefix(directoryPath, storeDir) {
        return "", errors.New("wrong bucket name")
    }

    filePath := filepath.Clean(filepath.Join(directoryPath, fileName))
    if !strings.HasPrefix(filePath, storeDir) {
        return "", errors.New("wrong backet or file name")
    }
    return filePath, nil
}

func (this *Controller) ValidateBucketPath(bucketName string) (string, error) {

    storeDir, _ := this.config.GetStoreDir()

    directoryPath := filepath.Clean(filepath.Join(storeDir, bucketName))
    if !strings.HasPrefix(directoryPath, storeDir) {
        return "", errors.New("wrong bucket name")
    }
    return directoryPath, nil
}

func (this *Controller) PageList(context *gin.Context) {

    /* Bind form */
    var page Page
    if err := context.Bind(&page); err != nil {
        sendError(context, err)
        return
    }

    /* Validate bucket */
    directoryPath, err := this.ValidateFilePath(page.Bucket, "")
    if err != nil {
        sendError(context, err)
        return
    }

    /* Check file path for existing */
    _, err = os.Stat(directoryPath)
    if err != nil {
        sendError(context, err)
        return
    }

    /* Validate pattern */
    pattern := "*"
    if page.Pattern != "" {
            if tools.PathLength(page.Pattern) > 1 {
                sendError(context, err)
                return
            }
            pattern = page.Pattern
    }

    filePath, err := this.ValidateFilePath(page.Bucket, pattern)
    if err != nil {
        sendError(context, err)
        return
    }

    /* List directory by pattern */
    fileNameList, err := filepath.Glob(filePath)
    if err != nil {
        sendError(context, err)
        return
    }

    list := []File{}
    for _, fileName := range fileNameList {
        fi, err := os.Stat(fileName)
        if err != nil {
            log.Printf("%s\n", err)
            continue
        }
        if !fi.Mode().IsRegular() {
            continue
        }
        list = append(list, File{
                Name: filepath.Base(fileName),
                Size: fi.Size(),
                ModTime: fi.ModTime().Format(time.RFC3339),
            })
    }

    /* Send result */
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
    page.Files = &subList
    page.Total = len(list)
    sendResult(context, &page)
}


type listForm struct {
    Bucket  string  `form:"bucket"  json:"bucket"`
    Pattern string  `form:"pattern" json:"pattern"`
}

func (this *Controller) List(context *gin.Context) {

    /* Bind form */
    var form listForm
    if err := context.Bind(&form); err != nil {
        sendError(context, err)
        return
    }

    /* Validate bucket */
    directoryPath, err := this.ValidateFilePath(form.Bucket, "")
    if err != nil {
        sendError(context, err)
        return
    }

    /* Check file path for existing */
    _, err = os.Stat(directoryPath)
    if err != nil {
        sendError(context, err)
        return
    }

    /* Validate pattern */
    pattern := "*"
    if form.Pattern != "" {
            if tools.PathLength(form.Pattern) > 1 {
                sendError(context, err)
                return
            }
            pattern = form.Pattern
    }

    filePath, err := this.ValidateFilePath(form.Bucket, pattern)
    if err != nil {
        sendError(context, err)
        return
    }

    /* List directory by pattern */
    fileNameList, err := filepath.Glob(filePath)
    if err != nil {
        sendError(context, err)
        return
    }

    list := []File{}
    for _, fileName := range fileNameList {
        fi, err := os.Stat(fileName)
        if err != nil {
            log.Printf("%s\n", err)
            continue
        }
        if !fi.Mode().IsRegular() {
            continue
        }
        list = append(list, File{
                Name: filepath.Base(fileName),
                Size: fi.Size(),
                ModTime: fi.ModTime().Format(time.RFC3339),
            })
    }
    /* Send result */
    sendResult(context, list)
}

type putForm struct {
    FileName    string          `form:"filename" binding:"required"`
    BucketName  string          `form:"bucket"`
    File *multipart.FileHeader  `form:"file"     binding:"required"`
}

func (this *Controller) Put(context *gin.Context) {

    /* Bind form */
    form := putForm{}
    if err := context.ShouldBind(&form); err != nil {
        sendError(context, err)
        return
    }

    /* Validate bucket name */
    directoryPath, err := this.ValidateBucketPath(form.BucketName)
    if err != nil {
        sendError(context, err)
        return
    }

    /* Validate file name */
    filePath, err := this.ValidateFilePath(form.BucketName, form.FileName)
    if err != nil {
        sendError(context, err)
        return
    }

    /* Store file */
    if err := os.MkdirAll(directoryPath, os.ModeDir | 0750); err != nil {
        sendError(context, err)
        return
    }

    file := form.File
    if err := context.SaveUploadedFile(file, filePath); err != nil {
        sendError(context, err)
        return
    }

    /* Check uploaded file */
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        sendError(context, err)
        return
    }

    if !fileInfo.Mode().IsRegular() {
        sendError(context, err)
        return
    }

    /* Send file info */
    var list []File
    list = append(list, File{
                Name: filepath.Base(filePath),
                Size: fileInfo.Size(),
                ModTime: fileInfo.ModTime().Format(time.RFC3339),
            })

    sendResult(context, list)
}

type getForm struct {
    FileName    string  `form:"filename" json:"filename" binding:"required" `
    BucketName  string  `form:"bucket"   json:"bucket"`
}

func (this *Controller) Get(context *gin.Context) {

    form := getForm{}
    if err := context.ShouldBind(&form); err != nil {
        log.Printf("%s\n", err)
        context.Status(http.StatusNotFound)
        return
    }

    /* Validate file name */
    filePath, err := this.ValidateFilePath(form.BucketName, form.FileName)
    if err != nil {
        log.Println(err)
        context.Status(http.StatusNotFound)
        return
    }

    /* Check real file */
    if !tools.FileExists(filePath) {
        err := errors.New(fmt.Sprintf("file path not found %s\n", filePath))
        log.Println(err)
        context.Status(http.StatusNotFound)
        return
    }
    context.FileAttachment(filePath, filepath.Base(filePath))
}

func (this *Controller) Down(context *gin.Context) {
    paramPath := context.Param("path")

    /* Validate file name */
    filePath, err := this.ValidateFilePath("", paramPath)
    if err != nil {
        log.Println(err)
        context.Status(http.StatusNotFound)
        return
    }

    if !tools.FileExists(filePath) {
        err := errors.New(fmt.Sprintf("file path not found %s\n", filePath))
        log.Println(err)
        context.Status(http.StatusNotFound)
        return
    }
    context.FileAttachment(filePath, filepath.Base(filePath))
}

type deleteForm struct {
    FileName    string  `form:"filename" json:"filename" binding:"required" `
    BucketName  string  `form:"bucket"   json:"bucket"`
}

func (this *Controller) Delete(context *gin.Context) {

    form := deleteForm{}
    if err := context.ShouldBind(&form); err != nil {
        sendError(context, err)
        return
    }

    storePath, _ := this.config.GetStoreDir()
    reqPath := filepath.Join(form.BucketName, form.FileName)
    fullPath := filepath.Clean(filepath.Join(storePath, reqPath))

    if !strings.HasPrefix(fullPath, storePath) {
        err := errors.New(fmt.Sprintf("wrong file name %s", reqPath))
        sendError(context, err)
        return
    }
    if !tools.FileExists(fullPath) {
        err := errors.New(fmt.Sprintf("wrong file name %s", reqPath))
        sendError(context, err)
        return
    }

    /* Real revive file */
    err := os.Remove(fullPath)
    if err != nil {
        sendError(context, err)
        return
    }

    /* Validate operation */
    if tools.FileExists(fullPath) {
        err := errors.New(fmt.Sprintf("wrong file name %s", reqPath))
        sendError(context, err)
        return
    }

    /* Clean directory if empty */
    _ = syscall.Rmdir(storePath)

    sendResult(context, []File{})
}


func (this *Controller) Hello(context *gin.Context) {
    sendMessage(context, "hello")
}

func New(config *config.Config) *Controller {
    return &Controller{
        config: config,
    }
}
