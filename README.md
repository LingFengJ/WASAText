# WASAText

Connect with your friends effortlessly using WASAText!  Send and receive messages, whether one-on-one  
or in groups, all from the convenience of your PC. Enjoy seamless conversations with text or GIFs and  
easily stay in touch through your private chats or group discussions.  

## Setup Instructions

Clone the repository:
```bash
git clone https://github.com/LingFengJ/WASAText.git
cd WASAText
```

## Install dependencies
Install Go dependencies
```bash
go mod tidy
go mod vendor
```

Then in webui folder:
```sh
npm install
```

## How to build container images (Docker)

### Backend

```sh
docker build -t wasatext-backend:latest -f Dockerfile.backend .
```

### Frontend

```sh
docker build -t wasatext-frontend:latest -f Dockerfile.frontend .
```

## How to run container images

### Backend

```sh
docker run -it --rm -p 3000:3000 wasatext-backend:latest
```

### Frontend

```sh
docker run -it --rm -p 8080:80 wasatext-frontend:latest
```

## How to build for Development
If you don't want to embed the WebUI into the final executable:
```bash
go build ./cmd/webapi
```
If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-embed
exit
# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

## How to run for development
You can launch the backend by only using:
```shell
go run ./cmd/webapi
```

### Launch the Frontend using Webui
Open a new tab and run the command:
```bash
cd webui
npm install
npm run dev
```