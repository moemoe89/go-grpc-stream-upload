package main

import (
	"context"
	"io"
	"log"
	"os"

	pb "github.com/moemoe89/go-grpc-upload/files/go/file"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	client := pb.NewFileServiceClient(conn)

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer stream.CloseSend()

	file, err := os.Open("client/logos.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var offset int64

	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		err = stream.Send(&pb.UploadRequest{
			Filename: "logos.jpg",
			Data:     buf[:n],
			Offset:   offset,
		})
		if err != nil {
			log.Fatal(err)
		}

		offset += int64(n)
	}
	
	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("### Successfully uploading file!! ###")
}
