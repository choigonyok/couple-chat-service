package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type ChatData struct {
	Text_body string `json:"text_body"`
        Writer_id string `json:"writer_id"`
        Write_time string `json:"write_time"`
	Is_answer int `json:"is_answer"`
	Chat_id int `json:"chat_id"`
	Question_id int `json:"question_id"`
}

type RequestData struct {
	Request_id int
	Requester_uuid string
	Requester_id string
	Target_uuid string
	Target_id string
	Request_time string
	
}

type UsrsData struct {
	ID string `json:"usr_id"`
	Password string `json:"usr_pw"`
	UUID string
	Conn_id int
	Order_usr int
 }

// Origin CORS 설정
func originConfig() cors.Config{
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("ORIGIN")} 
	// 허용할 오리진 설정, 원래 리액트의 port가 아니라 리액트가 있는 container의 port 번호를 origin allow 해줘야함
	// localhost:3000로 origin allow 하면 통신 안됨

	config.AllowMethods= []string{"GET"}
	config.AllowHeaders = []string{"Content-type"}
	config.AllowCredentials = true
	return config
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
func isConnected(c *gin.Context, db *sql.DB) bool {
	uuid, err := c.Cookie("uuid")
	fmt.Println(uuid)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("COOKIE LOADING TO CHECK CONNECTED ERROR OCCURED")
	}
	r, err := db.Query(`SELECT conn_id FROM usrs WHERE uuid = "`+uuid+`"`)
	defer r.Close()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("LOAD DB TO CHECK CONNECTED ERROR OCCURED")
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

func main() {	
// 커넥션 집합 슬라이스
	conns := make(map[string]*websocket.Conn)

// 환경변수 로딩
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ERROR #1 : ", err.Error())
	}

	ginEngine := gin.Default()

	config := originConfig()
	ginEngine.Use(cors.New(config)) 
	
	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp(mysql)/"+os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Println("ERROR #2 : ", err.Error())
	}
	defer db.Close()

// DB와 서버가 연결 되었는지 확인
	err = db.Ping()
	if err != nil {
		fmt.Println("ERROR #3 : ", err.Error())
	}

// TEST : DB에 usr 정보가 잘 저장되는지 테스트
	test := UsrsData{}
	tests := []UsrsData{}
	r, err := db.Query("SELECT id, password, uuid, conn_id, order_usr FROM usrs")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("USER TEST ERROR")
	}
	defer r.Close()
	for r.Next()  {
		r.Scan(&test.ID, &test.Password, &test.UUID, &test.Conn_id, &test.Order_usr)
		tests = append(tests, test)
	}
	fmt.Println("NOW STORED USR ID / PW / UUID / CONN_ID / ORDER : ", tests)

// TEST : DB에 chat data가 잘 저장되는지 테스트
	chattest := ChatData{}
	chattests := []ChatData{} 
	r, err = db.Query("SELECT chat_id, writer_id, write_time, text_body FROM chat")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("CHAT TEST ERROR")
	}
	defer r.Close()
	for r.Next() {
		r.Scan(&chattest.Chat_id,&chattest.Writer_id,&chattest.Write_time,&chattest.Text_body)
		chattests = append(chattests, chattest)
	}
	fmt.Println("NOW STORED CHAT : ", chattests)

// TEST : DB에 request data가 잘 저장되는지 테스트
	requesttest := RequestData{}
	requesttests := []RequestData{}
	r, err = db.Query("SELECT request_id, requester_uuid, requester_id, target_uuid, target_id, request_time FROM request")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("REQUEST TEST ERROR")
	}
	defer r.Close()
	for r.Next() {
		r.Scan(&requesttest.Request_id, &requesttest.Requester_uuid,&requesttest.Requester_id, &requesttest.Target_uuid,&requesttest.Target_id, &requesttest.Request_time)
		requesttests = append(requesttests, requesttest)
	}
	fmt.Println("NOW STORED REQUEST : ", requesttests)

