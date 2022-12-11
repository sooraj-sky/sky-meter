# sky-meter

sky-meter is an synthetic endpoint checker. You can deploy this on your infra and run checks from your infa and set alerts. Here we are using the go httptrace library.  
Currenly we have addded Database support. The endpoints and http output are now bing saved in Database. We also have a sentry integration to catch the runtime errors.
 Development is in progress

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
To add a URL to minitoring is pertty simple. Create **input.json** to add your endpoints to monitor. See an example of **input.json** below  
```sh
[
    {
      "url": "https://github.com",
      "timeout": 1200,
      "skip_ssl": false,
      "frequency" : 5,
      "group": "prod"
    },
        {
      "url": "https://stackoverflow.com",
      "timeout": 1200,
      "skip_ssl": false,
      "frequency" : 30,
      "group": "prod"
    }
]
```
> _url_ : url to monitor (string)   
> _timeout_ : Timeout of request in Millisecond (int)  
> _skip_ssl_ : set flase if you want to skip the ssl verification (bool)  
> _frequency_ : frequency of health check in secont (int)  
> _group_ : Add group properties (string)

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

we are sing concourse CI for  Main Branch

For release branch we have Github Actions

install : https://concourse-ci.org/install.html

