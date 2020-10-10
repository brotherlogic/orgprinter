package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
	ppb "github.com/brotherlogic/printer/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	ropb "github.com/brotherlogic/recordsorganiser/proto"
	rpb "github.com/brotherlogic/reminders/proto"
)

type org interface {
	listLocations(ctx context.Context) ([]*ropb.Location, error)
	listLocation(ctx context.Context, location string) (*ropb.Location, error)
	getRecord(ctx context.Context, iid int32) (*rcpb.Record, error)
}

type porg struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *porg) listLocations(ctx context.Context) ([]*ropb.Location, error) {
	conn, err := p.dial(ctx, "recordsorganiser")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	org := ropb.NewOrganiserServiceClient(conn)
	resp, err := org.GetOrganisation(ctx, &ropb.GetOrganisationRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetLocations(), err
}

func (p *porg) listLocation(ctx context.Context, location string) (*ropb.Location, error) {
	conn, err := p.dial(ctx, "recordsorganiser")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	org := ropb.NewOrganiserServiceClient(conn)
	resp, err := org.GetOrganisation(ctx, &ropb.GetOrganisationRequest{ForceReorg: true, Locations: []*ropb.Location{&ropb.Location{Name: location}}})
	if err != nil {
		return nil, err
	}
	return resp.GetLocations()[0], err
}

func (p *porg) getRecord(ctx context.Context, iid int32) (*rcpb.Record, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	coll := rcpb.NewRecordCollectionServiceClient(conn)
	resp, err := coll.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: iid})
	if err != nil {
		return nil, err
	}
	return resp.GetRecord(), err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	org       org
	runprint  bool
	lastprint []string
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		runprint: true,
	}
	s.org = &porg{dial: s.FDialServer}
	return s
}

func (s *Server) print(ctx context.Context, lines []string) error {
	if !s.runprint {
		s.lastprint = lines
		return nil
	}

	conn, err := s.FDialServer(ctx, "printer")
	if err != nil {
		return err
	}
	defer conn.Close()

	printer := ppb.NewPrintServiceClient(conn)
	_, err = printer.Print(ctx, &ppb.PrintRequest{Lines: lines, Origin: "orgprinter"})
	return err
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	rpb.RegisterReminderReceiverServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("orgprinter", false, true)
	if err != nil {
		return
	}

	fmt.Printf("%v", server.Serve())
}
