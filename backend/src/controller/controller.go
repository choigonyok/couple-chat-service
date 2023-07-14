package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"choigonyok.com/couple-chat-service-project-docker/src/model"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Test(){
	// model.DeleteAll()

	usrsData := model.UsrsData{}
	usrsDatas := []model.UsrsData{}
	chatData := model.ChatData{}
	chatDatas := []model.ChatData{}
	requestData := model.RequestData{}
	requestDatas := []model.RequestData{}
	beAboutToDeleteData := model.BeAboutToDeleteData{}
	beAboutToDeleteDatas := []model.BeAboutToDeleteData{}
	connectionData := struct {
		connection_id int
		first_usr string
		second_usr string
		start_date string
	}{}
	connectionDatas := []struct {
		connection_id int
		first_usr string
		second_usr string
		start_date string
	}{}
	questionData := struct {
		question_id int
		target_word string
		question_contents string
	}{}
	questionDatas := []struct {
		question_id int
		target_word string
		question_contents string
	}{}
	answerData := struct {
		answer_id int
		connection_id int
		first_answer string
		second_answer string
		answer_date string
		question_id int
	}{}
	answerDatas := []struct {
		answer_id int
		connection_id int
		first_answer string
		second_answer string
		answer_date string
		question_id int
	}{}
	exceptionData := struct {
		exception_id int
		connection_id int
		except_word string
	}{}
	exceptionDatas := []struct {
		exception_id int
		connection_id int
		except_word string
	}{}

	r, _ := model.TestUsrs()
	for r.Next() {
		r.Scan(&usrsData.UUID, &usrsData.ID, &usrsData.Password, &usrsData.Conn_id, &usrsData.Order_usr)
		usrsDatas = append(usrsDatas, usrsData)
	}
	fmt.Println("usrs DB : ", usrsDatas)

	r, _ = model.TestChat()
	for r.Next(){
		r.Scan(&chatData.Chat_id, &chatData.Writer_id, &chatData.Write_time, &chatData.Text_body, &chatData.Is_answer)
		chatDatas = append(chatDatas, chatData)
	}
	fmt.Println("chat DB : ", chatDatas)

	r, _ = model.TestRequest()
	for r.Next(){
		r.Scan(&requestData.Request_id,&requestData.Requester_uuid,&requestData.Requester_id,&requestData.Target_uuid, &requestData.Target_id,&requestData.Request_time)
		requestDatas = append(requestDatas, requestData)
	}
	fmt.Println("request DB : ", requestDatas)

	r, _ = model.TestConnection()
	for r.Next() {
		r.Scan(&connectionData.connection_id, &connectionData.first_usr, &connectionData.second_usr,&connectionData.start_date)
		connectionDatas = append(connectionDatas, connectionData)
	}
	fmt.Println("connection DB : ", connectionDatas)

	r, _ = model.TestQuestion()
	for r.Next(){
		r.Scan(&questionData.question_id, &questionData.target_word, &questionData.question_contents)
		questionDatas = append(questionDatas, questionData)
	}
	fmt.Println("question DB : ", questionDatas)

	r, _ = model.TestAnswer()
	for r.Next() {
		r.Scan(&answerData.answer_id, &answerData.connection_id, &answerData.first_answer, &answerData.second_answer, &answerData.answer_date, &answerData.question_id)
		answerDatas = append(answerDatas, answerData)
	}
	fmt.Println("answer DB : ", answerDatas)

	r, _ = model.TestExceptionWord()
	for r.Next() {
		r.Scan(&exceptionData.exception_id, &exceptionData.connection_id, &exceptionData.except_word)
		exceptionDatas = append(exceptionDatas, exceptionData)
	}
	fmt.Println("exceptionword DB : ", exceptionDatas)

	r, _ = model.TestBeAboutToDelete()
	for r.Next() {
		r.Scan(&beAboutToDeleteData.Delete_Date, &beAboutToDeleteData.Connection_id)
		beAboutToDeleteDatas = append(beAboutToDeleteDatas, beAboutToDeleteData)
	}
	fmt.Println("beabouttodelete DB : ", beAboutToDeleteDatas)
}

var conns = make(map[string]*websocket.Conn)
var timerMap  = make(map[int]*time.Timer)
	
var mutex = &sync.Mutex{}

