# WASAText

Connect with your friends effortlessly using WASAText! Send and receive messages, whether one-on-one  
or in groups, all from the convenience of your PC. Enjoy seamless conversations with text or GIFs and  
easily stay in touch through your private chats or group discussions.  

## Setup Instructions

1. Clone the repository:
```bash
git clone https://github.com/LingFengJ/WASAText.git
cd WASAText
```

```bash
go mod tidy
go mod vendor
```

```bash
cd cmd/webapi
go build .
```

2. Make sure your project has the correct directory structure and all files are committed and pushed to GitHub:

WASAText/  
├── go.mod  
├── go.sum  
├── README.md  
├── service/  
│   ├── api/  
│   │   ├── reqcontext/  
│   │   │   └── reqcontext.go  
│   │   └── api.go  
│   ├── database/  
│   └── globaltime/  
└── cmd/  
└── webapi/   
└── main.go    
