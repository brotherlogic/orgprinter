package main

import (
	"fmt"

	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
	ropb "github.com/brotherlogic/recordsorganiser/proto"
)

func (s *Server) runSalePrint(ctx context.Context) error {
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	sales, err := client.GetInventory(ctx, &rcpb.GetInventoryRequest{})
	if err != nil {
		return err
	}

	lines := []string{"For Sale:"}
	for _, sale := range sales.GetItems() {
		if sale.GetSalePrice() <= 5 {
			records, err := client.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_ReleaseId{sale.GetId()}})
			if err == nil {
				record, err := client.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: records[0].GetInstanceId})
				if err == nil {
					lines = append(lines, fmt.Sprintf("%v - %v.", sale.GetId(), record.GetRecord().GetRelease().GetTitle()))
				}
			}
		}
	}

	if len(lines) > 1 {
		s.print(ctx, lines)
	}

	return nil
}

func getLine(rec *rcpb.Record) string {
	if len(rec.GetRelease().GetArtists()) > 0 {
		return fmt.Sprintf("  %v - %v", rec.GetRelease().GetArtists()[0].GetName(), rec.GetRelease().GetTitle())
	}
	return fmt.Sprintf("  %v - %v", "Unknown", rec.GetRelease().GetTitle())
}

func (s *Server) runOrgPrint(ctx context.Context) error {
	orgs, err := s.org.listLocations(ctx)
	if err != nil {
		return err
	}

	for _, org := range orgs {
		if org.GetSlots() > 1 && org.GetInPlay() == ropb.Location_IN_PLAY {
			pieces, err := s.org.listLocation(ctx, org.GetName())
			if err != nil {
				return err
			}
			lines := []string{org.GetName()}
			for slot := int32(0); slot <= org.GetSlots(); slot++ {
				lines = append(lines, fmt.Sprintf(" Slot %v:", slot+1))

				lowest := int32(9999999)
				lowiid := int32(0)
				highest := int32(0)
				highiid := int32(0)

				for _, place := range pieces.GetReleasesLocation() {
					if place.GetSlot() == slot {
						if place.GetIndex() < lowest {
							lowest = place.GetIndex()
							lowiid = place.GetInstanceId()
						}
						if place.GetInstanceId() > highest {
							highest = place.GetIndex()
							highiid = place.GetInstanceId()
						}
					}
				}

				if lowiid > 0 {
					rec, err := s.org.getRecord(ctx, lowiid)
					if err != nil {
						return err
					}
					lines = append(lines, getLine(rec))
				}

				if highiid > 0 {
					rec, err := s.org.getRecord(ctx, highiid)
					if err != nil {
						return err
					}

					lines = append(lines, getLine(rec))
				}
			}

			s.print(ctx, lines)
		}

	}

	return nil
}
