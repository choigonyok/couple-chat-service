# 서버 도커파일 작성

FROM golang:1.18.1-buster

WORKDIR /app

COPY go.mod go.sum ./
* go.mod를 따로 먼저 카피해서 종속성을 설치해줘야 함
* 안그러면 나중에 종속성과 관련 없는 부분이 수정되어도 빌드하면 종속성까지 다 다시 다운 받게 됨
* go.sum도 여기 넣어줘야함 이거 안넣어주면 컨테이너 빌드할 때 오류남

RUN go mod download
* go mod tidy 같은 느낌

COPY . .
* go mod, go.sum 제외한 파일 모두 복사

RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]

* RUN go build -o /main 대신 go build main.go 로 빌드 가능
* 이렇게 하면 CMD ["/main"] 대신 CMD ["./main] 써야함
* 만약에 도커파일과 main.go 파일의 디렉토리가 다르면 RUN cd ~ 로 이동해서 빌드하도록 해줘야함
* RUN은 명령어 치는 거

* 도커파일을 이미지로 빌드할 때는 docker build -t 이름 . 이 맨 뒤에 .이 꼭 붙어야함
* dot은 dockerfile을 어디서 찾을지를 알게해주는 거
* 만약 지금 디렉토리에 dockerfile이 없으면 dockerfile이 위치한 디렉토리를 알려줘야함

#docker build -t 이미지이름 ->으로 이미지 빌드하고
#docker run -d -p 8080:8080 이미지이름 -> 으로 해당 이미지로 컨테이너를 만들고 실행

# 서버 Origin Allow 설정 FOR CORS

```go
config := cors.DefaultConfig()
config.AllowOrigins = []string{"http://localhost:2000"} 
// 허용할 오리진 설정, 원래 리액트의 port가 아니라 리액트가 있는 container의 port 번호를 origin allow 해줘야함
// localhost:3000로 origin allow 하면 통신 안됨

config.AllowMethods= []string{"GET"}
config.AllowHeaders = []string{"Content-type"}
config.AllowCredentials = true
eg.Use(cors.New(config)) 
// origin 설정하고 설정한 config를 gin engine에서 사용하겠다는 이 부분이 있어야 적용이 됨!
```

# 클라이언트 도커파일 설정

FROM node:16-alpine AS build
* 이미지를 node:10으로 당겨왔다가 1시간동안 뻘짓함
* 도커 이미지 빌드 중 npm run build에서 막히길래 한참 찾았음

WORKDIR /app

COPY ./package.json ./

RUN npm install
* package.json을 해야 코드의 일부를 수정할 때 매번 종속성까지 다시 다운받지 않을 수 있음

COPY . .

RUN npm run build
* axios 등 다른 패키지를 사용하고 도커 이미지를 만드려면 package.json에 해당 패키지를 작성해줘야 함

FROM nginx

EXPOSE 3000

COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf

COPY --from=build /app/build /usr/share/nginx/html 

# 클라이언트 API 요청 작성

```js
axios
  .get("http://localhost:1000/api/test")
  // api호출은 go port num인 8080이 아니라 container port num인 1000으로 요청해야 통신이 됨
  // localhost:8080으로 요청하면 통신 안됨
```

# 클라이언트 nginx 루트 설정

```nginx
server {
        listen 3000;

        location / {
                root /usr/share/nginx/html;
                # 정적 파일을 브라우저에 제공하기 위해. 리액트에서 빌드한 파일이 어디있는지 지정하는 거
                index index.html index.htm;
                # 처음 시작을 뭘로 할 건지 설정
                try_files $uri $uri/ /index.html;
                # SPA만 만들 수 있는 리액트에서 라우팅을 가능하게 해줌
        }
}
```