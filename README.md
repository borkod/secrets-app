# secrets-app
Sample secret sharing application build as part of technical review of Manning liveProject "Build a Secrets Sharing Web Application in Go"

This is a backend-only web application that will allow creating and sharing a one-time secrets.

To build the application from source code:
```
go build -o secrets-app
```

To run the application:
```
./secrets-app
```

To create a secret, you make a POST HTTP API request to the application containing the plain text (in a new terminal). The request JSON body should look like `{"plain_text":"SecretValue"}`. The response will contain an ID that you can then share with the recipient.

```
curl -X POST http://localhost:8080 -d '{\"plain_text\":\"My super secret123\"}' -H "Content-Type: application/json"
{"id":"c616584ac64a93aafe1c16b6620f5bcd"}
```

Once the recipient makes a GET HTTP API request with the provided ID to view the secret, the secret is not viewable again. The value is sent back as a JSON response, `{"data":"SecretValue"}`
```
curl http://localhost:8080/c616584ac64a93aafe1c16b6620f5bcd
{"data":"My super secret123"}
```

The application stores secrets in memory as a map and a file on local disk. The key for the map is an md5 hash of the secret text and the value is the secret itself. Mutex is used to lock and unlock the map and file when it is being manipulated to avoid race conditions. The file location is set using `DATA_FILE_PATH` environment variable.
