# Migration tool: from Heroku to Scalingo in a minute

## Roadmap

* First version includes:
  * Scalingo app creation
  * Get heroku app environment and set it to scalingo's new app
  * Heroku repository cloning + add of scalingo remote + push to scalingo master

## Installation

You can install the binary by:
* Downloading it manually here: [heroku2scalingo](https://github.com/Scalingo/heroku2scalingo/releases/latest)
  Unzipping it: `tar -xvf heroku2scalingo_0.1.1_linux_amd64.tar.gz`
  And placing it in one of your `$PATH`
* Building it from source:
  `git clone https://github.com/Scalingo/heroku2scalingo.git`
  `godep go build`
  And then placing `heroku2scalingo` to one of your `$PATH`

## Usage

```bash
heroku2scalingo <app_name>
```

The following operations will be performed in this order:
* Autenthication to Scalingo
* Authentication to Heroku API
* Creation of Scalingo app
* Get/Set environment variables
* `git clone` your heroku app repository
* `git push scalingo master` -> Auto-deployment using the Procfile
