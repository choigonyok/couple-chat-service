package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// 클라이언트에서 채팅 본문을 받는 구조체
type ChatBody struct {
	Chat_body string `json:"chat_body"`
}

type MessageData struct {
	Text_body string `json:"text_body"`
        Writer_id string `json:"writer_id"`
        Write_time string `json:"write_time"`
	Conn_id string `json:"conn_id"`
	Chat_id int `json:"chat_id"`
}

type RequestData struct {
	Request_id int
	Requester_uuid string
	Requester_id string
	Target_uuid string
	Target_id string
	Request_time string
	
}

type DeleteUUID struct {
	Uuid string `json:"uuid_delete"`
}

type UsrInfo struct {
	Usr_ID string `json:"usr_id"`
	Usr_PW string `json:"usr_pw"`
	Usr_ConnID string
 }

// Origin CORS 설정
func OriginConfig() cors.Config{
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("ORIGIN1")} 
	// 허용할 오리진 설정, 원래 리액트의 port가 아니라 리액트가 있는 container의 port 번호를 origin allow 해줘야함
	// localhost:3000로 origin allow 하면 통신 안됨

	config.AllowMethods= []string{"GET"}
	config.AllowHeaders = []string{"Content-type"}
	config.AllowCredentials = true
	return config
}

// TEST : TEST를 위한 db 스트럭쳐
type db_test struct {
	Id int	`json:"id"`
	// 클라이언트랑 통신하려면 field name을 UPPER로 적어야함
}

// 커넥션 별 uuid 생성
func GenerateUID() string {
	u := uuid.New()
	return u.String()
}

func checkIDandPWLength(idorpw string) bool {
	if len(idorpw) >= 21 {
		return false
	} else {
		return true
	}
}

// usr가 상대방과 연결된 상태인지 아닌지 체크
func isConnected(c *gin.Context, db *sql.DB) bool {
	// 함수명과 파라미터 띄어쓰면 오류 생김
	
	uuid, err := c.Cookie("uuid")
	fmt.Println(uuid)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("COOKIE LOADING TO CHECK CONNECTED ERROR OCCURED")
	}
	r, err := db.Query(`SELECT conn_id FROM usrs WHERE uuid = "`+uuid+`"`)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("LOAD DB TO CHECK CONNECTED ERROR OCCURED")
	}
	var conn_id string
	for r.Next() {
		r.Scan(&conn_id)
	}
	if conn_id == "0" {
		return false
	} else {
		return true
	}
}

func main() {	
// 커넥션 집합 슬라이스
	conns := make(map[string]*websocket.Conn)

// 환경변수 로딩
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("LOADING ENV FILE ERROR OCCURED")
	}

// TEST : 환경변수 잘 불러와지는지 테스트
	DB_PASSWORD, isExist := os.LookupEnv("DB_PASSWORD")
	fmt.Println(DB_PASSWORD)
	if isExist == false {
		fmt.Println("LOADING ENV VAR DB_PASSWORD ERROR OCCURED")
	}
	DB_NAME, isExist := os.LookupEnv("DB_NAME")
	if isExist == false {
		fmt.Println("LOADING ENV VAR DB_NAME ERROR OCCURED")
	}
	fmt.Println(DB_NAME)
	// GO의 환경변수 설정은 os 패키지 이용
	// godotenv 패키지를 이용해서 .env파일에 환경변수를 설정해주고
	// .env파일은 보안을 위해서 공유저장소에 올라가지 않도록 .gitignore에 설정
	// os.Getenv() 메서드는 환경변수가 없어도 empty, 있는데 설정이 안되어있어도 empty를 리턴하기 때문에
	// ambiguity를 줄이기 위해서 os.LookupEnv() 메서드를 사용
	// 환경변수의 존재 여부를 두 번째 파라미터 boolean으로 알려줌
	// os.Setenv()와 os.UnSetenv()로 환경변수를 생성/삭제 할 수 있음

	eg := gin.Default()
	// 엔진 생성

	config := OriginConfig()
	eg.Use(cors.New(config)) 
	// origin 설정하고 설정한 config를 gin engine에서 사용하겠다는 이 부분이 있어야 적용이 됨!
	
	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp(mysql)/"+os.Getenv("DB_NAME"))
	// 도커에서는 localhost가 안먹혀서 통신이 안됨
	// ip는 terminal에서 curl ifconfig.me 로 확인가능
	// https://covenant.tistory.com/198 보고 설정하기
	// 처음에 로컬 서버 ip 적었다가
	// 안돼서 docker inspect로 mysql 컨테이너 ip 확인하고 적었는데 계속 핑이 안감
	// 그래서 컨테이너 이름을 적어줌. 이러면 도커가 알아서 ip주소와 포트까지 연결해줌
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("DB CONNECTING ERROR OCCURED!")
	}
	defer db.Close()

