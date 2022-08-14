package main

import (
	"golang.org/x/net/context"

	rpb "github.com/brotherlogic/reminders/proto"
)

// Receive a request ot update
func (s *Server) Receive(ctx context.Context, req *rpb.ReceiveRequest) (*rpb.ReceiveResponse, error) {
	err := s.runOrgPrint(ctx)
	err2 := s.runSalePrint(ctx)
	if err != nil {
		return &rpb.ReceiveResponse{}, err
	} else if err2 != nil {
		return &rpb.ReceiveResponse{}, err2
	}

	return &rpb.ReceiveResponse{}, nil
}
