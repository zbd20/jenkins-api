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

##### Get job 

For example, job number 196 of project called `android-mobile`
``` Go
job := jenkinsApi.GetJob("android-mobile", 196)
```

##### Start job (new build)

With params:
``` Go
jenkinsApi.StartJob("android-mobile", map[string]interface{} {
	"branch" : "master",
	"build" : "staging",
})
```

##### Job details

- Get user that triggered the job: 

	`user, err := job.GetUser()`

- Get upstream job that triggered the job: 

	`upstream, err := job.GetUpstreamJob()`

- Get param values by param name: 

	`branchName, _ := job.GetParamString("branch")`

- Get tests results:

	``

[Documentation](https://godoc.org/github.com/medisafe/jenkins-api/jenkins)