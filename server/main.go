package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "github.com/moemoe89/go-grpc-upload/files/go/file"
	"github.com/moemoe89/go-helpers/pkg/diskstorage"

	"google.golang.org/grpc"
)

// server is used to implement File service.
type server struct {
	pb.UnimplementedFileServiceServer
}

// Upload uploads file to the server.
func (s *server) Upload(stream pb.FileService_UploadServer) error {
	var offset int64

	f, err := diskstorage.New()
	if err != nil {
		return err
	}

	filename := ""

	for {
		chunk, err := stream.Recv()

		if filename == "" {
			filename = chunk.GetFilename()
		}

		if err == io.EOF {

			if err := f.WriteFile("uploads/"+filename, os.FileMode(0644)); err != nil {
				return err
			}

			return stream.SendAndClose(&pb.Empty{})
		}

		if err != nil {
			return err
		}

		if offset != chunk.Offset {
			return fmt.Errorf("unexpected offset, got %d, want %d", chunk.Offset, offset)
		}

		offset += int64(len(chunk.Data))

		err = f.Write(chunk.Data)
		if err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterFileServiceServer(s, &server{})

	log.Printf("grpc server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
