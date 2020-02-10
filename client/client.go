/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package client

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

type Client struct {
}

const (
   listURI  string = "api/v1/file/list"
   putURI   string = "api/v1/file/put"
   getURI   string = "api/v1/file/get"
   dropURI  string = "api/v1/file/drop"
)

type ListForm struct {
    Bucket      string  `json:"bucket"`
    Pattern     string  `json:"pattern"`
}

type File struct {
    Name string     `json:"name"`
    Size int64      `json:"size"`
    ModTime string  `json:"modtime"`
}

type ListResult struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message"`
    Files       []File      `json:"result"`
}

func (this *Client) List(hostname, username, password, bucket, pattern string) (string, error) {

    var err error
    url := fmt.Sprintf("https://%s:%s@%s/%s", username, password, hostname, listURI)

    form := ListForm{
        Bucket: bucket,
        Pattern: pattern,
    }

    data, _ := json.Marshal(form)
    reader := bytes.NewReader([]byte(data))

    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: transCfg}

    resp, err := client.Post(url, "application/json", reader)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

//type PutForm struct {
//    Bucket      string  `json:"bucket"  form:"bucket"`
//    Filename    string  `json:"pattern" form:"filename"`
//}

func (this *Client) Put(hostname, username, password, bucket, filename string) (string, error) {

    var err error
    url := fmt.Sprintf("https://%s:%s@%s/%s", username, password, hostname, putURI)

    pipeOut, pipeIn := io.Pipe()
    writer := multipart.NewWriter(pipeIn)

    go func() {
        defer pipeIn.Close()
        defer writer.Close()

        _ = writer.WriteField("filename", filepath.Base(filename))
        _ = writer.WriteField("bucket", bucket)

        part, err := writer.CreateFormFile("file", filepath.Base(filename))
        if err != nil {
            return
        }
        file, err := os.Open(filename)
        if err != nil {
            return
        }
        defer file.Close()
        if _, err = io.Copy(part, file); err != nil {
            return
        }
    }()


    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
    }
    client := &http.Client{ Transport: transCfg }
    resp, err := client.Post(url, writer.FormDataContentType(), pipeOut)
    if err != nil {
        return "", err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(body), err
}

type GetForm struct {
    Bucket      string  `json:"bucket"  form:"bucket"`
    Filename    string  `json:"filename" form:"filename"`
}

func (this *Client) Get(hostname, username, password, bucket, filename string) (string, error) {

   var err error
    url := fmt.Sprintf("https://%s:%s@%s/%s", username, password, hostname,  getURI)

    form := GetForm{
        Bucket: bucket,
        Filename: filename,
    }

    data, _ := json.Marshal(form)
    reader := bytes.NewReader([]byte(data))


    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
    }
    client := &http.Client{ Transport: transCfg }

    resp, err := client.Post(url, "application/json", reader)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return resp.Status, nil
    }

    out, err := os.Create(filepath.Base(filename))
    if err != nil {
        return "", err
    }
    defer out.Close()
    //_, err = io.Copy(out, resp.Body)

    buf := make([]byte,  128 * 1024)
    _, _ = io.CopyBuffer(out, resp.Body, buf)
    return resp.Status, nil
}

type DropForm struct {
    Bucket      string  `json:"bucket"   form:"bucket"`
    Filename    string  `json:"filename" form:"filename"`
}

func (this *Client) Drop(hostname, username, password, bucket, filename string) (string, error) {

    var err error
    url := fmt.Sprintf("https://%s:%s@%s/%s", username, password, hostname,  dropURI)

    form := DropForm{
        Bucket: bucket,
        Filename: filename,
    }

    data, _ := json.Marshal(form)
    reader := bytes.NewReader([]byte(data))


    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
    }
    client := &http.Client{ Transport: transCfg }

    resp, err := client.Post(url, "application/json", reader)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    return string(body), nil
}

func printHeader(header http.Header) {
    for key, val := range header {
        fmt.Printf("%s: %s\n", key, strings.Join(val, " "))
    }
}

func New() *Client {
    return &Client{
    }
}
