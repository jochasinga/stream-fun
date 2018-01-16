# stream-fun
Video streaming service

# Install

Clone this repo and install `httprouter`.

```bash

$ git clone https://github.com/jochasinga/stream-server; cd stream-server
$ go get github.com/julienschmidt/httprouter

```

Copy an mp4 file and rename it to `stream-fun/test.mp4`. Then run

```bash

$ go build; ./stream-server

```

Browse to `http://localhost:8080`.

