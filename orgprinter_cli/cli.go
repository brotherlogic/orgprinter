package main

import (
	"fmt"
	"log"

	"github.com/brotherlogic/goserver/utils"

	pb "github.com/brotherlogic/reminders/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.BuildContext("recordader-cli", "recordadder")
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "orgprinter")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewReminderReceiverClient(conn)
	res, err := client.Receive(ctx, &pb.ReceiveRequest{})
	fmt.Printf("%v and %v\n", res, err)
}
