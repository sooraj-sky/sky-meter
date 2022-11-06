# sky-meter

sky-meter is an synthetic endpoint checker. You can deploy this on your infra and run checks from your infa and set alerts. Here we are using the go httptrace library.  
Currenly we have addded Database support. The endpoints and http output are now bing saved in Database. We also have a sentry integration to catch the runtime errors.


## Tested Environments
GO version: 1.18
Tested OS: Ubuntu 22.10, alpine(docker), Macos

We are highly recommending to run th app as docker container. 
See Docker Hub Image 
https://hub.docker.com/r/soorajsky/sky-meter

## Environment variables
Currenly we have two environment variables.  
1. sentry_dsn
2. PORT

- You can create sentry project and imput the **sentry_dsn** as env variable.  
- You can export the **PORT** variable to set the http port of the server

## Run the Code
Clone the code
```sh  
$ git clone https://github.com/sooraj-sky/sky-meter.git
```  
Run the postgres docker container
```sh  
$ docker-compose up -d
```  
Export sentry dsn  
```sh
$ export sentry_dsn="<yourDsnHere>"
```  
Export the port
```sh
$ export PORT=8080
```
Run the project
```sh    
$ go run cmd/main.go  
```

## CI

we are sing concourse CI for  Main Branch

For release branch we have Github Actions

install : https://concourse-ci.org/install.html

