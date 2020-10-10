package main

import (
	"golang.org/x/net/context"

	rpb "github.com/brotherlogic/reminders/proto"
)

// Receive a request ot update
func (s *Server) Receive(ctx context.Context, req *rpb.ReceiveRequest) (*rpb.ReceiveResponse, error) {
	err := s.runOrgPrint(ctx)
	return &rpb.ReceiveResponse{}, err
}
