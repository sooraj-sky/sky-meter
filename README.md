# sky-meter

/stats to see the output

default port : 8000

## CI

we are sing concourse for ci

install : https://concourse-ci.org/install.html

create credentials.yml and add the following variables

docker-hub-email: <your-hub-email>
docker-hub-username: <your-hub-user>
docker-hub-password: <your-hub-pass>

fly -t tutorial set-pipeline -p skymeter-build -c concourse-pipeline.yaml -l credentials.yml