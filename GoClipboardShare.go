package main

import (
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/tignioj/go-get-argsmap-from-commandline"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
)

const (
	cfUnicodetext = 13
)
//v1
func getText() (string, error) {
	content, err := clipboard.ReadAll()
	if err != nil {
		// BUGFIX: Check Clipboard is Empty
		u := syscall.MustLoadDLL("user32")
		isClipboardFormatAvailable := u.MustFindProc("IsClipboardFormatAvailable")
		if formatAvailable, _, _ := isClipboardFormatAvailable.Call(cfUnicodetext); formatAvailable == 0 {
			return "", nil
		}
		return "", errors.New("Get clipboard to server failure:" + err.Error())
	}
	return content, nil
}
func setText(content string) (string, error) {
	err1 := clipboard.WriteAll(content)
	if err1 != nil {
		return "", errors.New("Set clipboard to server failure:" + err1.Error())
	} else {
		return "Copying to server success!", nil
	}
}

func main() {
	var port = "9999"
	var json = `
{
  "-p": {
    "value": "9999",
    "usage": "Listen Port",
    "pattern": "\\d+",
    "err": "Invalid port"
  },
  "-h": {
    "usage": "show help",
    "must_have_value": false
  }
}
`
	obj, err := argsmap.NewCommandLineObjByJSON(json, os.Args)
	if err != nil {
		log.Fatal("clipboard: argument error:", err.Error())
	}
	v, err1 := obj.GetArg("-p")
	if err1 == nil {
		port = v
	}
	_, err2 := obj.GetArg("-h")
	if err2 == nil {
		obj.ShowHelp()
		return
	}

	http.HandleFunc("/", indexView)
	http.HandleFunc("/set", setClipboardView)
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		b, err := getText()
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(b))
	})
	log.Println("Starting share clipboard server: http://localhost:" + port)
	addr := ":" + port
	log.Fatal(http.ListenAndServe(addr, nil))
}

func setClipboardView(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		strToSet := r.FormValue("contentToSet")

		respText, err1 := setText(strToSet)
		if err1 != nil {
			w.Write([]byte("Set Error:" + err1.Error()))
		}
		w.Write([]byte(respText))
	} else {
		w.WriteHeader(403)
		w.Write([]byte("403:Not allow!"))
	}
}

func indexView(w http.ResponseWriter, r *http.Request) {
	contentToGet, err := getText()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	contentToGet = html.EscapeString(contentToGet)
	var index = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta charset="UTF-8">
    <title>Title</title>
    <style>
        .root {
            padding: 10px;
        }

        .text, .btn {
            width: 100%;
            margin: 3px 0 3px 0;
        }
        .text {
            border: 0;
        }

        .btn {
            font-size: 1.5em;
            padding: 10px 0 10px 0;
        }
    </style>
</head>
<body>
<div class="root">
    <h1 align="center">Clipboard Share</h1>
    <div>
        <span>INFO: </span><span id="info"></span>
    </div>
    <form onsubmit="return false">
        <textarea class="text" readonly rows="10" cols="100" name="contentToGet" id="contentToGet">` +
		contentToGet +
		`</textarea>
        <br/>
        <button class="btn" onclick="getClipBoardFromTextArea()">Copy</button>
        <br/>
        <textarea class="text" rows="10" cols="100" placeholder="Paste your content here" name="contentToSet"
                  id="contentToSet"></textarea>
        <br/>
        <button class="btn" onclick="setClipboardToServer()">Set Clipboard</button>
    </form>

</div>
<script>
    function setClipboardToServer() {
        var textDom = document.getElementById("contentToSet")
        showInfo(textDom.value)
        var xhr = new XMLHttpRequest()
        xhr.open("POST", "/set")
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

        val = encodeURIComponent(textDom.value)
        var data = "contentToSet=" + val

        xhr.send(data)
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status === 200) {
                showInfo(xhr.responseText)
            }
        }
    }

    function fallbackCopyTextToClipboard(text) {
        var textArea = document.createElement("textarea");
        textArea.value = text;

        // Avoid scrolling to bottom
        textArea.style.top = "0";
        textArea.style.left = "0";
        textArea.style.position = "fixed";

        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            var successful = document.execCommand('copy');
            var msg = successful ? 'successful' : 'unsuccessful';
            console.log('Fallback: Copying text command was ' + msg);
            showInfo('Fallback: Copying text command was ' + msg);
        } catch (err) {
            console.error('Fallback: Oops, unable to copy', err);
            showInfo('Fallback: Oops, unable to copy', err);
        }
        document.body.removeChild(textArea);
    }

    function copyTextToClipboard(text) {
        if (!navigator.clipboard) {
            fallbackCopyTextToClipboard(text);
            return;
        }
        navigator.clipboard.writeText(text).then(function () {
            // console.log('Async: Copying to clipboard was successful!');
            showInfo('Copying to clipboard was successful!')
        }, function (err) {
            // console.error('Async: Could not copy text: ', err);
            showInfo('Copying to clipboard was successful!' + err)
        });
    }
    function getClipBoardFromTextArea() {
            var copyTextarea = document.querySelector('#contentToGet');
            copyTextToClipboard(copyTextarea.value)
    }

    function showInfo(str) {
        var ele = document.getElementById("info")
        console.log(ele)
        ele.innerText = new Date().toLocaleTimeString() + ":" + str
    }
</script>
</body>
</html>
`
	//if 1 < 0 {
	//	fmt.Println(index)
	//}
	//form := readFile("index.html")
	//w.Write(form)
	w.Write([]byte(index))
}

func readFile(fileName string) []byte {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("read:", err.Error())
	}
	return b
}
