List of steps for adding a new endpoint/feature
1. update list.proto
2. From `list-service/`, run `bash proto/list/protoc.sh`
3. In service.go, add the new method placed into the proto file
4. Implement the DB access and data manipulation logic from the service file in the repository file.
5. From `list-service/`, run `go build . && go run simpletodo/list-service`
6. Implement the rest endpoint in `list-client/`, then run the client via `go run main.go`
7. Test locally; ensure the rest client is pointed at the local address, then use postman to test the new endpoint

8. Push the new container image; From `list-service/`, run `bash deploy.sh`
9. Delete the old container containing the server
10. from `list-client/`, run `gcloud app deploy` to push up the new REST client
11. Once all the new systems are confirmed to be set up, test the new endpoint with postman