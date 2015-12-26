# jenkins-api
Go library that talks with Jenkins API

#### Install

`go get github.com/medisafe/jenkins-api/jenkins`

#### Usage 

##### Initialize api connection
``` Go
jenkinsApi := Init(&Connection{
	Username: "sromku",
	AccessToken: "001122334455667788",
	BaseUrl: "http://jenkins.sample.com:8080",
})
```

##### Get all jobs

``` Go
jobs, err := jenkinsApi.GetJobs()
```

##### Get job

``` Go
job, err := jenkinsApi.GetJob("android-mobile")
```

##### Get build 

For example, build number 196 of job called `android-mobile`
``` Go
build, err := jenkinsApi.GetBuild("android-mobile", 196)
```

##### Start new build

With params:
``` Go
jenkinsApi.StartBuild("android-mobile", map[string]interface{} {
	"branch" : "master",
	"build" : "staging",
})
```

##### Build details

- Get user that triggered the build: 

	`user, err := build.GetUser()`

- Get upstream job that triggered the build: 

	`upstream, err := build.GetUpstreamJob()`

- Get param values by param name: 

	`branchName, err := build.GetParamString("branch")`

- Get tests results:

	`testResults, err := build.GetTestResults()`

[Documentation](https://godoc.org/github.com/medisafe/jenkins-api/jenkins)