// DB와 서버가 연결 되었는지 확인
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("PING TO DB REJECTED")
	}

// TEST : DB에 usr 정보가 잘 저장되는지 테스트
	test := UsrInfo{}
	tests := []UsrInfo{}
	r, err := db.Query("SELECT id, password, uuid FROM usrs")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("USER TEST ERROR")
	}
	for r.Next()  {
		r.Scan(&test.Usr_ID, &test.Usr_PW, &test.Usr_ConnID)
		tests = append(tests, test)
	}
	fmt.Println("NOW STORED USR ID AND PW AND CONN_ID: ", tests)

// TEST : DB에 chat data가 잘 저장되는지 테스트
	chattest := MessageData{}
	chattests := []MessageData{} 
	r, err = db.Query("SELECT chat_id, writer_id, write_time, text_body FROM chat")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("CHAT TEST ERROR")
	}
	for r.Next() {
		r.Scan(&chattest.Chat_id,&chattest.Writer_id,&chattest.Write_time,&chattest.Text_body)
		chattests = append(chattests, chattest)
	}
	fmt.Println("NOW STORED CHAT : ", chattests)

// TEST : DB에 request data가 잘 저장되는지 테스트
	requesttest := RequestData{}
	requesttests := []RequestData{}
	r, err = db.Query("SELECT request_id, requester_uuid, target_uuid, request_time FROM request")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("REQUEST TEST ERROR")
	}
	for r.Next() {
		r.Scan(&requesttest.Request_id, &requesttest.Requester_uuid, &requesttest.Target_uuid, &requesttest.Request_time)
		requesttests = append(requesttests, requesttest)
	}
	fmt.Println("NOW STORED REQUEST : ", requesttests)

// TEST : DB에 connection data가 잘 저장되는지 테스트
	connectiontest := struct {
		connection_id string
		first_usr string
		second_usr string
		start_date string
	}{}
	connectiontests := []struct {
		connection_id string
		first_usr string
		second_usr string
		start_date string
	}{}
	r, err = db.Query("SELECT connection_id, first_usr, second_usr, start_date FROM connection")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("CONNECTION TEST ERROR")
	}
	for r.Next() {
		r.Scan(&connectiontest.connection_id, &connectiontest.first_usr, &connectiontest.second_usr, &connectiontest.start_date)
		connectiontests = append(connectiontests, connectiontest)
	}
	fmt.Println("NOW STORED CONNECTION : ", connectiontests)	

	// _, err = db.Query("DELETE FROM chat")
	// _, err = db.Query("DELETE FROM connection")
	// _, err = db.Query("DELETE FROM usrs")
	// _, err = db.Query("DELETE FROM request")

// 회원가입 시 아이디 중복체크
	eg.POST("/api/id", func (c *gin.Context){
		temp := struct {
			InputID string `json:"input_id"`
		}{}
		
		err := c.ShouldBindJSON(&temp)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("BIDING ID FOR SIGNUP ERROR OCCURED")
		}
		r, err = db.Query(`SELECT * FROM usrs WHERE id = "`+temp.InputID+`"`)
		if err != nil {
			fmt.Println(err.Error())
		}
		if r.Next() {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(200)
		}
		
	})	

