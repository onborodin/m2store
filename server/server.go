/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package server

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "errors"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "html/template"
    "os"
    "os/user"
    "path/filepath"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

    "github.com/jessevdk/go-assets"

    "store/server/user-model"
    "store/server/user-controller"

    "store/server/store-controller"
    "store/server/bucket-controller"

    "store/daemon"
    "store/config"
    "store/tools"
    "store/bundle"
)

type Server struct {
    Config      *config.Config
    db          *sqlx.DB
    files       map[string]*assets.File
}


func New() *Server {
    return &Server{}
}

func (this *Server) Start() {
    var err error

    this.Config = config.New()
    this.Config.Read()
    //this.Config.Write()

    optForeground := flag.Bool("foreground", false, "run in foreground")
    flag.BoolVar(optForeground, "f", false, "run in foreground")

    optPort := flag.Int("port", this.Config.Port, "listen port")
    flag.IntVar(optPort, "p", this.Config.Port, "listen port")

    optDebug := flag.Bool("debug", this.Config.Debug, "debug mode")
    flag.BoolVar(optDebug, "d", false, "debug mode")

    optDevel := flag.Bool("devel", this.Config.Devel, "devel mode")
    flag.BoolVar(optDebug, "e", false, "devel mode")

    optWrite := flag.Bool("write", false, "write config")
    flag.BoolVar(optWrite, "w", false, "write config")

    exeName := filepath.Base(os.Args[0])

    flag.Usage = func() {
        fmt.Println("")
        fmt.Printf("usage: %s command [option]\n", exeName)
        fmt.Println("")
        flag.PrintDefaults()
        fmt.Println("")
    }
    flag.Parse()

    this.Config.Port = *optPort
    this.Config.Debug = *optDebug
    this.Config.Devel = *optDevel

    if *optWrite == true {
        fmt.Printf("write configuration to %s\n", this.Config.ConfigPath)
        err := this.Config.Write()
        if err != nil {
            fmt.Printf("write configuration error: %s\n", err)
            os.Exit(1)
        }
        os.Exit(0)
    }

    /* Daemonize process */
    if !*optForeground {
        daemon.ForkProcess()
    }

    /* Lookup user system info */
    user, err := user.Lookup(this.Config.User)
    if err != nil {
        fmt.Printf("user lookup error: %s\n", err)
        os.Exit(1)
    }

    /* Make process ID directory */
    err = os.MkdirAll(filepath.Dir(this.Config.PidPath), 0750)
    if err != nil {
        log.Printf("unable create rundir: %s\n", err)
        os.Exit(1)
    }

    /* Save process ID to file */
    if err := daemon.SaveProcessID(this.Config.PidPath); err != nil {
        fmt.Printf("%s; exit\n", err)
        os.Exit(1)
    }
    defer os.Remove(this.Config.PidPath)

    uid, err := strconv.Atoi(user.Uid)

    /* Make log directory */
    err = os.MkdirAll(filepath.Dir(this.Config.MessageLogPath), 0750)
    if err != nil {
        log.Printf("unable create message log dir: %s\n", err)
        os.Exit(1)
    }
    err = os.Chown(filepath.Dir(this.Config.MessageLogPath), uid, os.Getgid())
    if err != nil {
        log.Printf("unable chown log dir: %s\n", err)
        os.Exit(1)
    }

    /* Make store directory */
    err = os.MkdirAll(this.Config.StoreDir, 0750)
    if err != nil {
        log.Printf("unable create store dir: %s\n", err)
        os.Exit(1)
    }
    err = os.Chown(this.Config.StoreDir, uid, os.Getgid())
    if err != nil {
        log.Printf("unable chown store dir: %s\n", err)
        os.Exit(1)
    }

    if _, err := os.Stat(this.Config.StoreDir); err != nil {
        log.Printf("store dir not exists: %s\n", err)
        os.Exit(1)
    }

    /* Change effective user ID */
    if uid != 0 {
        err = syscall.Setuid(uid)
        if err != nil {
            log.Printf("set process user id error: %s\n", err)
            os.Exit(1)
        }
        if syscall.Getuid() != uid {
            log.Printf("set process user id error: %s\n", err)
            os.Exit(1)
        }
    }

    /* Redirect log to message file */
    file, err := daemon.RedirectLog(this.Config.MessageLogPath, *optDebug)
    if err != nil {
        fmt.Printf("%s; exit\n", err)
        os.Exit(1)
    }
    defer file.Close()

    /* Redirect standard IO */
    if !*optForeground {
        if _, err := daemon.RedirectIO(); err != nil {
            log.Printf("%s; exit\n", err)
            os.Exit(1)
        }
    }

    log.Printf("%s start :%d\n", exeName, this.Config.Port)

    err = this.Run()

    if err != nil {
        log.Printf("%s; exit\n", err)
        os.Exit(1)
    }
}

