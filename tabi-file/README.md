# TABI - Microservice for tabi-file

## Prerequisites

- [Go](https://golang.org/doc/install) ^1.12.4
- [Docker](https://docs.docker.com/install/) ^18.09.2
- [Docker Compose](https://docs.docker.com/compose/install/) ^1.23.2
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv1.html) ^1.16.14
- [go-swagger](https://goswagger.io/install.html#homebrewlinuxbrew) ^0.21.0

## Install serverless framework

- Global:
   ```
   npm install -g serverless
   ```
- Local:
   ```
   npm i --save-dev serverless
   ```
- Specific version:
   ```
   npm install -g serverless@3.33.0
   npm i --save-dev serverless@3.33.0

## Install plugin

- Hooks Plugin:
   ```
   npm install --save-dev serverless-hooks-plugin
   ```

## Getting started

copy `.env.sample` to `.env.local`

1. Set up go env
   ```
   export GOPRIVATE=github.com/namhoai1109/tabi
   ```
2. Initialize the app for the first time:
   ```
   make depends
   make mod
   make provision
   ```
3. Generate swagger API docs:
   ```
   make specs
   ```
4. Run the development server:
   ```
   make start
   ```
   
## Other cmd

- To build diagram
   ```
   make build.diagram
   ```
- To migrate db
   ```
   make migrate
   ```
- To undo migrate db
   ```
   make migrate.undo
   ```
- To deploy
   ```
   make deploy
   ```

Read `Makefile` for more

The application runs as an HTTP server at port 3000.


