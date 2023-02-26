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
| Variable       | Type    | Example         |
|----------------|---------|-----------------|
| DnsServer      | string  | 8.8.8.8         |
| Port           | string  | 8000            |
| EmailPass      | string  | youremailpass   |
| EmailFrom      | string  | from@gmail.com  |
| EmailPort      | string  | 583             |
| EmailServer    | string  | smtp.gmail.com  |
| OpsgenieSecret | string  | examplesecret   |
| SentryDsn      | string  | exapledsnvalue  |
| Mode           | string  | prod            |
| DbUrl          | string  | host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable             |




## Add URLs to check
To add a URL to minitoring is pertty simple. Create **settings.yml** to add your endpoints to monitor. See an example of **settings.yml** below  
```sh
opegenie:
  enabled: false
email:
  enabled: true
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
Run the postgres docker container (skip this step if you already have a database)
```sh  
$ docker-compose up -d
```  
Export ENV variables (Sentry will work only in dev mode)    
If Email is disbled on **settings.yml** the following variables are not needed.
1. EmailPass
2. EmailFrom
3. EmailPort
4. EmailServer

**SentryDsn** is only needed when **Mode=dev**

```sh
$ export DnsServer="8.8.8.8" #requied  
$ export DbUrl="host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable"  #requied          
$ export EmailPass="your-pass-here" #requied when Email is Enabled  
$ export EmailFrom="youremail@server.com" #requied when Email is Enabled     
$ export EmailPort="587" #requied when Email is Enabled     
$ export EmailServer="smtp.server-address-here.com" #requied when Email is Enabled   
$ export OpsgenieSecret="your-opsgenie-key-here" #requied when Opsgenie is Enabled on settings.yml
$ export Mode="dev"  
$ export SentryDsn="your-DSn-key-here" #requied when Mode="dev"           
```
Run the project
```sh    
$ go run cmd/main.go  
```

## CI

we are using concourse CI for  Main Branch

For release branch we have Github Actions




