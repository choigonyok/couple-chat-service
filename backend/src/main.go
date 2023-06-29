package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Origin CORS 설정
func OriginConfig() cors.Config{
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost"} 
	// 허용할 오리진 설정, 원래 리액트의 port가 아니라 리액트가 있는 container의 port 번호를 origin allow 해줘야함
	// localhost:3000로 origin allow 하면 통신 안됨

	config.AllowMethods= []string{"GET"}
	config.AllowHeaders = []string{"Content-type"}
	config.AllowCredentials = true
	return config
}

type db_test struct {
	Id int	`json:"id"`
	// 클라이언트랑 통신하려면 field name을 UPPER로 적어야함
}


func main() {
	eg := gin.Default()
	// 엔진 생성

	config := OriginConfig()
	eg.Use(cors.New(config)) 
	// origin 설정하고 설정한 config를 gin engine에서 사용하겠다는 이 부분이 있어야 적용이 됨!

	//db := makeDbConn()
	
	db, err := sql.Open("mysql", "root:password@tcp(mysql)/chat")
	// 도커에서는 localhost가 안먹혀서 통신이 안됨
	// ip는 terminal에서 curl ifconfig.me 로 확인가능
	// https://covenant.tistory.com/198 보고 설정하기
	// 처음에 로컬 서버 ip 적었다가
	// 안돼서 docker inspect로 mysql 컨테이너 ip 확인하고 적었는데 계속 핑이 안감
	// 그래서 컨테이너 이름을 적어줌. 이러면 도커가 알아서 ip주소와 포트까지 연결해줌
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("DB CONNECTING ERROR OCCURED")
	}
	defer db.Close()

// DB와 서버가 연결 되었는지 확인
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("PING TO DB REJECTED")
	}

// TEST용 table 생성
	_, err = db.Query(`CREATE TABLE test (id int)`)
	// err를 선언해놓고 에러처리 등으로 err를 사용하지 않으면 오류가 발생함
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("FAILED TO CREATE TEST TABLE ERROR OCCURED")
	}

// TEST를 위한 레코드 삽입
	_, _ = db.Query(`INSERT INTO test (id) VALUES (100)`)

// DB-SERVER 연결 확인용 테스트 API
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

// CLIENT-SERVER 연결 확인용 테스트 API
	eg.GET("/api/test", func (c *gin.Context){
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte("TEST"))
	})

	eg.Run(":8080")
}
