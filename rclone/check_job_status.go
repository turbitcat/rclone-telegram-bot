package rclone

import (
	"encoding/json"
)

type JobStatus struct {
	Duration         float32                   // time in seconds that the job ran for
	EndTime          string                    // time the job finished (e.g. "2018-10-26T18:50:20.528746884+01:00")
	Error            string                    // error from the job or empty string for no error
	Finished         bool                      // boolean whether the job has finished or not
	Id               int                       // as passed in above
	StartTime        string                    // time the job started (e.g. "2018-10-26T18:50:20.528336039+01:00")
	Success          bool                      // boolean - true for success false otherwise
	Eta              float32                   // estimated time in seconds until the group completes,
	Speed            float32                   // average speed in bytes per second since start of the group,
	TotalBytes       int                       // total number of bytes in the group,
	TransferredBytes int                       // total transferred bytes since the start of the group,
	Transferring     []StatsResponsTransfering // an array of currently active file transfers,
	Checking         []string                  // an array of names of currently active file checks
}

func (sc *ServerConfig) CheckJobStatus(jobid int) (*JobStatus, error) {
	body, err := sc.Do("job/status", JobStatusRequest{JobID: jobid})
	if err != nil {
		return nil, err
	}
	jResp := JobStatusResponse{}
	if err := json.Unmarshal(body, &jResp); err != nil {
		return nil, err
	}
	var coreStats *StatsResponse = nil
	if jResp.Group != "" && !jResp.Finished {
		body, err := sc.Do("core/stats", StatsRequest{jResp.Group})
		if err == nil {
			respons := StatsResponse{}
			if err := json.Unmarshal(body, &respons); err == nil {
				coreStats = &respons
			}
		}
	}
	jobStatus := JobStatus{
		Duration:  jResp.Duration,
		EndTime:   jResp.EndTime,
		Error:     jResp.Error,
		Finished:  jResp.Finished,
		Id:        jResp.Id,
		StartTime: jResp.StartTime,
		Success:   jResp.Success,
	}
	if coreStats != nil {
		jobStatus.Eta = coreStats.Eta
		jobStatus.Speed = coreStats.Speed
		jobStatus.TotalBytes = coreStats.TotalBytes
		jobStatus.TransferredBytes = coreStats.Bytes
		jobStatus.Transferring = coreStats.Transferring
		jobStatus.Checking = coreStats.Checking
	}
	return &jobStatus, nil
}
