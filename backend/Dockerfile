FROM golang:latest

WORKDIR /app

COPY . .
# go.mod를 따로 먼저 카피해서 종속성을 설치해줘야 함
# 안그러면 나중에 종속성과 관련 없는 부분이 수정되어도 빌드하면 종속성까지 다 다시 다운 받게 됨
# go.sum도 여기 넣어줘야함 이거 안넣어주면 컨테이너 빌드할 때 오류남
# 로컬에서 Dockerfile.dev로 개발하면서 필요한 종속성은 다 go.mod, go.sum 안에 넣어두기 때문에
# 도커파일에서 go get github.com/gin-gonic/gin 같은 명령어 작성 안해줘도 됨
# 근데 저 CompileDaemon은 왜 작성해야하는거? 거참

RUN go mod download
# go mod tidy 같은 느낌

# go mod, go.sum 제외한 파일 모두 복사

RUN cd src && go build main.go

EXPOSE 8080

CMD ["./src/main" ]

# RUN go build -o /main 대신 go build main.go 로 빌드 가능
# 이렇게 하면 CMD ["/main"] 대신 CMD ["./main] 써야함
# 만약에 도커파일과 main.go 파일의 디렉토리가 다르면 RUN cd ~ 로 이동해서 빌드하도록 해줘야함
# RUN은 명령어 치는 거

# 도커파일을 이미지로 빌드할 때는 docker build -t 이름 . 이 맨 뒤에 .이 꼭 붙어야함
# dot은 dockerfile을 어디서 찾을지를 알게해주는 거
# 만약 지금 디렉토리에 dockerfile이 없으면 dockerfile이 위치한 디렉토리를 알려줘야함

#docker build -t 이미지이름 ->으로 이미지 빌드하고
#docker run -d -p 8080:8080 이미지이름 -> 으로 해당 이미지로 컨테이너를 만들고 실행

# RUN cd src
# go build main.go
# EXPOSE 8080
# CMD ["./main" ]
# 이렇게 하면 src 안에 있는 main.go가 빌드가 안됨
# 로컬에서는 잘 되는데 왜인지 찾아보니
# 도커의 RUN command는 이전 RUN을 기억하지 않는다는거!
# 그래서 아래처럼 바꿔줌

# RUN cd src && go build main.go
# EXPOSE 8080
# CMD ["./main" ]
# 그랬더니 빌드는 됐는데, CMD에서 main 빌드 파일을 실행하려고 하니까
# 또 directory가 초기화돼서 main파일을 찾을 수 없다고 함
# 그래서 아래처럼 바꿔줌

# RUN cd src && go build main.go
# EXPOSE 8080
# CMD ["cd", "src", "./main" ]
# 이거도 안됨. 한 줄에 디렉토리 변경과 빌드파일 실행을 한 줄의 명령어로 할 수는 없는 듯힘
# 그래서 아래처럼 바꿔줌

# RUN cd src && go build main.go
# EXPOSE 8080
# WORKDIR /app/backend/src
# CMD ["./main" ]
# 이것도 빌드는 되는데 이미지로 컨테이너 실행하면 ./main을 찾을 수 없다고 함
# 그래서 아래처럼 바꿔줌

# RUN cd src && go build main.go
# EXPOSE 8080
# WORKDIR /app/src
# CMD ["./main" ]
# 이랬더니 됨!!!!!!!!!!!!!!!
# 도커파일에서 처음에 WORKDIR을 /app으로 설정해주고 COPY . . 을 할 때
# DOCKERFILE이 있는 디렉토리의 모든 파일과 폴더가 app으로 복사 되었고,
# 그래서 backend폴더는 컨테이너에는 복사되지 않고 app 다음 바로 src 폴더가 위치하게 된 것
# 그리고 WORKDIR을 /app/src로 설정해주면서 CMD가 제 위치에서 실행될 수 있었던 것!