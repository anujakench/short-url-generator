# Build and Deploy the application -
1. On a machine that has Windows 10 OS installed, install git bash, Go1.22.0.<br />
2. Open git bash terminal.<br />
3. Clone the github repository https://github.com/anujakench/short-url-generator.git on the machine(inside the home directory).<br />
4. cd to the repository directory.<br />
5. Follow the instructions here
```
https://etcd.io/docs/v3.4/install/
```
under "Build From Source" to install and build Etcd.<br />
6. Open git bash terminal.<br />
7. Go to etcd/bin directory inside short-url-generator.<br />
8. Run following command to start Etcd server<br />
```
etcd
```
9. Open git bash terminal and run the following command.<br />
```
go build -o short-url main.go
./short-url
```

# Manual steps to test -
1. Open git bash terminal and run the following command to create short URL.<br />
```
curl -H 'Content-Type: application/json' -d '{ "url": "https://gmail.com"}' -X POST http://localhost:8080/shortenurl
```
2. Run the following command to delete a short URL. Note: replace the short URL appropriately. It will be printed on the console in the logs.<br />
```
curl -H 'Content-Type: application/json' -d '{ "url": "http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t"}' -X POST http://localhost:8080/deleteurl
```
3. Run the following command to redirect to long URL. Note: replace the short URL appropriately. It will be printed on the console in the logs.<br />
```
curl -H 'Content-Type: application/json' -d '{ "url": "http://myshorturl.com/aHR0cHM6Ly9nbWFpbC5jb20="}' -X PUT http://localhost:8080/redirecturl
```
4. Run the following command to get access time of a short URL.<br />
```
curl -H 'Content-Type: application/json' -d '{ "url": "http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t", "accesstime": "1 hours"}' -X GET http://localhost:8080//urlaccessed
```
# Test Data Persistence -
1. Open git bash terminal.
2. cd to short-url-generator/etcd/bin.
3. Run etcdctl to access keys
   ```
   Anuja Kench@ANUJA-IBM MINGW64 ~/Desktop/short-url-generator/etcd/bin ((v3.4.28))
   $ ./etcdctl get http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t_url
   http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t_url
   http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t

   Anuja Kench@ANUJA-IBM MINGW64 ~/Desktop/short-url-generator/etcd/bin ((v3.4.28))
   $ ./etcdctl get http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t_longurl
   http://myshorturl.com/aHR0cHM6Ly9nb29nbGUuY29t_longurl
   https://google.com
  ```
