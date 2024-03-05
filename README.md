Anuja Kench@ANUJA-IBM MINGW64 ~/Desktop/short-url-generator (main)
$ go build -o short-url main.go

Anuja Kench@ANUJA-IBM MINGW64 ~/Desktop/short-url-generator (main)
$ ./short-url
URL shortening service has started.



Anuja Kench@ANUJA-IBM MINGW64 ~/Desktop
$ curl -H 'Content-Type: application/json' -d '{ "url": "https://google.com"}' -X POST http://localhost:8080/shortenurl
{"message":"Status OK"}
