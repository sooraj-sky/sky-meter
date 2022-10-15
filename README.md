# sky-meter

sky-meter is an synthetic endpoint checker. You can deploy this on your infra and run checks from your infa and set alerts.

### Tested Environment
GO version: 1.18

Tested OS: Ubuntu 22.10, alpine(docker)

We are highly recommending to run th app as docker contianer. 
See Docker Hub Image 
https://hub.docker.com/r/soorajsky/sky-meter


default port : 8000

## CI

we are sing concourse CI for  Main Branch

For release branch we have Github Actions


install : https://concourse-ci.org/install.html

## Environment variables

create credentials.yml and add the following variables

docker-hub-email: <your-hub-email>
docker-hub-username: <your-hub-user>
docker-hub-password: <your-hub-pass>

fly -t tutorial set-pipeline -p skymeter-build -c concourse-pipeline.yaml -l credentials.yml