package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) runOrgPrint(ctx context.Context) error {
	orgs, err := s.org.listLocations(ctx)
	if err != nil {
		return err
	}

	for _, org := range orgs {
		if org.GetSlots() > 1 {
			pieces, err := s.org.listLocation(ctx, org.GetName())
			if err != nil {
				return err
			}
			lines := []string{org.GetName()}
			for slot := int32(0); slot < org.GetSlots(); slot++ {
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

				rec, err := s.org.getRecord(ctx, lowiid)
				if err != nil {
					return err
				}
				lines = append(lines, fmt.Sprintf("  %v - %v", rec.GetRelease().GetArtists()[0].GetName(), rec.GetRelease().GetTitle()))

				rec, err = s.org.getRecord(ctx, highiid)
				if err != nil {
					return err
				}
				lines = append(lines, fmt.Sprintf("  %v - %v", rec.GetRelease().GetArtists()[0].GetName(), rec.GetRelease().GetTitle()))
			}

			s.print(ctx, lines)
		}

	}

	return nil
}
