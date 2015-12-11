# jenkins-api
Go library that talks with Jenkins API

#### Install

`go get github.com/Medisafe/jenkins-api/jenkins`

#### Usage 

Initialize api connection
``` Go
jenkinsApi := Init(&Connection{
	Username: "sromku",
	AccessToken: "001122334455667788",
	BaseUrl: "http://jenkins.sample.com:8080",
})
```

Get job (for example, job number 196 of project called `alpha`)
``` Go
job := jenkinsApi.GetJob("alpha", 196)
```

Get user or upstream job that triggered the job
``` Go
user, _ := job.GetUser()
upstream, _ := job.GetUpstreamJob()
```

Get param values by param name
``` Go
branchName, _ := job.GetParamString("branch")
```

[Documentation](https://godoc.org/github.com/Medisafe/jenkins-api/jenkins)