# m2store

It is a simple file store with simple API


### API

| URL                 |Method and arguments              | Result              |
|---------------------|----------------------------------|---------------------|
| /api/v1/file/list   | POST (bucket*, filename, pattern*) | application/json    |
| /api/v1/file/put    | POST multipart/form              | application/json    |
| /api/v1/file/get    | POST (bucket*, filename)         | octet/stream or 404 |
| /api/v1/file/drop   | POST (bucket*, filename)         | application/json    |
| /api/v1/file/down   | GET /path                        | octet/stream or 404 |
| /api/v1/bucket/list | GET                              | application/json |


** - optional arguments

### Result

    type Result struct {
        Error       bool        `json:"error"`
        Message     string      `json:"message"`
        Result      interface{} `json:"result"`
    }


Result can be the File, []File or empty File{}

    type File struct {
        Name string     `json:"name"`
        Size int64      `json:"size"`
        ModTime string  `json:"modtime"`
    }

    type Bucket struct {
        Name string     `json:"name"`
        Size int64      `json:"size"`
    }

### Examples

#### Bucket list

    curl -v http://user:12345@127.0.0.1:8080/api/v1/bucket/list

Response:

    {"error":false,"message":"success","result":[
        {"name":"","size":0},
        {"name":"foobar","size":11534336}
    ]}


#### Put file

    curl -v -F filename=data44.bin -F bucket=foobar -F file=@blob.bin
        http://user:1234@127.0.0.1:8080/api/v1/file/put

Response:

    {
        "error":false,
        "message":"success",
        "result":[
            { "name":"data44.bin","size":1048576,"modtime":"2019-12-13T12:16:10+02:00" }
        ]
    }

#### File list

    curl -v -X POST -H "Content-Type: application/json"
        -d '{ "bucket": "foobar", "pattern": "data*4*" }'
         http://user:12345@127.0.0.1:8080/api/v1/file/list

Response:

    {
        "error":false,
        "message":"success",
        "result":[
            {"name":"data14.bin","size":1048576,"modtime":"2019-12-12T13:56:23+02:00"},
            {"name":"data44.bin","size":1048576,"modtime":"2019-12-13T12:16:10+02:00"}
        ]
    }

Pattern and bucket is optional

#### Get file

    curl -X POST -H "Content-Type: application/json"
        -d '{ "bucket": "foobar", "filename":"data.bin" }'
        http://user:1234@127.0.0.1:8080//api/v1/file/get

Or, GET form with full path

    curl -X GET http://user:1234@127.0.0.1:8080/api/v1/file/down/foobar/data.bin

#### Drop file

    curl -X POST -H "Content-Type: application/json"
        -d '{ "bucket": "foobar", "filename":"data.bin" }'
        http://user:1234@127.0.0.1:8080/api/v1/file/drop


----