func ConnectDB(driverName, dbData string) {
	err := model.OpenDB(driverName, dbData)
	if err != nil {
		fmt.Println("ERROR #73 : ", err.Error())
	}
}

func UnConnectDB() {
	err := model.CloseDB()
	if err != nil {
		fmt.Println("ERROR #74 : ", err.Error())
	}
}

func getTimeNow() time.Time {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("ERROR #91 : ", err.Error())
	}
	now := time.Now()
	t := now.In(loc)
	return t
}


func LoadEnv(){
	// 환경변수 로딩
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ERROR #1 : ", err.Error())
	}
}

// ID, Password 유효성 검사
func checkIDandPWCorrect(ID string, PW string) bool {
	isIDCorrect, _ := regexp.MatchString("^[a-z][a-z0-9]+$",ID)
	isPWCorrect, _ := regexp.MatchString("^[a-z0-9]*$", PW)
	if len(ID) >= 21 {
		return false
	} else if len(PW) >= 21 {
		return false
	} else if !isIDCorrect {
		return false
	} else if !isPWCorrect {
		return false
	} else {
		return true
	}
}

// 회원가입	
func SignUpHandler(c *gin.Context) {

	signUpData := model.UsrsData{}
	err := c.ShouldBindJSON(&signUpData)
	if err != nil {
		fmt.Println("ERROR #6 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !checkIDandPWCorrect(signUpData.ID, signUpData.Password) {
		c.String(http.StatusBadRequest, "%v", "ID와 PW의 최대 길이는 20자로 제한됩니다. 또한 영어 소문자로 시작하는 영어소문자와 숫자의 조합만 유효합니다.")
		return
	}
	
	signUpData.UUID = uuid.New().String()

	err = model.InsertUsr(signUpData.ID, signUpData.Password, signUpData.UUID)
	if err != nil {
		fmt.Println("ERROR #9 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// 회원가입 시 아이디 중복체크
func IDCheckHandler(c *gin.Context){
	input := struct {
		ID string `json:"input_id"`
	}{}
	
	err := c.ShouldBindJSON(&input)
	if err != nil {
		fmt.Println("ERROR #4 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	isExist, err := model.CheckUsrByID(input.ID)
	if err != nil {
		fmt.Println("ERROR #5 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isExist == true {
		c.Writer.WriteHeader(http.StatusBadRequest)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

// 로그인
func LogInHandler(c *gin.Context){
	logInData := model.UsrsData{}
	err := c.ShouldBindJSON(&logInData)
	if err != nil {
		fmt.Println("ERROR #10 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	uuid, err := model.GetUUIDByIDandPW(logInData.ID, logInData.Password)
	if err != nil {
		fmt.Println("ERROR #11 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if uuid != "" {
		c.SetCookie("uuid", uuid, 60*60, "/", os.Getenv("ORIGIN"),false,true)
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
}

// 로그아웃
func LogOutHandler(c *gin.Context){
	LoadEnv()
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	c.SetCookie("uuid", uuid, -1, "/", os.Getenv("ORIGIN"), false, true)
	c.Writer.WriteHeader(http.StatusOK)
}

// 기존 로그인 되있던 상태인지 쿠키 확인	
func AlreadyLogInCheckHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn_id, err := model.SelectConnIDFromUsrsByUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #12 : ", err.Error())
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if conn_id == 0 {
		c.String(http.StatusOK, "%v", "NOT_CONNECTED")
	} else {
		c.String(http.StatusOK, "%v", "CONNECTED")
	}
}

func WithDrawalHandler(c *gin.Context){
	uuid, err1 := model.CookieExist(c)
	if err1 != nil {
		fmt.Println("ERROR #82 : ", err1.Error())
	}

	conn_id, err3 := model.SelectConnIDByUUID(uuid)
	if err3 != nil {
		fmt.Println("ERROR #84 : ", err3.Error())
	}
	
	if conn_id == 0 {
		err2 := model.DeleteUsrByUUID(uuid)
		if err2 != nil {
			fmt.Println("ERROR #83 : ", err2.Error())
		}
		c.SetCookie("uuid", "", -1, "/", os.Getenv("ORIGIN"), false, true)
		c.Writer.WriteHeader(http.StatusOK)
		return
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

// 상대방과의 커넥션 끊기
func CutConnectionHandler(c *gin.Context){
	uuid, err1 := model.CookieExist(c)
	if err1 != nil {
		fmt.Println("ERROR #85 : ", err1.Error())
	}

	conn_id, err2 := model.SelectConnIDByUUID(uuid)
	if err2 != nil {
		fmt.Println("ERROR #86 : ", err2.Error())
	}

	if timerMap[conn_id] != nil {
		c.String(http.StatusBadRequest, "%v", "ALREADY_REGISTER")
	} else {
		setConnDeleteTimer(uuid, conn_id)
		c.Writer.WriteHeader(http.StatusOK)
	}
}

// 커넥션 7일 후 종료를 위한 타이머 설정
func setConnDeleteTimer(uuid string, connection_id int) {
	setTime := getTimeNow().Add(3 * time.Minute)
	fmt.Println("SETTIME : ", setTime)
	fmt.Println("SETTIME : ", setTime)
	fmt.Println("SETTIME : ", setTime)
	timer := time.NewTimer(setTime.Sub(getTimeNow()))
	timerMap[connection_id] = timer
	go func() {
		<-timer.C
		fmt.Println("TIME END ! ")
		fmt.Println("TIME END ! ")
		fmt.Println("TIME END ! ")
		err := model.DeleteConnectionByConnID(uuid)
		if err != nil {
			fmt.Println("ERROR #90 : ", err.Error())
		}
		timerMap[connection_id] = nil
	}()
}

func RollBackConnectionHandler(c *gin.Context) {
	uuid, err1 := model.CookieExist(c)
	if err1 != nil {
		fmt.Println("ERROR #92 : ", err1.Error())
	}
	conn_id, err2 := model.SelectConnIDByUUID(uuid)
	if err2 != nil {
		fmt.Println("ERROR #93 : ", err2.Error())
	}
	if timerMap[conn_id] == nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	} else {
		timer := timerMap[conn_id]
		go func(){
			timer.Stop()
			<-timer.C
		}()
		timerMap[conn_id] = nil
		c.Writer.WriteHeader(http.StatusOK)
	}
}

// 상대방에게 connection 연결 요청
func ConnRequestHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	isExist, err := model.CheckRequestByRequesterUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #19 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isExist {
		c.String(http.StatusBadRequest, "%v", "ALREADY_REQUEST")
		return
	}

	id, err := model.SelectIDFromUsrsByUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #78 : ", err.Error())
	}

	input := struct {
		ID string `json:"input_id"`
	}{}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		fmt.Println("ERROR #20 : ", err.Error())
	}
	// 입력한 ID에 맞는 사용자 DATA DB에서 불러오기

	isExist, targetConnID, targetUUID, err := model.SelectConnIDandUUIDFromUsrsByID(input.ID)
	if err != nil {
		fmt.Println("ERROR #21 : ", err.Error())
	}
	if isExist {
		if targetUUID == uuid {
			c.String(http.StatusBadRequest, "%v", "NOT_YOURSELF")
		} else if targetConnID != 0 {
			c.String(http.StatusBadRequest, "%v", "ALREADY_CONNECTED")
		} else {
			// 요청된 정보를 DB에 저장
			err = model.InsertRequest(uuid, targetUUID, time.Now().Format("01/02 15:04"), id, input.ID)
			if err != nil {
				fmt.Println("ERROR #22 : ", err.Error())
			}
			c.Writer.WriteHeader(http.StatusOK)
		}
	} else {
	// ID가 존재하지 않는 ID면
		c.String(http.StatusBadRequest, "%v", "NOT_EXIST")
	}
}

// 현재 요청받은 request 목록 가져오기
func GetRecieveRequestHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestedDatas, err := model.SelectRecieveRequestByTargetUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #13 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshaledRequestedData, err := json.Marshal(requestedDatas)
	if err != nil {
		fmt.Println("ERROR #15 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.Write(marshaledRequestedData)
}


// 현재 신청중인 request 가져오기
func GetSendRequestHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestingData, err := model.SelectSendRequestByTargetUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #16 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshaledRequestingData, err := json.Marshal(requestingData)
	if err != nil {
		fmt.Println("ERROR #18 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	c.Writer.Write(marshaledRequestingData)
}

// 상대방과 연결 후, DB에 저장되어있던 자신과 상대 관련 요청 전체 삭제 + conn_id 생성
func DeleteRestRequestHandler(c *gin.Context){
	myUUID, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	target := struct {
		UUID string `json:"uuid_delete"`
	}{}

	err1 := c.ShouldBindJSON(&target)
	if err1 != nil {
		fmt.Println("ERROR #23 : ", err1.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err2 := model.InsertConnection(target.UUID, myUUID, time.Now().Format("2006/01/02"))
	if err2 != nil {
		fmt.Println("ERROR #24 : ", err2.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	connID, err3 := model.SelectConnectionIDByUsrsUUID(target.UUID, myUUID)
	if err3 != nil {
		fmt.Println("ERROR #25 : ", err3.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err4 := model.UpdateUsrsConnID(connID, target.UUID)
	if err4 != nil {
		fmt.Println("ERROR #26 : ", err4.Error())
		return
	}

	err5 := model.UpdateUsrsOrder(connID, myUUID)
	if err5 != nil {
		fmt.Println("ERROR #27 : ", err5.Error())
		return
	}

	err6 := model.DeleteRestRequest(target.UUID, myUUID)
	if err6 != nil {
		fmt.Println("ERROR #28 : ", err6.Error())
		return
	}

	err7 := model.InsertBeAboutToDelete(connID)
	if err7 != nil {
		fmt.Println("ERROR #24 : ", err7.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// 받은 요청 중 선택해서 요청을 삭제
func DeleteOneRequestHandler(c *gin.Context){
	request_id := c.Param("param")

	err := model.DeleteRequestByRequestID(request_id)
	if err != nil {
		fmt.Println("ERROR #29 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
}

// 그동안 답한 내용들을 모아서 보여주기 위한 API
func GetAnswerHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	connID, err := model.SelectConnIDByUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #30 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	answerDatas, err := model.GetAnswerandQuestionContentsByConnID(connID)
	if err != nil {
		fmt.Println("ERROR #31 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mashaledAnswerData, err := json.Marshal(answerDatas)
	if err != nil {
		fmt.Println("ERROR #33 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	c.Writer.Write(mashaledAnswerData)
}

// Websocket 프로토콜로 업그레이드 및 메시지 read/write
func UpgradeHandler(c *gin.Context){
	
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var upgrader  = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == os.Getenv("ORIGIN")
		    },
	}

	conn, err1 := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err1 != nil {
		fmt.Println("ERROR #34 : ", err1.Error())
		return
	}
	defer conn.Close()
	defer func(){
		mutex.Lock()
		conns[uuid] = nil
		mutex.Unlock()
	}()
	
	mutex.Lock()
	conns[uuid] = conn
	// conn객체를 읽어야함
	mutex.Unlock()
	
	// 클라이언트에 uuid 전달, 그래야 클라이언트에게 채팅을 표시할 때
	// 누가 보낸 채팅인지 UUID로 구분해서 표시할 수 있음
	jsonUUID := struct {
		UUID string `json:"uuid"`
	}{
		uuid,
	}
	err2 := conn.WriteJSON(jsonUUID)
	if err2 != nil {
		fmt.Println("ERROR #35 : ", err2.Error())
		return
	}

	first_uuid, second_uuid, conn_id, err3 := model.GetConnectionByUsrsUUID(uuid)
	if err3 != nil {
		fmt.Println("ERROR #36 : ", err3.Error())
		return
	}

	
	initialChats, err4 := model.SelectChatByUsrsUUID(first_uuid, second_uuid)
	if err4 != nil {
		fmt.Println("ERROR #37 : ", err4.Error())
		return 
	}
	
	err5 := conn.WriteJSON(initialChats)
	if err5 != nil {
		fmt.Println("ERROR #38 : ", err5.Error())
		return
	}

	// 이전에 대답 안하고 커넥션 종료된 question 있는지 확인
	order, err6 := model.GetUsrOrderByUUID(uuid)
	if err6 != nil {
		fmt.Println("ERROR #55 : ", err6.Error())
		return
	}

	question_id, err7 := model.QuestionIDOfEmptyAnswerByOrder(order, conn_id)
	if err7 != nil {
		fmt.Println("ERROR #79 : ", err7.Error())
		return
	}

	if question_id != 0 {
		_, questionContents, err8 := model.GetQuestionByQuestionID(question_id)
		if err8 != nil {
			fmt.Println("ERROR #80 : ", err8.Error())
			return
		}

		questiondata := model.ChatData{
			Text_body: questionContents,
			Writer_id: "question",
			Write_time: time.Now().Format("2006/01/02 03:04"),
			Is_answer: 1,
			Chat_id: 0,
			Question_id: question_id,
		}
		questiondatas := []model.ChatData{}
		questiondatas = append(questiondatas, questiondata)
		err := conn.WriteJSON(questiondatas)
		if err != nil {
			fmt.Println("ERROR #56 : ", err.Error())
			return
		}
	}

// 메시지를 읽고 쓰는 부분, 읽은 메시지는 DB에 저장됨
	for { 
		var chatData []model.ChatData

		err := conn.ReadJSON(&chatData)
		if err != nil {
			fmt.Println("ERROR #39 : ", err.Error())
			break;
		}

		// 일반채팅이면 chat table에 저장, question에 대한 답이면 answer table에 저장
		if chatData[0].Is_answer == 0 {
			err := model.InsertChat(chatData[0].Text_body, uuid, chatData[0].Write_time)
			// 어차피 커넥션 당 메시지 하나씩 전송 받으니까 slice index는 0으로 설정
			if err != nil {
				fmt.Println("ERROR #40 : ", err.Error())
			}	
		} else {
			recieveAnswer(uuid, conn_id, chatData, first_uuid)
		}

		target_conn := []*websocket.Conn{}

		mutex.Lock()
		if conns[first_uuid] != nil && conns[second_uuid] != nil {
			first_conn := conns[first_uuid]
			second_conn := conns[second_uuid]
			target_conn = append(target_conn, first_conn, second_conn)
		
		} else if conns[first_uuid] != nil {
			first_conn := conns[first_uuid]
			target_conn = append(target_conn, first_conn)
		
		} else {
			second_conn := conns[second_uuid]
			target_conn = append(target_conn, second_conn)	
		}
		mutex.Unlock()



		// 커넥션 연결이 안되어있으면 보내면 nil pointer 오류 생김
		// 모든 커넥션에 메시지 write
		
		if chatData[0].Is_answer != 1 {
			for _, item := range target_conn {
				err := item.WriteJSON(chatData)
				if err != nil {
					fmt.Println("ERROR #43 : ", err.Error())
				}
			}
		}

		sendQuestion(chatData, conn_id, target_conn)
	}
}

func sendQuestion(chatData []model.ChatData, conn_id int, target_conn []*websocket.Conn){
// 채팅 중 단어가 발견되면 단어 관련된 질문을 커플에게 던지는 기능
	// 1. 단어를 먼저 다 뽑아서
	var target_word, question_contents string
	var question_id int
	r, err := model.SelectQuetions()
	if err != nil {
		fmt.Println("ERROR #44 : ", err.Error())
	}
	for r.Next() {
		// 2. 방금 READ한 채팅에 단어가 있는지 돌면서 확인
		r.Scan(&target_word, &question_id, &question_contents)	
		if strings.Contains(chatData[0].Text_body, target_word) {
			// 3. 단어가 발견되면 이전에 답을 한 전적이 있는지 검색
			isExist, err := model.CheckAnswerByConnIDandQuestionID(conn_id, question_id)
			if err != nil {
				fmt.Println("ERROR #45 : ", err.Error())
				return
			}

			// 4. 단어도 발견됐고, 이전에 했던 질문도 아니면 질문 WRITE
			if !isExist {
				questiondata := model.ChatData{
					Text_body: question_contents,
					Writer_id: "question",
					Write_time: time.Now().Format("2006/01/02 03:04"),
					Is_answer: 1,
					Chat_id: 0,
					Question_id: question_id,
				}
				questiondatas := []model.ChatData{}
				questiondatas = append(questiondatas, questiondata)

				for _, item := range target_conn {
					err := item.WriteJSON(questiondatas)
					if err != nil {
						fmt.Println("ERROR #46 : ", err.Error())
					}
				}
				// 5. answer에 답 적기 (는 위에 READ에서 처리)
				err = model.InsertAnswer(chatData[0].Write_time, conn_id, question_id)
				if err != nil {
					fmt.Println("ERROR #42 : ", err.Error())
				}
			}
		}
	}
}

func recieveAnswer(uuid string, conn_id int, chatData []model.ChatData, first_uuid string){
	isExist, err1 := model.CheckAnswerByConnIDandQuestionID(conn_id, chatData[0].Question_id)
	if err1 != nil {
		fmt.Println("ERROR #41 : ", err1.Error())
		return
	}
	
	if isExist {
		var err2 error
		if first_uuid == uuid {
			err2 = model.UpdateFirstAnswerByQuestionID(chatData[0].Text_body, chatData[0].Question_id)
		} else {
			err2 = model.UpdateSecondAnswerByQuestionID(chatData[0].Text_body, chatData[0].Question_id)
		}
		if err2 != nil {
			fmt.Println("ERROR #50 : ", err2.Error())
		}
	}
}

func GetMostUsedWordsHandler(c *gin.Context){
	rankNumString := c.Param("ranknum")
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	rankNumInt, err := strconv.Atoi(rankNumString)
	if err != nil {
		fmt.Println("ERROR #59 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return	
	}
	firstUUID, secontUUID, _, err := model.GetConnectionByUsrsUUID(uuid)
	
	var ohterFrequentWords []string
	var err2 error
	if firstUUID == uuid {
		ohterFrequentWords, err2 = model.GetFrequentWords(secontUUID, rankNumInt)	
		if err2 != nil {
			fmt.Println("ERROR #80 : ", err2.Error())
		}
	} else {
		ohterFrequentWords, err2 = model.GetFrequentWords(firstUUID, rankNumInt)
		if err2 != nil {
			fmt.Println("ERROR #80 : ", err2.Error())
		}
	}
	myFrequentWords, err3 := model.GetFrequentWords(uuid, rankNumInt)
	if err3 != nil {
		fmt.Println("ERROR #80 : ", err3.Error())
	}

	if myFrequentWords == nil || ohterFrequentWords == nil {
		c.Writer.WriteHeader(http.StatusLengthRequired)
		return
	}

	sendData := struct {
		MyWords []string `json:"mywords"`
		OtherWords []string `json:"otherwords"`
	}{
		myFrequentWords,
		ohterFrequentWords,
	}

	marshaledData, err := json.Marshal(sendData)
	if err != nil {
		fmt.Println("ERROR #60 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return	
	} 

	c.Writer.Write(marshaledData)
}

func GetExceptWordsHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	conn_id, err1 := model.SelectConnIDByUUID(uuid)
	if err1 != nil {
		fmt.Println("ERROR #61 : ", err1.Error())
	}
	
	exceptWords, err2 := model.GetExceptWords(conn_id)
	if err2 != nil {
		fmt.Println("ERROR #62 : ", err2.Error())
	}
	if len(exceptWords) == 0 {
		c.Writer.WriteHeader(http.StatusNoContent)	
		return
	}
	marshaledData, err3 := json.Marshal(exceptWords)
	if err3 != nil {
		fmt.Println("ERROR #81 : ", err3.Error())
	}
	c.Writer.Write(marshaledData)
}

func InsertExceptWordHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	conn_id, err := model.SelectConnIDByUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #66 : ", err.Error())
	}

	Input := struct {
		Except_word string `json:"except_word"`
	}{}
	err1 := c.ShouldBindJSON(&Input)
	if err1 != nil {
		fmt.Println("ERROR #63 : ", err1.Error())
	}
	isExitst, err2 := model.CheckWordAlreadyExcepted(conn_id, Input.Except_word)
	if err2 != nil {
		fmt.Println("ERROR #63 : ", err2.Error())
	}
	if isExitst {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	} else {
		err2 := model.InsertExceptWord(conn_id, Input.Except_word)
		if err2 != nil {
		fmt.Println("ERROR #63 : ", err2.Error())
		}
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func DeleteExceptWordHandler(c *gin.Context){
	uuid, err := model.CookieExist(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	param := c.Param("param")

	conn_id, err2 := model.SelectConnIDByUUID(uuid)
	if err2 != nil {
		fmt.Println("ERROR #67 : ", err2.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err3 := model.DeleteExceptWord(conn_id, param)
	if err3 != nil {
		fmt.Println("ERROR #68 : ", err3.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}