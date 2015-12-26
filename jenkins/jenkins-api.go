// Go library that talks with Jenkins API.
//
// Author: sromku
package jenkins

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// Initialize Jenkins API
//
// For example:
//  jenkinsApi := Init(&Connection {
//    Username: "sromku",
//    AccessToken: "001122334455667788",
//    BaseUrl: "http://jenkins.sample.com:8080",
//  })
func Init(connection *Connection) *JenkinsApi {
	jenkinsApi := new(JenkinsApi)
	jenkinsApi.connection = connection
	jenkinsApi.client = &http.Client{}
	return jenkinsApi
}

// Get build of specific job
//
// For example, job number 196 of project called android-mobile:
//  job := jenkinsApi.GetJob("android-mobile", 196)
func (jenkinsApi *JenkinsApi) GetBuild(jobName string, num int) (*Build, error) {

	// build endpoint url
	url := fmt.Sprintf("%v/job/%v/%v/api/json", jenkinsApi.connection.BaseUrl, jobName, num)
	body, err := jenkinsApi.get(url)
	if err != nil {
		return nil, err
	}

	build := new(Build)
	err = json.Unmarshal(body, &build)
	if err != nil {
		return nil, err
	}

	return build, nil
}

// Get parameter of string type
//
// For example:
//  branchName, _ := job.GetParamString("branch")
func (build *Build) GetParamString(name string) (string, error) {
	for _, action := range build.Actions {
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
func (build *Build) GetParamInt(name string) (int, error) {
	for _, action := range build.Actions {
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
func (build *Build) GetParamBool(name string) (bool, error) {
	for _, action := range build.Actions {
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

// Get user that triggered this build
//  user, err := job.GetUser()
func (build *Build) GetUser() (*User, error) {
	for _, action := range build.Actions {
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

// Get upstream job that triggered this build
//  upstream, err := job.GetUpstreamJob()
func (build *Build) GetUpstreamJob() (*UpstreamJob, error) {
	for _, action := range build.Actions {
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

// The job can run tests part of the script. Get the tests count summary
//  testResults, err := build.GetTestResults()
func (build *Build) GetTestResults() (*TestResult, error) {
	for _, action := range build.Actions {
		if action.TestResult.TotalCount > 0 {
			return &action.TestResult, nil
		}
	}
	return nil, JenkinsApiError{ What: "No tests results for this job" }
}

// Start jenkins build and pass params.
//
// For example, start new build with two params:
//  jenkinsApi.StartBuild("android-mobile", map[string]interface{} {
//    "branch" : "master",
//    "build" : "staging",
//  })
func (jenkinsApi *JenkinsApi) StartBuild(job string, params map[string]interface{}) error {

	parameters := &Parameters{}
	if params != nil && len(params) > 0 {
		for k := range params {
			parameters.Params = append(parameters.Params, Parameter{ Name: k, Value: params[k]})
		}
	}

	var buildStr string
	if len(parameters.Params) > 0 {
		jsonbts, _ := json.Marshal(parameters)
		buildStr = string(jsonbts)
	}

	// build endpoint url
	url := fmt.Sprintf("%v/job/%v/build?json=%v", jenkinsApi.connection.BaseUrl, job, buildStr)
	err := jenkinsApi.post(url)
	if err != nil {
		return err
	}

	return nil
}

// Get job details
//  job, err := jenkinsApi.GetJob("android-mobile")
func (jenkinsApi *JenkinsApi) GetJob(jobName string) (*Job, error) {
	// build endpoint url
	url := fmt.Sprintf("%v/job/%v/api/json", jenkinsApi.connection.BaseUrl, jobName)
	body, err := jenkinsApi.get(url)
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

func (jenkinsApi *JenkinsApi) get(url string) ([]byte, error) {
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

	return body, nil
}

func (jenkinsApi *JenkinsApi) post(url string) error {
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