// 회원가입	
	eg.POST("/api/usr", func (c *gin.Context){
		data := UsrInfo{}
		err := c.ShouldBindJSON(&data)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("BINDING SIGNUP DATA ERROR OCCURED")
		} else {
			if !checkIDandPWLength(data.Usr_ID) {
				c.String(400, "%v", "ID의 최대 길이는 20자로 제한됩니다.")
				return
			}
			if !checkIDandPWLength(data.Usr_PW) {
				c.String(400, "%v", "PASSWORD의 최대 길이는 20자로 제한됩니다.")
				return
			}
			idCorrect, err := regexp.MatchString("^[a-z][a-z0-9]+$", data.Usr_ID)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("ID REGEXP ERROR OCCURED")
			}
			if !idCorrect {
				c.String(400, "%v", "ID는 첫 글자가 영어 소문자인, 영어 소문자와 숫자 조합의 1~20자로만 사용할 수 있습니다.")
				return
			}
			pwCorrect, err := regexp.MatchString("^[a-z0-9]*$", data.Usr_ID)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("ID REGEXP ERROR OCCURED")
			}
			if !pwCorrect {
				c.String(400, "%v", "PW는 첫 글자가 영어 소문자인, 영어 소문자와 숫자 조합의 1~20자로만 사용할 수 있습니다.")
				return
			}
			// 사용자의 uuid를 생성
			uuid := GenerateUID()

			// DB에 사용자 데이터 저장
			_, err = db.Query(`INSERT INTO usrs (id, password, uuid, conn_id) VALUES ("`+data.Usr_ID+`", "`+data.Usr_PW+`", "`+uuid+`", 0)`)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("STORING SIGNUP DATA TO DB ERROR OCCURED")
			} else {
				c.Writer.WriteHeader(http.StatusOK)
			}
		}
	})

// 로그인
	eg.POST("/api/log", func (c *gin.Context){
		data := UsrInfo{}
		err := c.ShouldBindJSON(&data)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOGIN DATA BINDING ERROR OCCURED")
		}
		r, err := db.Query(`SELECT uuid FROM usrs WHERE id = "`+data.Usr_ID+`" and password = "`+data.Usr_PW+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("FINDING USR DATA TO LOGIN FROM DB ERROR OCCURED")
		}
		if r.Next() {
			var uuid_data string
			r.Scan(&uuid_data)
			c.SetCookie("uuid", uuid_data, 60*60, "/", os.Getenv("ORIGIN1"),false,true)	
			c.Writer.WriteHeader(200)
		} else {
			c.Writer.WriteHeader(400)
		}
	})

// 로그아웃
	eg.DELETE("/api/log", func (c *gin.Context){
		uuid, err := c.Cookie("uuid")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOADING COOKIE TO DELETE ERROR OCCURED")
		}
		c.SetCookie("uuid", uuid, -1, "/", os.Getenv("ORIGIN1"), false, true)
		c.String(200, "로그아웃 되었습니다.")
	})

// 기존 로그인 되있던 상태인지 쿠키 확인	
	eg.GET("/api/log", func (c *gin.Context){
		// 쿠키의 uuid 확인
		uuid, err := c.Cookie("uuid")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOADING COOKIE TO CHECK LOGIN ERROR OCCURED")
			c.Writer.WriteHeader(400)
			return
		}

		// 그 uuid가 진짜 db에 있는 회원정보인지 확인
		r, err := db.Query(`SELECT * FROM usrs WHERE uuid = "`+uuid+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("CHECKING DB TO CHECK LOGGED IN ERROR OCCURED")
		}

		// 있으면 상대와 connection이 연결된 상태인지 확인 후 응답
		if r.Next() {
			if isConnected(c, db) {
				c.String(200, "%v", "CONNECTED")
			} else {
				c.String(200, "%v", "NOT_CONNECTED")
			}	
		} else {
			c.Writer.WriteHeader(500)
		}
	})

