# Clipboard Share
Share your clipboard text to local network.

# Usage
## S1 clone this repo
```shell script
git clone https://github.com/tignioj/go-clipboard-share.git
```
## S2 build source code
```shell script
cd go-clipboard-share
go build
```
After executing the command above, you will see `go-clipboard-share.exe` in this directory.
## S3 Run it.
- Double click go-clipboard-share.exe
- Open http://localhost:9999

# Argument
## list
- `-p`: Listen port. default 9999
- `-h`: Show help
## example
### Set listen port to 8080
```shell script
./go-clipboard-share.exe -p 8080
```

# API
## `/get` Get server clipboard text raw string
- Request method: GET
### example: 
- http://localhost:9999/get

## `/set` Set server clipboard text
- Request method: POST
- Content-Type: application/x-www-form-urlencoded 
- Date-format: key=value
Ajax example: 
```javascript
    var xhr = new XMLHttpRequest()
    xhr.open("POST", "/set")
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    var data = "contentToSet=" + yourTextToSet
    xhr.send(data)
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            showInfo(xhr.responseText)
        }
    }
```


# NOTICE
 - These two files `index.html` and `help.json`  are not required, it just a template fo test. I have already set it to go source code as variable.
 
 # THANKS
 - `Clipboard`: https://github.com/atotto/clipboard