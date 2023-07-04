package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"

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
	Conn_id int `json:"conn_id"`
	Chat_id int `json:"chat_id"`
}

type UsrInfo struct {
	Usr_ID string `json:"usr_id"`
	Usr_PW string `json:"usr_pw"`
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
func GenerateUserID() string {
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


func main() {	
// 커넥션 집합 슬라이스
	var conns []*websocket.Conn

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
	r, err := db.Query("SELECT id, password FROM usrs")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("USR TEST ERROR")
	}
	for r.Next()  {
		r.Scan(&test.Usr_ID, &test.Usr_PW)
		tests = append(tests, test)
	}
	fmt.Println("NOW STORED USR ID AND PW : ", tests)

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
			uuid := GenerateUserID()

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

		// 전체 커넥션 슬라이스에 커넥션 추가
		conns = append(conns, conn)
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

		// 기존 저장되어있던 채팅 DB에서 불러와서 표시
		initialChat := MessageData{}
		initialChats := []MessageData{}
		r, err := db.Query(`SELECT chat_id, writer_id, write_time, text_body FROM chat`)
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
			
			// 모든 커넥션에 메시지 write 
			for index, item := range conns {
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
