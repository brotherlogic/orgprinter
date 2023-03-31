package main

import recordcollection_client "github.com/brotherlogic/recordcollection/client"

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.runprint = false
	s.rcclient = recordcollection_client.RecordCollectionClient{Test: true}
	s.org = &testOrg{}
	return s
}
