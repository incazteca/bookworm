# bookworm

Upload a page of your favorite book and get a response with the text, word
count of page, and a word count of each word.

You can also specify a case insesitive filter that will remove words that have
the string from the response. For example: Having the word "blue" as a filter
will remove "bluetooth", "blueberry", and "madeupblueword".

# Installation
Bookworm is go project, Please refer to the [golang installation guide](https://golang.org/doc/install)
for getting started with Go.

`go get` will not work due to the repo being in bitbucket. Instead execute the
following commands:

```
mkdir -p $GOPATH/src/bitbucket.org/incazteca
cd $GOPATH/src/bitbucket.org/incazteca
git clone https://bitbucket.org/incazteca/bookworm
```

Once bookworm is cloned you can run `go install` so that an executable command
named `bookworm` will be in your `$GOPATH/bin` directory. Execute `bookworm` to
start the server which will be available on port 8080.

# Usage

The bookworm server is available on port 8080.

At the moment there is only one endpoint which is for uploading files

POST to `example.com:8080/file/upload` with the following parameters:

    - file field, named "file". Use to specify file to upload
    - text field, named "filter". Use to specify a filter if you want.

An example request is done here with curl:

Without filter
`curl -F "file=@/path/to/file" -w "%{http_code}" localhost:8080/file/upload &`

with filter
`curl -F "file=@/path/to/file" -F "filter=blue" -w "%{http_code}" localhost:8080/file/upload &`

Response Content-Type will be application/json

If you find yourself wanting to upload a small file while a larger one is uploading
just send it over. It will process without waiting for the larger file to finish.

Also attempts to use GET on the `/file/upload` endpoint will be met with 404.

A 413 Error will also be received if an attempt is made of uploading a file larger
than 10MB

# Limitations

Requests are currently not saved so you will have to resubmit a request if you want
the results from an older request.

Also requests are global. There are no endpoints yet for user specific requests 
until older requests can be saved.
