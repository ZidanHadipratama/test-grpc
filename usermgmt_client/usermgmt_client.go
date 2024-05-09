package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/tech-with-moss/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func createUser(ctx context.Context, client pb.UserManagementClient, name string, age int32) {
	r, err := client.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	if r.GetSuccess() {
		log.Printf(`User Details:
NAME: %s
AGE: %d
ID: %d`, r.GetUser().GetName(), r.GetUser().GetAge(), r.GetUser().GetId())
	} else {
		log.Printf("Failed to create user: %s", r.GetMessage())
	}
}

func getUser(ctx context.Context, client pb.UserManagementClient, id int32) {
	r, err := client.GetUser(ctx, &pb.UserId{Id: id})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	if r.GetSuccess() {
		log.Printf(`User Details:
NAME: %s
AGE: %d
ID: %d`, r.GetUser().GetName(), r.GetUser().GetAge(), r.GetUser().GetId())
	} else {
		log.Printf("Failed to get user: %s", r.GetMessage())
	}
}

func getUsersByName(ctx context.Context, client pb.UserManagementClient, name string) {
	r, err := client.GetUsersByName(ctx, &pb.UserName{Name: name})
	if err != nil {
		log.Fatalf("could not get users: %v", err)
	}
	if r.GetSuccess() {
		for _, user := range r.GetUsers() {
			log.Printf(`User Details:
NAME: %s
AGE: %d
ID: %d`, user.GetName(), user.GetAge(), user.GetId())
		}
	} else {
		log.Printf("Failed to get users: %s", r.GetMessage())
	}
}

func deleteUser(ctx context.Context, client pb.UserManagementClient, id int32) {
	r, err := client.DeleteUser(ctx, &pb.UserId{Id: id})
	if err != nil {
		log.Fatalf("could not delete user: %v", err)
	}
	if r.GetSuccess() {
		log.Printf("Successfully deleted user")
	} else {
		log.Printf("Failed to delete user: %s", r.GetMessage())
	}
}

func updateUser(ctx context.Context, client pb.UserManagementClient, id int32, name string, age int32) {
	r, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{User: &pb.User{Id: id, Name: name, Age: age}})
	if err != nil {
		log.Fatalf("could not update user: %v", err)
	}
	if r.GetSuccess() {
		log.Printf(`User Details:
NAME: %s
AGE: %d
ID: %d`, r.GetUser().GetName(), r.GetUser().GetAge(), r.GetUser().GetId())
	} else {
		log.Printf("Failed to update user: %s", r.GetMessage())
	}
}

func main() {
	createFlag := flag.Bool("create", false, "Create new user")
	getFlag := flag.Int("get", -1, "Get user by ID")
	getNameFlag := flag.String("getname", "", "Get users by name")
	deleteFlag := flag.Int("delete", -1, "Delete user by ID")
	nameFlag := flag.String("name", "", "Name of the user")
	ageFlag := flag.Int("age", 0, "Age of the user")
	updateFlag := flag.Int("update", -1, "Update user by ID")
	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if *createFlag {
		createUser(ctx, c, *nameFlag, int32(*ageFlag))
	} else if *getFlag != -1 {
		getUser(ctx, c, int32(*getFlag))
	} else if *getNameFlag != "" {
		getUsersByName(ctx, c, *getNameFlag)
	} else if *deleteFlag != -1 {
		deleteUser(ctx, c, int32(*deleteFlag))
	} else if *updateFlag != -1 {
		updateUser(ctx, c, int32(*updateFlag), *nameFlag, int32(*ageFlag))
	} else {
		log.Fatalf("Invalid command")
	}
}
