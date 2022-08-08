package rclone

// json struct:
// core/stats
type StatsRequest struct {
	Group string `json:"group,omitempty"`
}

// json struct:
// core/stats
type StatsResponse struct {
	Bytes          int                       // total transferred bytes since the start of the group,
	Checks         int                       // number of files checked,
	Deletes        int                       // number of files deleted,
	ElapsedTime    float32                   // time in floating point seconds since rclone was started,
	Errors         int                       // number of errors,
	Eta            float32                   // estimated time in seconds until the group completes,
	FatalError     bool                      // boolean whether there has been at least one fatal error,
	LastError      string                    // last error string,
	Renames        int                       // number of files renamed,
	RetryError     bool                      // boolean showing whether there has been at least one non-NoRetryError,
	Speed          float32                   // average speed in bytes per second since start of the group,
	TotalBytes     int                       // total number of bytes in the group,
	TotalChecks    int                       // total number of checks in the group,
	TotalTransfers int                       // total number of transfers in the group,
	TransferTime   float32                   // total time spent on running jobs,
	Transfers      int                       // number of transferred files,
	Transferring   []StatsResponsTransfering // an array of currently active file transfers,
	Checking       []string                  // an array of names of currently active file checks
}

// json struct:
// core/stats
type StatsResponsTransfering struct {
	TransferredBytes int     `json:"bytes"` // total transferred bytes for this file,
	Eta              float32 // estimated time in seconds until file transfer completion
	Name             string  // name of the file,
	Percentage       int     // progress of the file transfer in percent,
	Speed            float32 // average speed over the whole transfer in bytes per second,
	SpeedAvg         float32 // current speed in bytes per second as an exponentially weighted moving average,
	ToTalBytes       int     `json:"size"` // size of the file in bytes
}