// 현재 요청받은 request 목록 가져오기
	eg.GET("/api/request/recieved", func (c *gin.Context){
		uuid, err := c.Cookie("uuid")	
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOAD COOKIE TO LOAD REQUEST LIST ERROR OCCURED")
		}

		// usr가 요청한 커넥션 표시, 요청 받은 것과 달리 요청은 한 번만 할 수 있어서 slice 안함
		requesting_data := RequestData{}
		requesting_datas := []RequestData{}
		// usr가 요청받은 커넥션 표시, 요청을 여러개 받을 수 있어서 slice 사용함
		r, err := db.Query(`SELECT requester_uuid, request_time, request_id FROM request WHERE target_uuid = "`+uuid+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("DATA WHO REQUEST TO ME FROM DB ERROR OCCURED")
		}
		for r.Next() {
			r.Scan(&requesting_data.Requester_uuid, &requesting_data.Request_time, &requesting_data.Request_id)

			// uuid에 맞는 id 찾기
			rr, err := db.Query(`SELECT id FROM usrs WHERE uuid = "`+requesting_data.Requester_uuid+`"`)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("FINDING ID BASED UUID ERROR OCCURED")
			}
			var id string
			rr.Next()
			rr.Scan(&id)
			requesting_data.Requester_id = id

			requesting_datas = append(requesting_datas, requesting_data)
		}

		data, err := json.Marshal(requesting_datas)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("MARSHALING REQUESTING DATA ERROR OCCURED")
		}

		c.Writer.Write(data)
	})


// 현재 신청중인 request 가져오기
	eg.GET("/api/request/send", func (c *gin.Context){
		uuid, err := c.Cookie("uuid")	
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOAD COOKIE TO LOAD REQUEST LIST ERROR OCCURED")
		}

		r, err = db.Query(`SELECT target_uuid, request_time FROM request WHERE requester_uuid = "`+uuid+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("DATA REQUESTING FROM DB ERROR OCCURED")
		}

		requesting_data := RequestData{}
		for r.Next(){
			r.Scan(&requesting_data.Target_uuid, &requesting_data.Request_time)
		}

		// uuid 기반으로 id 찾기
		r, err = db.Query(`SELECT id FROM usrs WHERE uuid = "`+requesting_data.Target_uuid+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("FINDING ID BASED UUID ERROR OCCURED")
		}
		var id string
		r.Next()
		r.Scan(&id)
		requesting_data.Target_id = id

		data, err := json.Marshal(requesting_data)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("MARSHALING REQUESTING DATA ERROR OCCURED")
		}
		
		c.Writer.Write(data)
	})	

// 상대방에게 connection 연결 요청	
	eg.POST("/api/request", func (c *gin.Context){
		uuid, err := c.Cookie("uuid")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOAD COOKIE TO CHECK CALL YOURSELF ERROR OCCURED")
		}
		// 이미 요청한 상태인지 확인
		r, err := db.Query(`SELECT * FROM request WHERE requester_uuid = "`+uuid+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("CHECK WHETHER ALREADY REQUEST ERROR OCCURED")
		}
		if r.Next() {
			c.String(400, "%v", "ALREADY_REQUEST")
			return
		}

		data := struct {
			UsrID string `json:"input_id"`
		}{}
		err = c.ShouldBindJSON(&data)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("JSON BINDING TO REQUEST ERROR OCCURED")
		}

		// 입력한 ID에 맞는 사용자 DATA DB에서 불러오기
		r, err = db.Query(`SELECT id, conn_id, uuid FROM usrs WHERE id = "`+data.UsrID+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("FINDING ID, CONN_ID TO SEND REQUEST ERROR OCCURED")
		} 

		// ID가 존재하는 ID면 이미 연결되어있진 않은지 conn_id를 확인
		if r.Next() {
			var id_temp string
			var conn_id_temp string
			var uuid_temp string

			r.Scan(&id_temp, &conn_id_temp, &uuid_temp)
			if uuid_temp == uuid {
				c.String(400, "%v", "NOT_YOURSELF")
			} else if conn_id_temp != "0" {
				c.String(400, "%v", "ALREADY_CONNECTED")
				return
			} else {
				// 요청된 정보를 DB에 저장
				_, err = db.Query(`INSERT INTO request (requester_uuid, target_uuid, request_time) VALUES ("`+uuid+`", "`+uuid_temp+`", "`+time.Now().Format("01/02 15:04")+`")`)
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println("STORING REQUEST DATA TO DB ERROR OCCURED")
				}
				c.Writer.WriteHeader(200)
				return
			}
		} else {
		// ID가 존재하지 않는 ID면
			c.String(400, "%v", "NOT_EXIST")
		}
	})

