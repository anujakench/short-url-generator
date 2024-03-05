Build and Deploy the application -
On a Windows 10 OS, install Git bash, Go1.22.0.
Open git bash terminal.
Cd to the repository directory.
Clone the github repository https://github.com/anujakench/short-url-generator.git on the machine(like i home directory).
Run  go build -o short-url main.go
Run ./short-url
Open another terminal.
Run curl -H 'Content-Type: application/json' -d '{ "url": "https://gmail.com"}' -X POST http://localhost:8080/shortenurl to create short URL
Run curl -H 'Content-Type: application/json' -d '{ "url": "http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t"}' -X POST http://localhost:8080/deleteurl to delete a short URL. Note: replace the short URL appropriately. It will be printed on the console in the logs.
Run curl -H 'Content-Type: application/json' -d '{ "url": "http://myshorturl.com/aHR0cHM6Ly9nbWFpbC5jb20="}' -X PUT http://localhost:8080/redirecturl to redirect to long URL. Note: replace the short URL appropriately. It will be printed on the console in the logs.
curl -H 'Content-Type: application/json' -d '{ "url": "http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t", "accesstime": "1 hours"}' -X GET http://localhost:8080//urlaccessed to get access time of a short URL.
