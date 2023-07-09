package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"choigonyok.com/couple-chat-service-project-docker/src/model"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Test(){
	//model.DeleteAll()

	usrsData := model.UsrsData{}
	usrsDatas := []model.UsrsData{}
	chatData := model.ChatData{}
	chatDatas := []model.ChatData{}
	requestData := model.RequestData{}
	requestDatas := []model.RequestData{}
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
	fmt.Println(connectionDatas)

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
}

var conns = make(map[string]*websocket.Conn)

func ConnectDB(driverName, dbData string) {
	model.OpenDB(driverName, dbData)
}

func UnConnectDB() {
	model.CloseDB()
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

// usr가 상대방과 연결된 상태인지 아닌지 체크
func isConnected(uuid string) bool {
	r, err := model.SelectConnIDFromUsrsByUUID(uuid)
	defer r.Close()
	if err != nil {
		fmt.Println("ERROR #49 : ", err.Error())
	}
	var conn_id int
	for r.Next() {
		r.Scan(&conn_id)
	}
	if conn_id == 0 {
		return false
	} else {
		return true
	}
}

// 쿠키가 있는지 확인
func cookieExist(c *gin.Context) string {
	uuid, err := c.Cookie("uuid")	
	if err != nil {
		fmt.Println("ERROR #14 : ", err.Error())
		c.Writer.WriteHeader(http.StatusUnauthorized)
	}
	return uuid
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

	r, err := model.SelectUsrByID(input.ID)
	if err != nil {
		fmt.Println("ERROR #5 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Close()

	if r.Next() {
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
	r, err := model.SelectUUIDFromUsrsByIDandPW(logInData.ID, logInData.Password)
	if err != nil {
		fmt.Println("ERROR #11 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Close()

	if r.Next() {
		var cookieValue string
		r.Scan(&cookieValue)
		c.SetCookie("uuid", cookieValue, 60*60, "/", os.Getenv("ORIGIN"),false,true)
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
}

// 로그아웃
func LogOutHandler(c *gin.Context){
	LoadEnv()
	uuid := cookieExist(c)
	c.SetCookie("uuid", uuid, -1, "/", os.Getenv("ORIGIN"), false, true)
	c.Writer.WriteHeader(http.StatusOK)
}

// 기존 로그인 되있던 상태인지 쿠키 확인	
func AlreadyLogInCheckHandler(c *gin.Context){
	uuid := cookieExist(c)

	r, err := model.SelectUsrByUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #12 : ", err.Error())
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer r.Close()

	if r.Next() {
		if isConnected(uuid) {
			c.String(http.StatusOK, "%v", "CONNECTED")
		} else {
			c.String(http.StatusOK, "%v", "NOT_CONNECTED")
		}	
	}
}

// 상대방에게 connection 연결 요청	
func ConnRequestHandler(c *gin.Context){
	uuid := cookieExist(c)

	r1, err := model.SelectRequestByRequesterUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #19 : ", err.Error())
	}
	defer r1.Close()
	
	if r1.Next() {
		c.String(http.StatusBadRequest, "%v", "ALREADY_REQUEST")
		return
	}

	var id string
	r2, _ := model.SelectIDFromUsrsByUUID(uuid)
	r2.Next()
	r2.Scan(&id)

	input := struct {
		ID string `json:"input_id"`
	}{}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		fmt.Println("ERROR #20 : ", err.Error())
	}
	// 입력한 ID에 맞는 사용자 DATA DB에서 불러오기

	r3, err := model.SelectConnIDandUUIDFromUsrsByID(input.ID)
	if err != nil {
		fmt.Println("ERROR #21 : ", err.Error())
	}
	defer r3.Close()
	// ID가 존재하는 ID면 이미 연결되어있진 않은지 conn_id를 확인
	if r3.Next() {
		var targetConnID int
		var targetUUID string
		err := r3.Scan(&targetConnID, &targetUUID)

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
	uuid := cookieExist(c)

	requestingData := model.RequestData{}
	requestingDatas := []model.RequestData{}

	r, err := model.SelectRecieveRequestByTargetUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #13 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Close()

	for r.Next() {
		r.Scan(&requestingData.Requester_id, &requestingData.Requester_uuid, &requestingData.Request_time, &requestingData.Request_id)
		requestingDatas = append(requestingDatas, requestingData)
	}
	marshaledRequestedData, err := json.Marshal(requestingDatas)
	if err != nil {
		fmt.Println("ERROR #15 : ", err.Error())
	}

	c.Writer.Write(marshaledRequestedData)
}


// 현재 신청중인 request 가져오기
func GetSendRequestHandler(c *gin.Context){
	uuid := cookieExist(c)

	r, err := model.SelectSendRequestByTargetUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #16 : ", err.Error())
	}
	defer r.Close()

	requestingData := model.RequestData{}
	for r.Next(){
		r.Scan(&requestingData.Target_uuid, &requestingData.Request_time, &requestingData.Target_id)
	}

	marshaledRequestingData, err := json.Marshal(requestingData)
	if err != nil {
		fmt.Println("ERROR #18 : ", err.Error())
	}
	
	c.Writer.Write(marshaledRequestingData)
}

// 상대방과 연결 후, DB에 저장되어있던 자신과 상대 관련 요청 전체 삭제 + conn_id 생성
func DeleteRestRequestHandler(c *gin.Context){
	myUUID  := cookieExist(c)
	
	target := struct {
		UUID string `json:"uuid_delete"`
	}{}

	err := c.ShouldBindJSON(&target)
	if err != nil {
		fmt.Println("ERROR #23 : ", err.Error())
		return
	}

	_, err = model.InsertConnection(target.UUID, myUUID, time.Now().Format("2006/01/02"))
	if err != nil {
		fmt.Println("ERROR #24 : ", err.Error())
	}

	r, err := model.SelectConnectionIDByUsrsUUID(target.UUID, myUUID)
	if err != nil {
		fmt.Println("ERROR #25 : ", err.Error())
	}
	defer r.Close()

	r.Next()
	var connID int
	r.Scan(&connID)

	_, err = model.UpdateUsrsConnID(connID, target.UUID)
	if err != nil {
		fmt.Println("ERROR #26 : ", err.Error())
		return
	}

	_, err = model.UpdateUsrsOrder(connID, myUUID)
	if err != nil {
		fmt.Println("ERROR #27 : ", err.Error())
		return
	}

	_, err = model.DeleteRestRequest(target.UUID, myUUID)
	if err != nil {
		fmt.Println("ERROR #28 : ", err.Error())
		return
	}
}

// 받은 요청 중 선택해서 요청을 삭제
func DeleteOneRequestHandler(c *gin.Context){
	request_id := c.Param("param")

	err := model.DeleteRequestByRequestID(request_id)
	if err != nil {
		fmt.Println("ERROR #29 : ", err.Error())
	}
}

// 그동안 답한 내용들을 모아서 보여주기 위한 API
func GetAnswerHandler(c *gin.Context){
	uuid := cookieExist(c)

	r1, err := model.SelectConnIDByUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #30 : ", err.Error())
	}
	defer r1.Close()

	r1.Next()
	var connID string
	r1.Scan(&connID)

	type AnswerData struct {
		QuestionContents string `json:"question_contents"`
		FirstAnswer string `json:"first_answer"`
		SecondAnswer string `json:"second_answer"`
		AnswerDate string `json:"answer_date"`
	}

	answerData := AnswerData{}
	answerDatas := []AnswerData{}

	var questionID string
	r2, err := model.SelectAnswerByConnID(connID)
	if err != nil {
		fmt.Println("ERROR #31 : ", err.Error())
	}
	defer r2.Close()
	for r2.Next() {
		r2.Scan(&answerData.FirstAnswer, &answerData.SecondAnswer, &answerData.AnswerDate, &questionID)
		r3, err := model.SelectQuestionContentsByQuestionID(questionID)
		if err != nil {
			fmt.Println("ERROR #32 : ", err.Error())
		}
		defer r3.Close()

		r3.Next()
		r3.Scan(&answerData.QuestionContents)

		if answerData.FirstAnswer == "not-written" || answerData.SecondAnswer == "not-written" {
			continue
		}

		answerDatas = append(answerDatas, answerData)
	}

	mashaledAnswerData, err := json.Marshal(answerDatas)
	if err != nil {
		fmt.Println("ERROR #33 : ", err.Error())
	}
	
	c.Writer.Write(mashaledAnswerData)
}

// Websocket 프로토콜로 업그레이드 및 메시지 read/write
func UpgradeHandler(c *gin.Context){
	uuid := cookieExist(c)

	var upgrader  = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == os.Getenv("ORIGIN")
		    },
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("ERROR #34 : ", err.Error())
		return
	}
	defer conn.Close()
	defer func(){
		conns[uuid] = nil
	}()

	conns[uuid] = conn
	// conn객체를 읽어야함
	
	// 클라이언트에 uuid 전달, 그래야 클라이언트에게 채팅을 표시할 때
	// 누가 보낸 채팅인지 UUID로 구분해서 표시할 수 있음
	jsonUUID := struct {
		UUID string `json:"uuid"`
	}{
		uuid,
	}
	err = conn.WriteJSON(jsonUUID)
	if err != nil {
		fmt.Println("ERROR #35 : ", err.Error())
		return
	}

	r1, err := model.SelectConnectionByUsrsUUID(uuid)
	if err != nil {
		fmt.Println("ERROR #36 : ", err.Error())
	}
	defer r1.Close()

	r1.Next()
	var first_uuid, second_uuid string
	var conn_id int
	r1.Scan(&first_uuid, &second_uuid, &conn_id)

	// 기존 저장되어있던 채팅 DB에서 불러와서 표시
	initialChat := model.ChatData{}
	initialChats := []model.ChatData{}
	r2, err := model.SelectChatByUsrsUUID(first_uuid, second_uuid)
	// LATER : 나중에 여러 conn 구현하면 쿼리문에 조건절이랑 conn_id 넣기
	if err != nil {
		fmt.Println("ERROR #37 : ", err.Error())
	}
	for r2.Next() {
		r2.Scan(&initialChat.Chat_id, &initialChat.Writer_id, &initialChat.Write_time, &initialChat.Text_body)
		initialChats = append(initialChats, initialChat)
	}
	err = conn.WriteJSON(initialChats)
	if err != nil {
		fmt.Println("ERROR #38 : ", err.Error())
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
			r, err := model.SelectAnswerByConnIDandQuestionID(conn_id, question_id)
			if err != nil {
				fmt.Println("ERROR #45 : ", err.Error())
			}
			defer r.Close()
			// 4. 단어도 발견됐고, 이전에 했던 질문도 아니면 질문 WRITE
			if !r.Next() {
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
			}
		}
	}
}

func recieveAnswer(uuid string, conn_id int, chatData []model.ChatData, first_uuid string){
	r3, err := model.SelectAnswerByConnIDandQuestionID(conn_id, chatData[0].Question_id)
	if err != nil {
		fmt.Println("ERROR #41 : ", err.Error())
	}
	defer r3.Close()

	if r3.Next() {
		if first_uuid == uuid {
			err = model.UpdateFirstAnswerByQuestionID(chatData[0].Text_body, chatData[0].Question_id)
		} else {
			err = model.UpdateSecondAnswerByQuestionID(chatData[0].Text_body, chatData[0].Question_id)
		}
	} else {
		err = model.InsertAnswer(chatData[0].Write_time, conn_id, chatData[0].Question_id)
		if err != nil {
			fmt.Println("ERROR #42 : ", err.Error())
		}
		if first_uuid == uuid {
			err = model.UpdateFirstAnswerByQuestionID(chatData[0].Text_body, chatData[0].Question_id)
		} else {
			err = model.UpdateSecondAnswerByQuestionID(chatData[0].Text_body, chatData[0].Question_id)
		}
		if err != nil {
			fmt.Println("ERROR #50 : ", err.Error())
		}
	}
}