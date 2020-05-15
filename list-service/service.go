package main

import (
	"context"
	"fmt"
	pb "simpletodo/list-service/proto/list"

	"github.com/go-acme/lego/log"
)

type service struct {
	repo repository
}

func (s *service) CreateList(ctx context.Context, req *pb.CreateListsRequest) (*pb.CreateListResponse, error) {
	log.Printf("Recieved CreateList request with params: context - %v, ListName - %v, ListUser - %v \n", ctx, req.ListName, req.ListUser)
	newListId, err := s.repo.InsertNewList(req.ListName, req.ListUser)
	if err != nil {
		log.Printf("Error encountered in CreateList: %v\n", err)
		return nil, err
	}

	return &pb.CreateListResponse{ListId: newListId}, nil
}

func (s *service) AddTasksToList(ctx context.Context, req *pb.AddTasksToListRequest) (*pb.AddTasksToListResponse, error) {
	log.Printf("Recieved AddTasksToList request with params: context - %v, list ID - %v, ListUser - %v \n", ctx, req.ListId, req.TaskDescriptions)
	var newTaskIds []int64
	for _, taskDescription := range req.TaskDescriptions {
		newTaskId, err := s.repo.InsertNewTask(req.ListId, taskDescription)
		if err != nil {
			log.Printf("Error encountered in CreateList: %v\n", err)
			return nil, err
		}
		newTaskIds = append(newTaskIds, newTaskId)
	}
	return &pb.AddTasksToListResponse{TaskIds: newTaskIds}, nil
}

func (s *service) GetLists(ctx context.Context, req *pb.MultiListRequest) (*pb.MultiList, error) {
	log.Printf("Recieved GetLists request\n")
	var lists []*pb.List
	for _, listId := range req.ListIds {
		list, err := s.repo.GetList(listId)
		if err != nil {
			log.Printf("Error encountered in GetList info: %v\n", err)
			return nil, err
		}
		tasks, err := s.repo.GetTaskLists(listId)
		if err != nil {
			log.Printf("Error encountered in Getting tasks for list: %v\n", err)
			return nil, err
		}
		list.ListTasks = tasks
		lists = append(lists, list)
	}
	return &pb.MultiList{Lists: lists}, nil
}

func (s *service) ToggleTasks(ctx context.Context, req *pb.ToggleTaskRequest) (res *pb.ToggleTaskResponse, err error) {
	log.Printf("Recieved ToggleTasks request\n")
	var updatedTasks []*pb.Task
	for _, taskId := range req.TaskIds {
		toggledTask, err := s.repo.ToggleTask(taskId)
		if err != nil {
			log.Printf("Error encountered in ToggleTasks for task ID %v: %v\n", taskId, err)
			return nil, err
		}
		updatedTasks = append(updatedTasks, toggledTask)
	}
	return &pb.ToggleTaskResponse{ListTasks: updatedTasks}, nil
}

func (s *service) DeleteLists(ctx context.Context, req *pb.MultiListRequest) (res *pb.BaseResponse, err error) {
	for _, listId := range req.ListIds {
		if err := s.repo.DeleteList(listId); err != nil {
			log.Printf("Error encountered in DeleteLists for list ID %v: %v\n", listId, err)
			return nil, err
		}
	}
	return &pb.BaseResponse{Message: fmt.Sprintf("Successful deletion of lists\n")}, nil
}
