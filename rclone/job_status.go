package rclone

import "encoding/json"

// json struct:
// job/status
type JobStatusRequest struct {
	JobID int `json:"jobid"`
}

// json struct:
// job/status
type JobStatusResponse struct {
	Duration  float32        // time in seconds that the job ran for
	EndTime   string         // time the job finished (e.g. "2018-10-26T18:50:20.528746884+01:00")
	Error     string         // error from the job or empty string for no error
	Finished  bool           // boolean whether the job has finished or not
	Group     string         //
	Id        int            // as passed in above
	StartTime string         // time the job started (e.g. "2018-10-26T18:50:20.528336039+01:00")
	Success   bool           // boolean - true for success false otherwise
	Output    map[string]any // output of the job as would have been returned if called synchronously
	Progress  string         // output of the progress related to the underlying job
}

func (rs *RcloneServer) CheckJobStatus(jobid int) (*JobStatusResponse, error) {
	body, err := rs.Do(JobStatus, JobStatusRequest{JobID: jobid})
	if err != nil {
		return nil, err
	}
	jResp := JobStatusResponse{}
	if err := json.Unmarshal(body, &jResp); err != nil {
		return nil, err
	}
	return &jResp, nil
}