// 상대방과 연결 후, DB에 저장되어있던 자신과 상대 관련 요청 전체 삭제 + conn_id 생성
	eg.PUT("/api/request", func (c *gin.Context){
		// 승인usr의 request data를 삭제하기 위한 쿠키
		firstUUID, err := c.Cookie("uuid")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOADING COOKIE TO DELETE REQUEST ERROR OCCURED")
			return
		}
		var data DeleteUUID
		err = c.ShouldBindJSON(&data)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("BINDING JSON TO DELETE REQUEST ERROR OCCURED")
			return
		}
		conn_id := GenerateUID()

		_, err = db.Query(`UPDATE usrs SET conn_id = "`+conn_id+`" WHERE uuid = "`+firstUUID+`" or uuid = "`+data.Uuid+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("UPDATE CONN_ID ERROR OCCURED")
			return
		}

		_, err = db.Query(`INSERT INTO connection (first_usr, second_usr, start_date) VALUES ("`+data.Uuid+`", "`+firstUUID+`", "`+time.Now().Format("2006/01/02")+`")`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("INSERTING CONNECTION DATA INTO DB ERROR OCCURED")
		}

		_, err = db.Query(`DELETE FROM request WHERE requester_uuid = "`+data.Uuid+`" or target_uuid = "`+data.Uuid+`" or requester_uuid = "`+firstUUID+`" or target_uuid = "`+firstUUID+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("DELETE REQUESTS FOR CONNECTION ERROR OCCURED")
			return
		}
	})

// 받은 요청 중 선택해서 요청을 삭제
	eg.DELETE("/api/request/:param", func (c *gin.Context){
		request_id := c.Param("param")

		_, err = db.Query(`DELETE FROM request WHERE request_id = `+request_id)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("DELETE SPECIFIC REQUEST FROM DB ERROR OCCURED")
		}
	})

	
