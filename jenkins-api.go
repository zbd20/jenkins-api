package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"log"
)

// Initialize Jenkins API
func Init(connection *Connection) *JenkinsApi {
	jenkinsApi := new(JenkinsApi)
	jenkinsApi.connection = connection
	jenkinsApi.client = &http.Client{}
	return jenkinsApi
}

// Get job of specific project and by job number
func (jenkinsApi *JenkinsApi) GetJob(project string, num int) *Job {

	// build endpoint url
	url := fmt.Sprintf("%v/job/%v/%v/api/json", jenkinsApi.connection.BaseUrl, project, num)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	r.SetBasicAuth(jenkinsApi.connection.Username, jenkinsApi.connection.AccessToken)
	resp, err := jenkinsApi.client.Do(r)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		log.Fatal("status: 401")
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	job := new(Job)
	err = json.Unmarshal(body, &job)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return job
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

// Custom error
type JenkinsApiError struct {
	What string
}

func (e JenkinsApiError) Error() string {
	return fmt.Sprintf("%v", e.What)
}