# stream-fun
Video streaming service

# Build

Clone this repo and install `httprouter` and `redigo/redis`.

```bash

$ git clone https://github.com/jochasinga/stream-server; cd stream-server
$ go get github.com/julienschmidt/httprouter
$ go get github.com/garyburd/redigo/redis

```

Copy one or more mp4 files and rename them to `stream-fun/assets/my-hyphenated-title.mp4`. The filename must be hyphenated in order to be properly formatted into a readable title (in the above 
case it would become "My Hyphenated Title").

Then build the project with

```bash

$ go build

```

# Run Redis

Redis server needs to listen on port 6379 (default). Use your OS's package manager
to install Redis and run it like so:

```bash

# mac OS
$ brew update && brew install redis
$ redis-server

```

Finally, run the freshly built executable

````bash

$ `./stream-fun`

```

Browse to `http://localhost:8080` to see the movie playing.

