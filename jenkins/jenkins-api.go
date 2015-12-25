package jenkins

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// Initialize Jenkins API
func Init(connection *Connection) *JenkinsApi {
	jenkinsApi := new(JenkinsApi)
	jenkinsApi.connection = connection
	jenkinsApi.client = &http.Client{}
	return jenkinsApi
}

// Get job of specific project and by job number
func (jenkinsApi *JenkinsApi) GetJob(project string, num int) (*Job, error) {

	// build endpoint url
	url := fmt.Sprintf("%v/job/%v/%v/api/json", jenkinsApi.connection.BaseUrl, project, num)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r.SetBasicAuth(jenkinsApi.connection.Username, jenkinsApi.connection.AccessToken)
	resp, err := jenkinsApi.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, JenkinsApiError{ What: fmt.Sprintf("Status code: %v", resp.StatusCode) }
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	job := new(Job)
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// Get parameter of string type
func (job *Job) GetParamString(name string) (string, error) {
	for _, action := range job.Actions {
		params := action.Parameters
		if len(params) > 0 {
			for _, param := range params {
				if param.Name == name {
					if val, ok := param.Value.(string); ok {
						return val, nil
					} else {
						return "", JenkinsApiError{ What: fmt.Sprintf("The value of '%v' isn't of string type", name) }
					}
				}
			}
		}
	}
	return "", JenkinsApiError{ What: fmt.Sprintf("Param '%v' wasn't found", name) }
}

// Get parameter of int type
func (job *Job) GetParamInt(name string) (int, error) {
	for _, action := range job.Actions {
		params := action.Parameters
		if len(params) > 0 {
			for _, param := range params {
				if param.Name == name {
					if val, ok := param.Value.(int); ok {
						return val, nil
					} else {
						return 0, JenkinsApiError{ What: fmt.Sprintf("The value of '%v' isn't of int type", name) }
					}
				}
			}
		}
	}
	return 0, JenkinsApiError{ What: fmt.Sprintf("Param '%v' wasn't found", name) }
}

// Get parameter of bool type
func (job *Job) GetParamBool(name string) (bool, error) {
	for _, action := range job.Actions {
		params := action.Parameters
		if len(params) > 0 {
			for _, param := range params {
				if param.Name == name {
					if val, ok := param.Value.(bool); ok {
						return val, nil
					} else {
						return false, JenkinsApiError{ What: fmt.Sprintf("The value of '%v' isn't of bool type", name) }
					}
				}
			}
		}
	}
	return false, JenkinsApiError{ What: fmt.Sprintf("Param '%v' wasn't found", name) }
}

// Get user that triggered this job
func (job *Job) GetUser() (*User, error) {
	for _, action := range job.Actions {
		causes := action.Causes
		if len(causes) > 0 {
			for _, cause := range causes {
				if cause.User.UserId != "" {
					return &cause.User, nil
				}
			}
		}
	}
	return nil, JenkinsApiError{ What: "User wasn't found for this job, maybe upstream job triggered this job" }
}

// Get upstream job that triggered this job
func (job *Job) GetUpstreamJob() (*UpstreamJob, error) {
	for _, action := range job.Actions {
		causes := action.Causes
		if len(causes) > 0 {
			for _, cause := range causes {
				if cause.UpstreamJob.UpstreamProject != "" {
					return &cause.UpstreamJob, nil
				}
			}
		}
	}
	return nil, JenkinsApiError{ What: "Upstream job wasn't found for this job, maybe user triggered this job" }
}

// The job can run tests part of the script. Get the tests count summary.
func (job *Job) GetTestResults() (*TestResult, error) {
	for _, action := range job.Actions {
		if action.TestResult.TotalCount > 0 {
			return &action.TestResult, nil
		}
	}
	return nil, JenkinsApiError{ What: "No tests results for this job" }
}

// Start jenkins job and pass params.
func (jenkinsApi *JenkinsApi) StartJob(project string, params map[string]interface{}) error {

	parameters := &Parameters{}
	if params != nil && len(params) > 0 {
		for k := range params {
			parameters.Params = append(parameters.Params, Parameter{ Name: k, Value: params[k]})
		}
	}

	var jsonStr string
	if len(parameters.Params) > 0 {
		jsonbts, _ := json.Marshal(parameters)
		jsonStr = string(jsonbts)
	}

	// build endpoint url
	url := fmt.Sprintf("%v/job/%v/build?json=%v", jenkinsApi.connection.BaseUrl, project, jsonStr)
	r, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	r.SetBasicAuth(jenkinsApi.connection.Username, jenkinsApi.connection.AccessToken)
	resp, err := jenkinsApi.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return JenkinsApiError{ What: fmt.Sprintf("Status code: %v", resp.StatusCode) }
	}

	return nil
}

// Custom error
type JenkinsApiError struct {
	What string
}

func (e JenkinsApiError) Error() string {
	return fmt.Sprintf("%v", e.What)
}