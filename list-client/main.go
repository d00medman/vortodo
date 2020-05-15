package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"context"

	pb "simpletodo/list-service/proto/list"

	"google.golang.org/grpc"
)

const (
	// deployedAddress = "34.67.1.155:80"
	localAddress = "localhost:50051"
)

func grpcClientConnect() (client pb.ListServiceClient, conn *grpc.ClientConn) {
	address := os.Getenv("GRPC_SERVER_ADDRESS")
	if address == "" {
		address = localAddress
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to GRPC client at address %v: %v", address, err)
	} else {
		log.Printf("Connected to address %v", address)
	}
	client = pb.NewListServiceClient(conn)
	return client, conn
}

func createListHandler(w http.ResponseWriter, r *http.Request) {
	client, conn := grpcClientConnect()
	defer conn.Close()

	decoder := json.NewDecoder(r.Body)
	var createListRequest *pb.CreateListsRequest

	if err := decoder.Decode(&createListRequest); err != nil {
		panic(err)
	}
	log.Println(createListRequest)

	response, err := client.CreateList(context.Background(), createListRequest)
	if err != nil {
		log.Fatalf("Could not Create list: %v", err)
	}
	responseMsg := fmt.Sprintf("Created new list with the following ID: %v", response.ListId)
	log.Println(responseMsg)
	fmt.Fprintf(w, responseMsg)
}

func addTaskToListHandler(w http.ResponseWriter, r *http.Request) {
	client, conn := grpcClientConnect()
	defer conn.Close()

	decoder := json.NewDecoder(r.Body)
	var addTasksToListRequest *pb.AddTasksToListRequest

	if err := decoder.Decode(&addTasksToListRequest); err != nil {
		panic(err)
	}
	log.Println(addTasksToListRequest)

	response, err := client.AddTasksToList(context.Background(), addTasksToListRequest)
	if err != nil {
		log.Fatalf("Could not Create list: %v", err)
	}
	responseMsg := fmt.Sprintf("Created new task with the following ID: %v", response.TaskIds)
	log.Println(responseMsg)
	fmt.Fprintf(w, responseMsg)
}

func getListsHandler(w http.ResponseWriter, r *http.Request) {
	client, conn := grpcClientConnect()
	defer conn.Close()

	decoder := json.NewDecoder(r.Body)
	var multiListRequest *pb.MultiListRequest

	if err := decoder.Decode(&multiListRequest); err != nil {
		panic(err)
	}
	log.Println(multiListRequest)

	response, err := client.GetLists(context.Background(), multiListRequest)
	if err != nil {
		log.Fatalf("Could not Create list: %v", err)
	}
	responseMsg := fmt.Sprintf("Created new task with the following ID: %v", response.Lists)
	log.Println(responseMsg)
	fmt.Fprintf(w, responseMsg)
}

func toggleTasksHandler(w http.ResponseWriter, r *http.Request) {
	client, conn := grpcClientConnect()
	defer conn.Close()

	decoder := json.NewDecoder(r.Body)
	var toggleTasksRequest *pb.ToggleTaskRequest
	if err := decoder.Decode(&toggleTasksRequest); err != nil {
		panic(err)
	}
	log.Println(toggleTasksRequest)

	response, err := client.ToggleTasks(context.Background(), toggleTasksRequest)
	if err != nil {
		log.Fatalf("Could not toggle tasks: %v", err)
	}
	responseMsg := fmt.Sprintf("The following tasks have been toggled: %v", response.ListTasks)
	log.Println(responseMsg)
	fmt.Fprintf(w, responseMsg)
}

func deleteListHandler(w http.ResponseWriter, r *http.Request) {
	client, conn := grpcClientConnect()
	defer conn.Close()

	decoder := json.NewDecoder(r.Body)
	var deleteListsRequest *pb.MultiListRequest
	if err := decoder.Decode(&deleteListsRequest); err != nil {
		panic(err)
	}
	log.Println(deleteListsRequest)

	response, err := client.DeleteLists(context.Background(), deleteListsRequest)
	if err != nil {
		log.Fatalf("Could not delete lists: %v", err)
	}
	log.Println(response.Message)
	fmt.Fprintf(w, response.Message)
}

func main() {
	http.HandleFunc("/createList", createListHandler)
	http.HandleFunc("/addTaskToList", addTaskToListHandler)
	http.HandleFunc("/getLists", getListsHandler)
	http.HandleFunc("/toggleTasks", toggleTasksHandler)
	http.HandleFunc("/deleteLists", deleteListHandler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