// Websocket 프로토콜로 업그레이드 및 메시지 read/write
	eg.GET("/ws", func(c *gin.Context){
		uuid, err := c.Cookie("uuid")
		if err != nil {
			c.String(400, "로그인을 한 이후에 서비스를 이용할 수 있습니다.")
			return
		}

		var upgrader  = websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return origin == os.Getenv("ORIGIN1")
			    },
			    // Websocket의 Origin 검증은 서버에서 진행
			    // 브라우저는 호스트 상관없이 막 요청 보냄
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("WEBSOCKET UPGRADING ERROR OCCURED")
			return
		}
		defer conn.Close()

		// 이 usr의 uuid를 키로 넣으면 현재 conn이 값으로 나오는 map
		conns[uuid] = conn
		// conn객체를 읽어야함
		
		// LATER : 커넥션 종료/삭제되면 슬라이스에서도 제외해야 함
		
		
		// 클라이언트에 uuid 전달, 그래야 클라이언트에게 채팅을 표시할 때
		// 누가 보낸 채팅인지 UUID로 구분해서 표시할 수 있음
		json_uuid := struct {
			Uuid string `json:"uuid"`
		}{
			uuid,
		}
		err = conn.WriteJSON(json_uuid)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("SENDING UUID TO CLIENT ERROR OCCURED")
			return
		}

		r, err = db.Query(`SELECT first_usr, second_usr FROM connection WHERE first_usr = "`+uuid+`" or second_usr = "`+uuid+`"`)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("FINDING CONNECITON UUID ERROR OCCURED")
		}
		r.Next()
		var first_uuid, second_uuid string
		r.Scan(&first_uuid, &second_uuid)

		// 기존 저장되어있던 채팅 DB에서 불러와서 표시
		initialChat := MessageData{}
		initialChats := []MessageData{}
		r, err := db.Query(`SELECT chat_id, writer_id, write_time, text_body FROM chat WHERE writer_id = "`+first_uuid+`" or writer_id = "`+second_uuid+`" ORDER BY chat_id ASC`)
		// LATER : 나중에 여러 conn 구현하면 쿼리문에 조건절이랑 conn_id 넣기
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("LOADING INITIAL CHATS FROM DB ERROR OCCURED")
		}
		for r.Next() {
			r.Scan(&initialChat.Chat_id, &initialChat.Writer_id, &initialChat.Write_time, &initialChat.Text_body)
			initialChats = append(initialChats, initialChat)
		}
		err = conn.WriteJSON(initialChats)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("WRITING INITIAL CHATS TO CLIENTS ERROR OCCURED")
		}

		// 메시지를 읽고 쓰는 부분, 읽은 메시지는 DB에 저장됨
		for { 
			var messageData []MessageData
			// 메시지 하나씩 주고받는데 slice로 메시지 read하는 이유
			// : 기존 DB에 저장되어있던 메시지를 보낼 때 slice 형태로 전송하는데
			// 클라이언트에서 기존 메시지나, 새로운 입력 메시지나 하나의 코드로 처리할 수 있게 하려고 이렇게 작성함
			err := conn.ReadJSON(&messageData)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("READING FROM CONNECTION ERROR OCCURED")
				break;
			}
			//invalid character 'o' looking for beginning of value 에러 발생
			// ReadJSON에서 문제 발생
			// 출처 : https://austindewey.com/2020/12/11/troubleshooting-invalid-character-looking-for-beginning-of-value/
			// json패키지가 json형식이 아닌 스트링을 언마샬링하려고 할 때 발생하는 에러
			// 리액트코드 원인 : newSocket.send(JSON.stringify(sendData)); 객체만 만들고 객체를 json형식으로 변환을 안시켜줬음

			// DB에 메시지 저장
			_, err = db.Query(`INSERT INTO chat (text_body, writer_id, write_time) VALUES ("`+messageData[0].Text_body+`", "`+uuid+`", "`+messageData[0].Write_time+`")`)
			// 어차피 커넥션 당 메시지 하나씩 전송 받으니까 slice index는 0으로 설정
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("ADD CHAT TO DB ERROR OCCURED")
			}
			
			first_conn := conns[first_uuid]
			second_conn := conns[second_uuid]

			target_conn := []*websocket.Conn{}
			target_conn = append(target_conn, first_conn, second_conn)

			// 모든 커넥션에 메시지 write 
			for index, item := range target_conn {
				err := item.WriteJSON(messageData)
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println(index, "TH CONN WRITING ERROR OCCURED")
				}
				fmt.Println("index : ", index)
			}

			// err = conn.WriteJSON(messageData)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// 	fmt.Println("WRITING TO CONN ERROR OCCURED")
			// }
		}
	})
	
	eg.Run(":8080")
}

// uuid와 conn 사이의 키-값 저장을 위해 redis 도입하면 좋을 것 같음
// 서버 중지되면 conn 데이터 분실 위험

// 도커로 개발환경 구성해서 개발하면 정확히 어느 부분에서 에러가 발생한 건지 확인이 어려움
// 계속 fmt.Println(err.Error())를 반복해서 쓰니 코드 가독성도 안좋아지고 불필요하게 많은 코드가 작성됨
// 테스트 코드 도입의 필요성


// 백그라운드로 커넥션 실행하기
// 브라우저에 focus가 안되어있기만 해도
// websocket: close 1006 (abnormal closure): unexpected EOF
// 라면서 커넥션이 close되어서 코드 문제인지 네트워크 문제인지 구분이 안가서 개발하기가 힘듦