func (this *Server) Run() error {

    /* embedded assets init */
    this.files = bundle.Assets.Files

    var err error

    dbUrl := fmt.Sprintf("%s", this.Config.PasswordPath)

    this.db, err = sqlx.Open("sqlite3", dbUrl)
    if err != nil {
        return err
    }

    /* Check DB connection */
    err = this.db.Ping()
    if err != nil {
        return err
    }

    //fmt.Println("debug mode:", this.Config.Debug)
    if this.Config.Debug{
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }
    gin.DisableConsoleColor()

    accessLogFile, err := os.OpenFile(this.Config.AccessLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
    if err != nil {
      return err
    }
    gin.DefaultWriter = io.MultiWriter(accessLogFile, os.Stdout)
    //gin.DefaultWriter = ioutil.Discard

    router := gin.New()

    /* Dump req/res */
    if this.Config.Debug{
        router.Use(RequestLogMiddleware())
        router.Use(ResponseLogMiddleware())
    }

    //router.Use(gin.Logger())
    router.Use(gin.LoggerWithFormatter(logFormatter()))
    router.Use(gin.Recovery())

    /* Read templates */
    if this.Config.Devel {
        /* Filesystem variant */
        router.LoadHTMLGlob(filepath.Join(this.Config.LibDir, "public/index.html"))
    } else {
        /* Embedded variant */
        data, err := ioutil.ReadAll(this.files["/public/index.html"])
        if err != nil {
            return err
        }
        tmpl, err := template.New("index.html").Parse(string(data))
        router.SetHTMLTemplate(tmpl)
    }

    store := cookie.NewStore([]byte("ds79asd9a7d9sa7d9sa87d"))
    store.Options(sessions.Options{
        MaxAge: 3600 * 4,
        Path:   "/",
    })
    router.Use(sessions.Sessions("session", store))


    router.GET("/", this.Index)

    userController := userController.New(this.Config, this.db)

    router.POST("/api/v2/user/login", userController.Login)
    router.POST("/api/v2/user/logout", userController.Logout)

    humanGroup := router.Group("/api/v2")
    humanGroup.Use(this.sessionAuthMiddleware)

    humanGroup.POST("/user/list", userController.List)
    humanGroup.POST("/user/create", userController.Create)
    humanGroup.POST("/user/delete", userController.Delete)
    humanGroup.POST("/user/update", userController.Update)

    botGroup := router.Group("/api/v1")
    botGroup.Use(this.uniAuthMiddleware)

    bucketController := bucketController.New(this.Config)
    botGroup.GET("/bucket/list", bucketController.List)
    botGroup.POST("/bucket/list", bucketController.List)
    botGroup.POST("/bucket/pagelist", bucketController.PageList)

    botGroup.GET("/bucket/hello", bucketController.Hello)

    storeController := storeController.New(this.Config)
    botGroup.POST("/file/list", storeController.List)
    botGroup.POST("/file/pagelist", storeController.PageList)
    botGroup.POST("/file/put", storeController.Put)
    botGroup.POST("/file/get", storeController.Get)
    botGroup.POST("/file/drop", storeController.Drop)
    botGroup.GET("/file/hello", storeController.Hello)
    botGroup.GET("/file/down/*path", storeController.Down)

    /* No route handler */
    router.NoRoute(this.NoRoute)

    daemon.SetSignalHandler()

    return router.RunTLS(":" + fmt.Sprintf("%d", this.Config.Port), this.Config.CertPath, this.Config.KeyPath)
}

func (this *Server) Index(context *gin.Context) {
    context.HTML(http.StatusOK, "index.html", nil)
}

func (this *Server) NoRoute(context *gin.Context) {

    requestPath := context.Request.URL.Path

    if this.Config.Devel {
        /* Filesystem assets */
        publicDir := filepath.Join(this.Config.LibDir, "public")
        filePath := filepath.Clean(filepath.Join(publicDir, requestPath))
        if !strings.HasPrefix(filePath, publicDir) {
            err := errors.New(fmt.Sprintf("wrong file patch %s\n", filePath))
            log.Println(err)
            context.HTML(http.StatusOK, "index.html", nil)
            return
        }
        /* for frontend handle: If file not found send index.html */
        if !tools.FileExists(filePath) {
            err := errors.New(fmt.Sprintf("file path not found %s\n", filePath))
            log.Println(err)
            context.HTML(http.StatusOK, "index.html", nil)
            return
        }
        context.File(filePath)
    } else {
        /* Embedded assets variant */
        file := this.files[filepath.Join("/public", requestPath)] //io.Reader
        if file == nil {
            err := errors.New(fmt.Sprintf("file path not found %s, send index", requestPath))
            log.Println(err)
            context.HTML(http.StatusOK, "index.html", nil)
            return
        }
        http.ServeContent(context.Writer, context.Request, requestPath, file.ModTime(), file)
    }
}

func logFormatter() func(param gin.LogFormatterParams) string {
    return func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s %s %s %s %s %d %d %s\n",
            param.TimeStamp.Format(time.RFC3339),
            param.ClientIP,
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.BodySize,
            param.Latency,
        )
    }
}

