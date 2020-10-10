package main

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.runprint = false
	s.org = &testOrg{}
	return s
}
