package main

import (
	"fmt"
	"os"

	"choigonyok.com/couple-chat-service-project-docker/src/controller"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func LoadEnv(){
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ERROR #1 : ", err.Error())
	}
}

func originConfig() cors.Config{
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("ORIGIN")} 
	// 허용할 오리진 설정, 원래 리액트의 port가 아니라 리액트가 있는 container의 port 번호를 origin allow 해줘야함
	// localhost:3000로 origin allow 하면 통신 안됨

	config.AllowMethods= []string{"GET", "POST", "DELETE", "PUT"}
	config.AllowHeaders = []string{"Content-type"}
	config.AllowCredentials = true
	return config
}

func main() {
	LoadEnv()	// 환경변수 로딩
	
	e := gin.Default()
	
	config := originConfig()	// Origin 설정
	e.Use(cors.New(config)) 	// Origin 적용
	
	controller.ConnectDB("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp(mysql)/"+os.Getenv("DB_NAME"))	// DB 초기 연결
	defer controller.UnConnectDB()
	
	e.POST("/api/usr", controller.SignUpHandler)								// 회원가입
	e.DELETE("/api/usr", controller.WithDrawalHandler)							// 회원탈퇴
	e.PUT("/api/usr", controller.ChangePasswordHandler)							// 비밀번호 변경
	
	e.POST("/api/id", controller.IDCheckHandler)								// 회원가입 시 아이디 중복체크

	e.DELETE("/api/conn", controller.CutConnectionHandler)						// 커넥션 끊기
	e.PUT("/api/conn", controller.RollBackConnectionHandler)					// 커넥션 재연결

	e.POST("/api/log", controller.LogInHandler)									// 로그인
	e.DELETE("/api/log", controller.LogOutHandler)								// 로그아웃
	e.GET("/api/log", controller.AlreadyLogInCheckHandler)						// 기존 로그인 되있던 상태인지 쿠키 확인

	e.POST("/api/request", controller.ConnRequestHandler)						// 상대방에게 connection 연결 요청	
	e.GET("/api/request/recieved", controller.GetRecieveRequestHandler)			// 현재 요청받은 request 목록 가져오기
	e.GET("/api/request/send", controller.GetSendRequestHandler)				// 현재 신청중인 request 가져오기
	e.PUT("/api/request", controller.DeleteRestRequestHandler)					// 상대방과 연결 후, DB에 저장되어있던 자신과 상대 관련 요청 전체 삭제 + conn_id 생성
	e.DELETE("/api/request/:param", controller.DeleteOneRequestHandler)			// 받은 요청 중 선택해서 요청을 삭제

	e.GET("/api/answer", controller.GetAnswerHandler)							// 그동안 답한 내용들을 모아서 보여주기 위한 API

	e.POST("/api/file", controller.InsertFileHandler)							// 채팅으로 보낸 파일 서버에 저장
	e.GET("/api/file/:chatID", controller.GetFileHandler)						// chatpage 렌더링용 썸네일 이미지 불러오기
	e.GET("/api/file/name/:chatID", controller.GetFileNameHandler)				// 파일이름+확장자 찾기

	e.GET("/api/chat/word/:param", controller.GetChatWordHandler)				// 단어 기반 채팅 검색
	e.GET("/api/chat/date", controller.GetChatDateHandler)						// 날짜 기반 채팅 검색

	e.POST("/api/anniversary", controller.InsertAnniversaryHandler)				// 일정, 기념일 추가
	e.GET("/api/anniversary", controller.GetAnniversaryHandler)					// 일정, 기념일 불러오기
	e.GET("/api/anniversary/dday", controller.GetDDayHandler)					// D-DAY 불러오기
	e.DELETE("/api/anniversary/:id", controller.DeleteAnniversaryHandler)		// 일정, 기념일 삭제

	e.GET("/api/rank/:ranknum", controller.GetMostUsedWordsHandler)				// 사용자가 가장 많이 사용한 단어 랭킹 보여주기

	e.GET("/ws", controller.UpgradeHandler)										// Websocket 프로토콜로 업그레이드 및 메시지 read/write

	e.GET("/api/except", controller.GetExceptWordsHandler)						// FrequentUsedWords에서 제외된 단어 불러오기
	e.POST("/api/except", controller.InsertExceptWordHandler)					// FrequentUsedWords에서 제외할 단어 입력받기
	e.DELETE("/api/except/:param", controller.DeleteExceptWordHandler)			// FrequentUsedWords에서 단어 제외 취소하기

	controller.Test() // DB 저장 데이터 출력 TEST
	
	e.Run(":8080")
}