// TEST : DB에 connection data가 잘 저장되는지 테스트
	connectiontest := struct {
		connection_id int
		first_usr string
		second_usr string
		start_date string
	}{}
	connectiontests := []struct {
		connection_id int
		first_usr string
		second_usr string
		start_date string
	}{}
	r, err = db.Query("SELECT connection_id, first_usr, second_usr, start_date FROM connection")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("CONNECTION TEST ERROR")
	}
	defer r.Close()
	for r.Next() {
		r.Scan(&connectiontest.connection_id, &connectiontest.first_usr, &connectiontest.second_usr, &connectiontest.start_date)
		connectiontests = append(connectiontests, connectiontest)
	}
	fmt.Println("NOW STORED CONNECTION : ", connectiontests)	

// TEST : DB에 connection data가 잘 저장되는지 테스트
	questiontest := struct {
		question_id int
		target_word string
		question_contents string
	}{}
	questiontests := []struct {
		question_id int
		target_word string
		question_contents string
	}{}
	r, err = db.Query("SELECT question_id, target_word, question_contents FROM question")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("QUESTION TEST ERROR")
	}
	defer r.Close()
	for r.Next() {
		r.Scan(&questiontest.question_id, &questiontest.target_word, &questiontest.question_contents)
		questiontests = append(questiontests, questiontest)
	}
	fmt.Println("NOW STORED QUESTION : ", questiontests)	

// TEST : DB에 connection data가 잘 저장되는지 테스트
	answertest := struct {
		answer_id int
		connection_id int
		question_id int
		first_answer string
		second_answer string
		answer_date string
	}{}
	answertests := []struct {
		answer_id int
		connection_id int
		question_id int
		first_answer string
		second_answer string
		answer_date string
	}{}
	r, err = db.Query("SELECT answer_id, connection_id, question_id, first_answer, second_answer, answer_date FROM answer")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("ANSWER TEST ERROR")
	}
	defer r.Close()
	for r.Next() {
		r.Scan(&answertest.answer_id, &answertest.connection_id, &answertest.question_id, &answertest.first_answer, &answertest.second_answer, &answertest.answer_date)
		answertests = append(answertests, answertest)
	}
	fmt.Println("NOW STORED ANSWER : ", answertests)		
	

