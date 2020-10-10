package main

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	gdpb "github.com/brotherlogic/godiscogs"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	ropb "github.com/brotherlogic/recordsorganiser/proto"
	rpb "github.com/brotherlogic/reminders/proto"
)

type testOrg struct {
	failLists bool
	failList  bool
	failLow   bool
	failHigh  bool
}

func (t *testOrg) listLocations(ctx context.Context) ([]*ropb.Location, error) {
	if t.failLists {
		return nil, fmt.Errorf("Built to fail")
	}
	return []*ropb.Location{
		&ropb.Location{
			Name:  "Records",
			Slots: 2,
		},
	}, nil
}

func (t *testOrg) listLocation(ctx context.Context, location string) (*ropb.Location, error) {
	if t.failList {
		return nil, fmt.Errorf("Built to fail")
	}
	return &ropb.Location{
		ReleasesLocation: []*ropb.ReleasePlacement{
			&ropb.ReleasePlacement{
				Index:      1,
				Slot:       0,
				InstanceId: 1,
			},
			&ropb.ReleasePlacement{
				Index:      2,
				Slot:       0,
				InstanceId: 2,
			},
			&ropb.ReleasePlacement{
				Index:      1,
				Slot:       1,
				InstanceId: 1,
			},
			&ropb.ReleasePlacement{
				Index:      2,
				Slot:       1,
				InstanceId: 2,
			},
		},
	}, nil
}

func (t *testOrg) getRecord(ctx context.Context, iid int32) (*rcpb.Record, error) {
	if iid == 2 {
		if t.failHigh {
			return nil, fmt.Errorf("Built to fail")
		}
		return &rcpb.Record{
			Release: &gdpb.Release{
				Title: "Last",
				Artists: []*gdpb.Artist{
					&gdpb.Artist{
						Name: "Artist",
					},
				},
			},
		}, nil
	}
	if t.failLow {
		return nil, fmt.Errorf("Built to fail")
	}
	return &rcpb.Record{
		Release: &gdpb.Release{
			Title: "First",
			Artists: []*gdpb.Artist{
				&gdpb.Artist{
					Name: "Artist",
				},
			},
		},
	}, nil
}

func TestIntegration(t *testing.T) {
	s := InitTest()
	_, err := s.Receive(context.Background(), &rpb.ReceiveRequest{})

	if err != nil {
		t.Fatalf("Unable to run basic process: %v", err)
	}

	if len(s.lastprint) != 7 {
		t.Fatalf("Not the right number of lines: %v", len(s.lastprint))
	}

	if s.lastprint[0] != "Records" ||
		s.lastprint[1] != " Slot 1:" ||
		s.lastprint[2] != "  Artist - First" ||
		s.lastprint[3] != "  Artist - Last" {
		t.Errorf("Bad print request: %v", s.lastprint)
	}
}

func TestFailLists(t *testing.T) {
	s := InitTest()
	s.org = &testOrg{failLists: true}
	_, err := s.Receive(context.Background(), &rpb.ReceiveRequest{})

	if err == nil {
		t.Fatalf("Failure did not fail")
	}
}

func TestFailList(t *testing.T) {
	s := InitTest()
	s.org = &testOrg{failList: true}
	_, err := s.Receive(context.Background(), &rpb.ReceiveRequest{})

	if err == nil {
		t.Fatalf("Failure did not fail")
	}
}

func TestFailListLow(t *testing.T) {
	s := InitTest()
	s.org = &testOrg{failLow: true}
	_, err := s.Receive(context.Background(), &rpb.ReceiveRequest{})

	if err == nil {
		t.Fatalf("Failure did not fail")
	}
}

func TestFailListHigh(t *testing.T) {
	s := InitTest()
	s.org = &testOrg{failHigh: true}
	_, err := s.Receive(context.Background(), &rpb.ReceiveRequest{})

	if err == nil {
		t.Fatalf("Failure did not fail")
	}
}
