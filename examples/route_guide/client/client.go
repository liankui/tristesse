package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/liankui/tristesse/examples/route_guide/routeguide"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"time"
)

var serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")

func main() {
	// Simple RPC
	var opts []grpc.DialOption
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewRouteGuideClient(conn)
	feature, err := client.GetFeature(context.Background(), &pb.Point{Latitude: 409146138, Longitude: -746188906})
	if err != nil {
		log.Println(err)
	}
	log.Println(feature)
	fmt.Println("----------------------")

	// Server-side streaming RPC
	rect := &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	}  // initialize a pb.Rectangle
	stream, err := client.ListFeatures(context.Background(), rect)
	if err != nil {
		log.Println(err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		log.Println(feature)
	}
	fmt.Println("----------------------")

	// Client-side streaming RPC
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points.", len(points))
	stream3, err := client.RecordRoute(context.Background())
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}
	for _, point := range points {
		if err := stream3.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}
	reply, err := stream3.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)
	fmt.Println("----------------------")

	// Bidirectional streaming RPC
	stream4, err := client.RouteChat(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream4.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	//for _, note := range notes {
	//	if err := stream4.Send(note); err != nil {
	//		log.Fatalf("Failed to send a note: %v", err)
	//	}
	//}
	stream4.CloseSend()
	<-waitc

}


func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}