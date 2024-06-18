package jobs

// Interface that defines the methods that a job should have
type CoreJobExcecutor interface {
	Start(fileName string)
	AssignProcessId()
	ExportProcessJsonFiles()
}

type JobStatus struct {
	ProcessId string
	Started   bool
	Finished  bool
}