// 회원가입	
	ginEngine.POST("/api/usr", func (c *gin.Context){
		signUpData := UsrsData{}
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
		
		uuid := uuid.New().String()

		_, err = db.Query(`INSERT INTO usrs (id, password, uuid, conn_id) VALUES ("`+signUpData.ID+`", "`+signUpData.Password+`", "`+uuid+`", 0)`)
		if err != nil {
			fmt.Println("ERROR #9 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.Writer.WriteHeader(http.StatusOK)
	})

// 회원가입 시 아이디 중복체크
	ginEngine.POST("/api/id", func (c *gin.Context){
		input := struct {
			ID string `json:"input_id"`
		}{}
		
		err := c.ShouldBindJSON(&input)
		if err != nil {
			fmt.Println("ERROR #4 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		r, err = db.Query(`SELECT * FROM usrs WHERE id = "`+input.ID+`"`)
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
	})	

// 로그인
	ginEngine.POST("/api/log", func (c *gin.Context){
		logInData := UsrsData{}
		err := c.ShouldBindJSON(&logInData)
		if err != nil {
			fmt.Println("ERROR #10 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		r, err := db.Query(`SELECT uuid FROM usrs WHERE id = "`+logInData.ID+`" and password = "`+logInData.Password+`"`)
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
	})

// 로그아웃
	ginEngine.DELETE("/api/log", func (c *gin.Context){
		uuid := cookieExist(c)

		c.SetCookie("uuid", uuid, -1, "/", os.Getenv("ORIGIN"), false, true)
		c.Writer.WriteHeader(http.StatusOK)
	})

// 기존 로그인 되있던 상태인지 쿠키 확인	
	ginEngine.GET("/api/log", func (c *gin.Context){
		uuid := cookieExist(c)

		r, err := db.Query(`SELECT * FROM usrs WHERE uuid = "`+uuid+`"`)
		if err != nil {
			fmt.Println("ERROR #12 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Close()

		if r.Next() {
			if isConnected(c, db) {
				c.String(http.StatusOK, "%v", "CONNECTED")
			} else {
				c.String(http.StatusOK, "%v", "NOT_CONNECTED")
			}	
		}
	})

// 상대방에게 connection 연결 요청	
	ginEngine.POST("/api/request", func (c *gin.Context){
		uuid := cookieExist(c)

		r1, err := db.Query(`SELECT * FROM request WHERE requester_uuid = "`+uuid+`"`)
		if err != nil {
			fmt.Println("ERROR #19 : ", err.Error())
		}
		defer r.Close()
		
		if r1.Next() {
			c.String(http.StatusBadRequest, "%v", "ALREADY_REQUEST")
			return
		}

		var id string
		r2, _ := db.Query(`SELECT id FROM usrs WHERE uuid = "`+uuid+`"`)
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
		r3, err := db.Query(`SELECT conn_id, uuid FROM usrs WHERE id = "`+input.ID+`"`)
		if err != nil {
			fmt.Println("ERROR #21 : ", err.Error())
		}
		defer r3.Close()

		// ID가 존재하는 ID면 이미 연결되어있진 않은지 conn_id를 확인
		if r3.Next() {
			var targetConnID int
			var targetUUID string
			r2.Scan(&targetConnID, &targetUUID)
			
			if targetUUID == uuid {
				c.String(http.StatusBadRequest, "%v", "NOT_YOURSELF")
			} else if targetConnID != 0 {
				c.String(http.StatusBadRequest, "%v", "ALREADY_CONNECTED")
			} else {
				// 요청된 정보를 DB에 저장
				_, err = db.Query(`INSERT INTO request (requester_uuid, target_uuid, request_time, requester_id, target_id) VALUES ("`+uuid+`", "`+targetUUID+`", "`+time.Now().Format("01/02 15:04")+`", "`+id+`", "`+input.ID+`")`)
				if err != nil {
					fmt.Println("ERROR #22 : ", err.Error())
				}
				c.Writer.WriteHeader(http.StatusOK)
			}
		} else {
		// ID가 존재하지 않는 ID면
			c.String(http.StatusBadRequest, "%v", "NOT_EXIST")
		}
	})

// 현재 요청받은 request 목록 가져오기
	ginEngine.GET("/api/request/recieved", func (c *gin.Context){
		uuid := cookieExist(c)

		requestingData := RequestData{}
		requestingDatas := []RequestData{}

		r, err := db.Query(`SELECT requester_id, requester_uuid, request_time, request_id FROM request WHERE target_uuid = "`+uuid+`"`)
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
	})


// 현재 신청중인 request 가져오기
	ginEngine.GET("/api/request/send", func (c *gin.Context){
		uuid := cookieExist(c)

		r, err := db.Query(`SELECT target_uuid, request_time, target_id FROM request WHERE requester_uuid = "`+uuid+`"`)
		if err != nil {
			fmt.Println("ERROR #16 : ", err.Error())
		}
		defer r.Close()

		requestingData := RequestData{}
		for r.Next(){
			r.Scan(&requestingData.Target_uuid, &requestingData.Request_time, &requestingData.Target_id)
		}

		marshaledRequestingData, err := json.Marshal(requestingData)
		if err != nil {
			fmt.Println("ERROR #18 : ", err.Error())
		}
		
		c.Writer.Write(marshaledRequestingData)
	})	

// 상대방과 연결 후, DB에 저장되어있던 자신과 상대 관련 요청 전체 삭제 + conn_id 생성
	ginEngine.PUT("/api/request", func (c *gin.Context){
		myUUID  := cookieExist(c)
		
		targetUUID := struct {
			UUID string `json:"uuid_delete"`
		}{}

		err = c.ShouldBindJSON(&targetUUID)
		if err != nil {
			fmt.Println("ERROR #23 : ", err.Error())
			return
		}

		_, err = db.Query(`INSERT INTO connection (first_usr, second_usr, start_date) VALUES ("`+targetUUID.UUID+`", "`+myUUID+`", "`+time.Now().Format("2006/01/02")+`")`)
		if err != nil {
			fmt.Println("ERROR #24 : ", err.Error())
		}

		r, err = db.Query(`SELECT connection_id FROM connection WHERE first_usr = "`+targetUUID.UUID+`" and second_usr = "`+myUUID+`"`)
		if err != nil {
			fmt.Println("ERROR #25 : ", err.Error())
		}
		defer r.Close()

		r.Next()
		var connID int
		r.Scan(&connID)

		_, err = db.Query(`UPDATE usrs SET order_usr = 1, conn_id = `+strconv.Itoa(connID)+` WHERE uuid = "`+targetUUID.UUID+`"`)
		if err != nil {
			fmt.Println("ERROR #26 : ", err.Error())
			return
		}

		_, err = db.Query(`UPDATE usrs SET order_usr = 2, conn_id = `+strconv.Itoa(connID)+` WHERE uuid = "`+myUUID+`"`)
		if err != nil {
			fmt.Println("ERROR #27 : ", err.Error())
			return
		}

		_, err = db.Query(`DELETE FROM request WHERE requester_uuid = "`+targetUUID.UUID+`" or target_uuid = "`+targetUUID.UUID+`" or requester_uuid = "`+myUUID+`" or target_uuid = "`+myUUID+`"`)
		if err != nil {
			fmt.Println("ERROR #28 : ", err.Error())
			return
		}
	})

// 받은 요청 중 선택해서 요청을 삭제
	ginEngine.DELETE("/api/request/:param", func (c *gin.Context){
		request_id := c.Param("param")

		_, err = db.Query(`DELETE FROM request WHERE request_id = `+request_id)
		if err != nil {
			fmt.Println("ERROR #29 : ", err.Error())
		}
	})

// 그동안 답한 내용들을 모아서 보여주기 위한 API
	ginEngine.GET("/api/answer", func (c *gin.Context){
		uuid := cookieExist(c)

		r1, err := db.Query(`SELECT connection_id FROM connection WHERE first_usr = "`+uuid+`" or second_usr = "`+uuid+`"`)
		if err != nil {
			fmt.Println("ERROR #30 : ", err.Error())
		}
		defer r.Close()

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
		r2, err := db.Query(`SELECT first_answer, second_answer, answer_date, question_id FROM answer WHERE connection_id = "`+connID+`"`)
		if err != nil {
			fmt.Println("ERROR #31 : ", err.Error())
		}
		defer r2.Close()
		for r2.Next() {
			r.Scan(&answerData.FirstAnswer, &answerData.SecondAnswer, &answerData.AnswerDate, &questionID)
			r3, err := db.Query(`SELECT question_contents FROM question WHERE question_id = `+questionID)
			if err != nil {
				fmt.Println("ERROR #32 : ", err.Error())
			}
			defer r3.Close()

			r3.Next()
			r3.Scan(&answerData.QuestionContents)
			answerDatas = append(answerDatas, answerData)
		}

		mashaledAnswerData, err := json.Marshal(answerDatas)
		if err != nil {
			fmt.Println("ERROR #33 : ", err.Error())
		}
		
		c.Writer.Write(mashaledAnswerData)
	})
	
// Websocket 프로토콜로 업그레이드 및 메시지 read/write
	ginEngine.GET("/ws", func(c *gin.Context){
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

		r1, err := db.Query(`SELECT first_usr, second_usr, connection_id FROM connection WHERE first_usr = "`+uuid+`" or second_usr = "`+uuid+`"`)
		if err != nil {
			fmt.Println("ERROR #36 : ", err.Error())
		}
		defer r.Close()

		r1.Next()
		var first_uuid, second_uuid string
		var conn_id int
		r1.Scan(&first_uuid, &second_uuid, &conn_id)

		// 기존 저장되어있던 채팅 DB에서 불러와서 표시
		initialChat := ChatData{}
		initialChats := []ChatData{}
		r2, err := db.Query(`SELECT chat_id, writer_id, write_time, text_body FROM chat WHERE writer_id = "`+first_uuid+`" or writer_id = "`+second_uuid+`" ORDER BY chat_id ASC`)
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
			var chatData []ChatData

			err := conn.ReadJSON(&chatData)
			if err != nil {
				fmt.Println("ERROR #39 : ", err.Error())
				break;
			}

			// 일반채팅이면 chat table에 저장, question에 대한 답이면 answer table에 저장
			if chatData[0].Is_answer == 0 {
				_, err = db.Query(`INSERT INTO chat (text_body, writer_id, write_time) VALUES ("`+chatData[0].Text_body+`", "`+uuid+`", "`+chatData[0].Write_time+`")`)
				// 어차피 커넥션 당 메시지 하나씩 전송 받으니까 slice index는 0으로 설정
				if err != nil {
					fmt.Println("ERROR #40 : ", err.Error())
				}	
			} else {
				r3, err := db.Query(`SELECT * FROM answer WHERE connection_id = `+strconv.Itoa(conn_id)+` and question_id = `+strconv.Itoa(chatData[0].Question_id))
				if err != nil {
					fmt.Println("ERROR #41 : ", err.Error())
				}
				defer r3.Close()

				if r3.Next() {
					if first_uuid == uuid {
						_, err = db.Query(`UPDATE answer SET first_answer = "`+chatData[0].Text_body+`" WHERE question_id = `+strconv.Itoa(chatData[0].Question_id))
					} else {
						_, err = db.Query(`UPDATE answer SET second_answer = "`+chatData[0].Text_body+`" WHERE question_id = `+strconv.Itoa(chatData[0].Question_id))
					}
				} else {
					_, err = db.Query(`INSERT INTO answer (connection_id, question_id, answer_date) VALUES (`+strconv.Itoa(conn_id)+`,`+strconv.Itoa(chatData[0].Question_id)+`, "`+chatData[0].Write_time+`")`)
					if err != nil {
						fmt.Println("ERROR #42 : ", err.Error())
					}
					if first_uuid == uuid {
						_, err = db.Query(`UPDATE answer SET first_answer = "`+chatData[0].Text_body+`" WHERE question_id = `+strconv.Itoa(chatData[0].Question_id))
					} else {
						_, err = db.Query(`UPDATE answer SET second_answer = "`+chatData[0].Text_body+`" WHERE question_id = `+strconv.Itoa(chatData[0].Question_id))
					}
					if err != nil {
						fmt.Println(err.Error())
						fmt.Println("ADD CHAT TO DB ERROR OCCURED")
					}
				}
			}
			
			first_conn := conns[first_uuid]
			second_conn := conns[second_uuid]

			target_conn := []*websocket.Conn{}
			target_conn = append(target_conn, first_conn, second_conn)
			// 커넥션 연결이 안되어있으면 보내면 nil pointer 오류 생김

			// 모든 커넥션에 메시지 write
			if chatData[0].Is_answer != 1 {
				for index, item := range target_conn {
					err := item.WriteJSON(chatData)
					if err != nil {
						fmt.Println("ERROR #43 : ", err.Error())
					}
					fmt.Println("index : ", index)
				}
			}
			
	// 채팅 중 단어가 발견되면 단어 관련된 질문을 커플에게 던지는 기능
			// 1. 단어를 먼저 다 뽑아서
			var target_word, question_contents string
			var question_id int
			r, err := db.Query("SELECT target_word, question_id, question_contents FROM question ORDER BY question_id ASC")
			if err != nil {
				fmt.Println("ERROR #44 : ", err.Error())
			}
			for r.Next() {
				// 2. 방금 READ한 채팅에 단어가 있는지 돌면서 확인
				r.Scan(&target_word, &question_id, &question_contents)	
				if strings.Contains(chatData[0].Text_body, target_word) {
					// 3. 단어가 발견되면 이전에 답을 한 전적이 있는지 검색
					fmt.Println(target_word)
					r, err := db.Query(`SELECT * FROM answer WHERE connection_id = `+strconv.Itoa(conn_id)+` and question_id = `+strconv.Itoa(question_id))
					if err != nil {
						fmt.Println("ERROR #45 : ", err.Error())
					}
					defer r.Close()
					// 4. 단어도 발견됐고, 이전에 했던 질문도 아니면 질문 WRITE
					if !r.Next() {
						questiondata := ChatData{
							question_contents,
							"question",
							time.Now().Format("http.StatusOK6/01/02 03:04"),
							1,
							0,
							question_id,
						}
						questiondatas := []ChatData{}
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
	})
	ginEngine.Run(":8080")
}