func (this *Server) sessionAuthMiddleware(context *gin.Context) {
    session := sessions.Default(context)

    username := session.Get("username")
    if username == nil || len(username.(string)) == 0 {

        result := Result{
            Error: true,
            Message: "wrong session autentification",
            Result: "",
        }
        context.JSON(http.StatusUnauthorized, result)
        context.Abort()
        return
    }
    context.Next()
}


type Result struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message"`
    Result      interface{} `json:"result,omitempty"`
}

func (this *Server) basicAuthMiddleware(context *gin.Context) {
    authHeader := context.Request.Header.Get("Authorization")
    userName, password, err := parseAuthBasicHeader(authHeader)
    if err != nil {
        result := Result{
            Error: true,
            Message: fmt.Sprintf("parse auth header error: %s", err),
            Result: "",
        }
        context.JSON(http.StatusUnauthorized, result)
        context.Abort()
        return
    }

    if !this.authenticateUser(userName, password) {
        result := Result{
            Error: true,
            Message: fmt.Sprintf("parse auth header error: %s", err),
            Result: "",
        }
        context.JSON(http.StatusUnauthorized, result)
        context.Abort()
        return
    }
    context.Next()
}


func (this *Server) uniAuthMiddleware(context *gin.Context) {

    session := sessions.Default(context)
    username := session.Get("username")
    if username != nil && len(username.(string)) > 0 {
        context.Next()
        return
    }

    authHeader := context.Request.Header.Get("Authorization")

    userName, password, err := parseAuthBasicHeader(authHeader)
    if err != nil {
        result := Result{
            Error: true,
            Message: fmt.Sprintf("parse auth header error: %s", err),
            Result: "",
        }
        context.JSON(http.StatusUnauthorized, result)
        context.Abort()
        return
    }

    if !this.authenticateUser(userName, password) {
        result := Result{
            Error: true,
            Message: fmt.Sprintf("wrong basic authorization"),
            Result: "",
        }
        context.JSON(http.StatusUnauthorized, result)
        context.Abort()
        return
    }
    context.Next()
}


func (this *Server) authenticateUser(username string, password string) bool {
    user := userModel.New(this.db)
    theUser := userModel.User{
        Username: username,
        Password: password,
    }
    err := user.Check(&theUser)
    if err != nil {
        log.Printf("autentification error: %s", err)
        return false
    }
    return true
}

func parseAuthBasicHeader(header string) (string, string, error) {
    auth := strings.SplitN(header, " ", 2)
    authType := strings.TrimSpace(auth[0])
    if authType != "Basic" {
        return "", "", errors.New("authentification type is different from basic")
    }
    authPair := strings.TrimSpace(auth[1])

    pairEncoded, err := base64.StdEncoding.DecodeString(authPair)
    if err != nil {
        return "", "", err
    }
    pair := strings.SplitN(string(pairEncoded), ":", 2)
    if len(pair) < 2 {
        return "", "", errors.New("wrong authentification pair")
    }

    login := strings.TrimSpace(pair[0])
    pass := strings.TrimSpace(pair[1])

    if len(login) == 0 {
        return "", "", errors.New("autentification username is null")
    }
    if len(pass) == 0 {
        return "", "", errors.New("autentification password is null")
    }
    return login, pass, nil
}

func RequestLogMiddleware() gin.HandlerFunc {
    return func(context *gin.Context) {

        var requestBody []byte
        if context.Request.Body != nil {
            requestBody, _ = ioutil.ReadAll(context.Request.Body)
        }

        contentType := context.GetHeader("Content-Type")
        contentType = strings.ToLower(contentType)

        buffer := bytes.NewBuffer(nil)
        json.Indent(buffer, requestBody, "", "    ")

        if strings.Contains(contentType, "application/json") {
            log.Print("request:\n", buffer.String())
        }

        context.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
        context.Next()
    }
}

func ResponseLogMiddleware() gin.HandlerFunc {

    return func(context *gin.Context) {
        contentType := context.GetHeader("Content-Type")
        contentType = strings.ToLower(contentType)

        writer := &LogWriter{
            body: bytes.NewBuffer(nil),
            ResponseWriter: context.Writer,
        }
        context.Writer = writer

        context.Next()

        buffer := bytes.NewBuffer(nil)
        json.Indent(buffer, writer.body.Bytes(), "", "    ")

        if strings.Contains(contentType, "application/json") {
            log.Print("response:\n", buffer.String())
        }
    }
}

type LogWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (this LogWriter) Write(data []byte) (int, error) {
    this.body.Write(data)
    return this.ResponseWriter.Write(data)
}

func (this LogWriter) WriteString(data string) (int, error) {
    this.body.WriteString(data)
    return this.ResponseWriter.WriteString(data)
}
