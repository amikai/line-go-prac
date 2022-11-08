# line-go-prac

## About The Project
This is homework about implementing echo line bot.
The line bot will echo your text message.

![IMG_5280](https://user-images.githubusercontent.com/13815956/200524117-a29b550d-68ea-4290-aeed-05c6363950db.jpg)

## Built With
The following are the packages used in this Golang project.

- [gin-gonic](https://github.com/gin-gonic/gin): as http framework
- [go-migrate](https://github.com/golang-migrate/migrate): to create mongo collections
- [cobra](https://github.com/spf13/cobra): to create subcommands
- [viper](https://github.com/spf13/viper): to read environment variables and config file.
- [golangci](https://github.com/golangci/golangci-lint): to check the style


## Setup line bot server 
Before starting this program, please make sure that your docker is running.

1. Get the line channel token and secret from Line Developer Console.

2. Get token and secret from step 1 and set `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_TOKEN` as environment variable.
```
export LINE_CHANNEL_SECRET=YOUR_CHANNEL_SECRET
export LINE_CHANNEL_TOKEN=YOUR_CHANNEL_TOKEN
```

3. Run line bot server using docker-compose
```
make dc.run
```

The server serve on `9999` port. 
The api endpoint:
```
http://URL:9999/linebot/webhook
http://URL:9999/linebot/user/${userID}?count=${size}&after=${timestamp} -> list the user image by userID 
```

4. Use ngrok to expose to public network and past the endpoint in Line Developer Console.
```
ngork http 9999
```
example: If ngrok provided url is `https://XXXX.jp.ngrok.io`, You need to paste follwing webhook url to Line Developer Console
```
https://XXXX.jp.ngrok.io//linebot/webhook
```

5. Chat with line, the server will echo your text message.



## Others
### Style Check
This project uses [golangci](https://github.com/golangci/golangci-lint) to check the style.
1. Check golang coding style in docker 
```
make d.lint
```

2. Check golang coding style in host machine (need to install golangci-lint first)
```
make lint
```

### Build Image
1. Build docker image for line bot server
```
make image
```



## Future Work
1. Instrument the service by opentelemetry.
2. Use Kubernetes to manage it instead of docker-compose.
3. Write test for core logic.

