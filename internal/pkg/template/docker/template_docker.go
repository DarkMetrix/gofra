package docker

// docker file config
type DockerFileInfo struct {
	Author string
	Project string
}

var DockerFileTemplate string = `
FROM centos:latest
MAINTAINER {{.Author}}

COPY ./build /app/{{.Project}}/bin/
COPY ./configs /app/{{.Project}}/configs/

WORKDIR /app/{{.Project}}/bin

ENTRYPOINT ["./{{.Project}}"]
`

