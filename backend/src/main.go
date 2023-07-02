package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

type ReadData struct {
	Text_body string `json:"text_body"`
        Writer_id string `json:"writer_id"`
        Write_time string `json:"write_time"`
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


func main() {

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

	//db := makeDbConn()
	
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

// TEST : TEST를 위한 레코드 삽입
	_, _ = db.Query(`DELETE FROM test WHERE id = 100`)
	_, _ = db.Query(`INSERT INTO test (id) VALUES (100)`)

// TEST : DB-SERVER 연결 확인용 테스트 API
	eg.GET("/api/usr", func (c *gin.Context){
		data, err := db.Query("SELECT id FROM test")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("READING DATA FROM DB ERROR OCCURED")

		}
		var test_data db_test
		var test_datas []db_test
		for data.Next() {
			data.Scan(&test_data.Id)
			test_datas = append(test_datas, test_data)
		}
		fmt.Println(test_datas)
		send_data, err := json.Marshal(test_datas)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("DATA MARSHALING ERROR OCCURED")
		}
		c.Writer.WriteHeader(200)
		c.Writer.Write(send_data)
	})

// TEST : CLIENT-SERVER 연결 확인용 테스트 API
	eg.GET("/api/test", func (c *gin.Context){
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte("TEST"))
	})	
	
// Websocket 프로토콜로 업그레이드
	eg.GET("/ws", func(c *gin.Context){

		
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

		// 사용자에게 uuid를 생성해서 전달
		uuid := GenerateUserID()
		fmt.Println("USRID : ", uuid)
		conn.WriteJSON(struct{
			ID string `json:"created_id"`
		}{
			uuid,
		})
		// usrID := struct {
		// 	Usr_id string `json:"usr_id"`
		// }{
		// 	uuid,
		// }
		
		
		for {
			var read_data ReadData
			err := conn.ReadJSON(&read_data)
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
			fmt.Println("READ_TEXT : ", read_data.Text_body)
			fmt.Println("READ_ID : ", read_data.Writer_id)
			fmt.Println("READ_TIME : ", read_data.Write_time)
			
			
			err = conn.WriteJSON(read_data)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("WRITING TO CONN ERROR OCCURED")
			}
		}
		
	})
	
	eg.Run(":8080")
}
