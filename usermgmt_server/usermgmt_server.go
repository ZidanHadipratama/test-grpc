package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/ZidanHadipratama/UTS_Mochamad-Zidan-Hadipratama_5027221052/usermgmt"
	"google.golang.org/grpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	port = ":50051"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.UserResponse, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100))
	user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}

	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		return &pb.UserResponse{User: user, Success: false, Message: "Failed to insert into MongoDB: " + err.Error()}, nil
	}

	log.Printf("Inserted a single document: %v", insertResult.InsertedID)

	return &pb.UserResponse{User: user, Success: true, Message: "Successfully inserted into MongoDB"}, nil
}

func (s *UserManagementServer) GetUser(ctx context.Context, in *pb.UserId) (*pb.UserResponse, error) {
	filter := bson.D{{"id", in.GetId()}}
	var result pb.User
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return &pb.UserResponse{User: nil, Success: false, Message: "Failed to get user from MongoDB: " + err.Error()}, nil
	}

	return &pb.UserResponse{User: &result, Success: true, Message: "Successfully got user from MongoDB"}, nil
}

func (s *UserManagementServer) GetUsersByName(ctx context.Context, in *pb.UserName) (*pb.UsersResponse, error) {
	filter := bson.D{{"name", in.GetName()}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return &pb.UsersResponse{Users: nil, Success: false, Message: "Failed to get users from MongoDB: " + err.Error()}, nil
	}
	var users []*pb.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}

	return &pb.UsersResponse{Users: users, Success: true, Message: "Successfully got users from MongoDB"}, nil
}

func (s *UserManagementServer) DeleteUser(ctx context.Context, in *pb.UserId) (*pb.UserResponse, error) {
	filter := bson.D{{"id", in.GetId()}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return &pb.UserResponse{User: nil, Success: false, Message: "Failed to delete user from MongoDB: " + err.Error()}, nil
	}

	return &pb.UserResponse{User: nil, Success: true, Message: "Successfully deleted user from MongoDB"}, nil
}

func (s *UserManagementServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	filter := bson.D{{"id", in.GetUser().GetId()}}
	update := bson.D{{"$set", bson.D{{"name", in.GetUser().GetName()}, {"age", in.GetUser().GetAge()}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &pb.UserResponse{User: nil, Success: false, Message: "Failed to update user in MongoDB: " + err.Error()}, nil
	}

	return &pb.UserResponse{User: in.GetUser(), Success: true, Message: "Successfully updated user in MongoDB"}, nil
}

func connectToMongoDB(ctx context.Context) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	database = client.Database("test")
	collection = database.Collection("users")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("server listening at %v", lis.Addr())

	connectToMongoDB(context.Background())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
