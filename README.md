# sky-meter
[![CodeQL](https://github.com/sooraj-sky/sky-meter/actions/workflows/codeql.yml/badge.svg)](https://github.com/sooraj-sky/sky-meter/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/sooraj-sky/sky-meter/actions/workflows/dependency-review.yml/badge.svg?branch=main)](https://github.com/sooraj-sky/sky-meter/actions/workflows/dependency-review.yml)

Sky-meter is a synthetic endpoint checker. You can deploy this on your infra and run checks from your infa and set alerts. Here we are using the go httptrace library.  
Currenly we have addded Database support. The endpoints and http output are now being saved in Database. We also have a sentry integration to catch the runtime errors.
 Development is in progress
 ### [Visit Our Website](https://sky-meter.skywalks.in)   
### [Visit pkg.go.dev](https://pkg.go.dev/github.com/sooraj-sky/sky-meter)

 ## Alerting
 We have integrated SMTP and Opsgenie, more integrations are in pipeline
 Currently the project is under developmet. You may have to experience some glitches at this moment.

## Tested Environments
GO version: 1.18  
Postgres : 15.0 
### Tested OS
- Ubuntu 22.10 
- alpine(docker)
- Macos

We are highly recommending to run th app as docker container. 
See Docker Hub Image 
https://hub.docker.com/r/soorajsky/sky-meter

## Environment variables
Currenly we have two environment variables.  
1. sentry_dsn
2. PORT

- You can create sentry project and imput the **sentry_dsn** as env variable.  
- You can export the **PORT** variable to set the http port of the server

## Add URLs to check
To add a URL to minitoring is pertty simple. Create **settings.yml** to add your endpoints to monitor. See an example of **settings.yml** below  
```sh
opegenie:
- enabled: false
email:
- enabled: true
  server: smtp.gmail.com
  port: 587
  sender: testemail@gmail.com
groups:
- name: prod
  emails:
     - reviceremail@gmail.com
     - reciver@yahoo.com
- name: dev
  emails:
     - reviceremail@gmail.com
     - reciver@yahoo.com
domains:
- name: https://skywalks.in
  enabled: true
  timeout: 10
  skip_ssl: false
  frequency: 10
  group: dev
- name: https://sky-meter.skywalks.in
  enabled: true
  timeout: 10
  skip_ssl: false
  frequency: 60
  group: dev
- name: https://github.com
  enabled: true
  timeout: 10
  skip_ssl: false
  frequency: 60
  group: prod

- name: https://githcccubs.com
  enabled: true
  timeout: 10
  skip_ssl: false
  frequency: 60
  group: prod
```
> _timeout_ : Timeout of request in Millisecond (int)  
> _skip_ssl_ : set flase if you want to skip the ssl verification (bool)  
> _frequency_ : frequency of health check in secont (int)  
> _group_ : Group settings

## Run the Code
Clone the code
```sh  
$ git clone https://github.com/sooraj-sky/sky-meter.git
$ cd sky-meter
```  
Run the postgres docker container
```sh  
$ docker-compose up -d
```  
Added Env Option: You can enable sentry by adding

Export sentry dsn  
```sh
$ export mode="dev"
$ export sentry_dsn="<yourDsnHere>"
```  
Export the port
```sh
$ export PORT=8080
```
Export opsgenieSecret
```sh
$ export opsgeniesecret="your-opsgenie-api-keyhere"
```
Export email passsword
```sh
export emailpass="your-email-pass-here"
```
Run the project
```sh    
$ go run cmd/main.go  
```

## CI

we are using concourse CI for  Main Branch

For release branch we have Github Actions




