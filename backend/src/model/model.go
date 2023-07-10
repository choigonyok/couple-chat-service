package model

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

type QuestionData struct {
	Question_id int
	Target_word string
	Question_contents string
}

type AnswerData struct {
	Answer_id int 
	Connection_id int
	QuestionContents string `json:"question_contents"`
	FirstAnswer string `json:"first_answer"`
	SecondAnswer string `json:"second_answer"`
	AnswerDate string `json:"answer_date"`
	Question_id int
}

 var db *sql.DB

func OpenDB(driverName, dataSourceName string){
	database, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		fmt.Println("ERROR #2 : ", err.Error())
	}

	db = database

	// DB와 서버가 연결 되었는지 확인
	err = db.Ping()
	if err != nil {
		fmt.Println("ERROR #3 : ", err.Error())
	}
}

func CloseDB() {
	db.Close()
}

func InsertUsr(id, password, uuid string) error {
	_, err := db.Query(`INSERT INTO usrs (id, password, uuid, conn_id) VALUES ("`+id+`", "`+password+`", "`+uuid+`", 0)`)
	return err
}


func SelectUsrByID(id string) (*sql.Rows, error){
	r, err := db.Query(`SELECT * FROM usrs WHERE id = "`+id+`"`)
	return r, err
}

func SelectUUIDFromUsrsByIDandPW(id, password string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT uuid FROM usrs WHERE id = "`+id+`" and password = "`+password+`"`)
	return r, err
}

func SelectUsrByUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT * FROM usrs WHERE uuid = "`+uuid+`"`)
	return r, err
}

func SelectConnIDFromUsrsByUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT conn_id FROM usrs WHERE uuid = "`+uuid+`"`)
	return r, err
}

func SelectRequestByRequesterUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT * FROM request WHERE requester_uuid = "`+uuid+`"`)
	return r, err
}

func SelectIDFromUsrsByUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT id FROM usrs WHERE uuid = "`+uuid+`"`)
	return r, err
}

func SelectConnIDandUUIDFromUsrsByID(id string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT conn_id, uuid FROM usrs WHERE id = "`+id+`"`)
	return r, err
}

func InsertRequest(requester_uuid, target_uuid, request_time, requester_id, target_id  string) error {
	_, err := db.Query(`INSERT INTO request (requester_uuid, target_uuid, request_time, requester_id, target_id) VALUES ("`+requester_uuid+`", "`+target_uuid+`", "`+request_time+`", "`+requester_id+`", "`+target_id+`")`)
	return err
}

func SelectRecieveRequestByTargetUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT requester_id, requester_uuid, request_time, request_id FROM request WHERE target_uuid = "`+uuid+`"`)
	return r, err
}

func SelectSendRequestByTargetUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT target_uuid, request_time, target_id FROM request WHERE requester_uuid = "`+uuid+`"`)
	return r, err
}

func InsertConnection(first_usr, second_usr, start_date string) (*sql.Rows, error) {
	r, err := db.Query(`INSERT INTO connection (first_usr, second_usr, start_date) VALUES ("`+first_usr+`", "`+second_usr+`", "`+start_date+`")`)
	return r, err
}

func SelectConnectionIDByUsrsUUID(first_usr, second_usr string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT connection_id FROM connection WHERE first_usr = "`+first_usr+`" and second_usr = "`+second_usr+`"`)
	return r, err
}

func UpdateUsrsConnID(conn_id int, targetUUID string) (*sql.Rows, error) {
	r, err := db.Query(`UPDATE usrs SET order_usr = 1, conn_id = `+strconv.Itoa(conn_id)+` WHERE uuid = "`+targetUUID+`"`)
	return r, err
}

func UpdateUsrsOrder(conn_id int, myUUID string) (*sql.Rows, error) {
	r, err := db.Query(`UPDATE usrs SET order_usr = 2, conn_id = `+strconv.Itoa(conn_id)+` WHERE uuid = "`+myUUID+`"`)
	return r, err
}

func DeleteRestRequest(requester_uuid, target_uuid string) (*sql.Rows, error) {
	r, err := db.Query(`DELETE FROM request WHERE requester_uuid = "`+requester_uuid+`" or target_uuid = "`+requester_uuid+`" or requester_uuid = "`+target_uuid+`" or target_uuid = "`+target_uuid+`"`)
	return r, err
}

func DeleteRequestByRequestID(request_id string) error {
	_, err := db.Query(`DELETE FROM request WHERE request_id = `+request_id)
	return err
}

func SelectConnIDByUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT connection_id FROM connection WHERE first_usr = "`+uuid+`" or second_usr = "`+uuid+`"`)
	return r, err
}

func SelectAnswerByConnID(connection_id string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT first_answer, second_answer, answer_date, question_id FROM answer WHERE connection_id = "`+connection_id+`"`)
	return r, err
}

func SelectQuestionContentsByQuestionID(question_id string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT question_contents FROM question WHERE question_id = `+question_id)
	return r, err
}

func SelectConnectionByUsrsUUID(uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT first_usr, second_usr, connection_id FROM connection WHERE first_usr = "`+uuid+`" or second_usr = "`+uuid+`"`)
	return r, err
}

func SelectChatByUsrsUUID(first_uuid, second_uuid string) (*sql.Rows, error) {
	r, err := db.Query(`SELECT chat_id, writer_id, write_time, text_body FROM chat WHERE writer_id = "`+first_uuid+`" or writer_id = "`+second_uuid+`" ORDER BY chat_id ASC`)
	return r, err
}

func InsertChat(text_body, writer_id, write_time string) error {
	_, err := db.Query(`INSERT INTO chat (text_body, writer_id, write_time) VALUES ("`+text_body+`", "`+writer_id+`", "`+write_time+`")`)
	return err
}

func SelectAnswerByConnIDandQuestionID(connection_id, question_id int) (*sql.Rows, error) {
	r, err := db.Query(`SELECT * FROM answer WHERE connection_id = `+strconv.Itoa(connection_id)+` and question_id = `+strconv.Itoa(question_id))
	return r, err
}

func UpdateFirstAnswerByQuestionID(first_answer string, question_id int) error {
	_, err := db.Query(`UPDATE answer SET first_answer = "`+first_answer+`" WHERE question_id = `+strconv.Itoa(question_id))
	return err
}

func UpdateSecondAnswerByQuestionID(first_answer string, question_id int) error {
	_, err := db.Query(`UPDATE answer SET second_answer = "`+first_answer+`" WHERE question_id = `+strconv.Itoa(question_id))
	return err
}

func InsertAnswer(answer_date string, connection_id, question_id int) error {
	_, err := db.Query(`INSERT INTO answer (connection_id, question_id, answer_date) VALUES (`+strconv.Itoa(connection_id)+`,`+strconv.Itoa(question_id)+`, "`+answer_date+`")`)
	return err
}

func SelectQuetions() (*sql.Rows, error) {
	r, err := db.Query("SELECT target_word, question_id, question_contents FROM question ORDER BY question_id ASC")
	return r, err
}

func GetUsrOrderByUUID(uuid string) int {
	r, err := db.Query(`SELECT order_usr FROM usrs WHERE uuid = "`+uuid+`"`)
	if err != nil {
		fmt.Println("ERROR #51 : ", err.Error())
	}
	r.Next()
	var order_usr int
	r.Scan(&order_usr)
	return order_usr
}

func QuestionIDOfEmptyAnswerByOrder(order int) int {
	question_id := 0
	if order == 1 {
		r, err := db.Query(`SELECT question_id FROM ANSWER WHERE first_answer = "not-written"`)
		if err != nil {
			fmt.Println("ERROR #52 : ", err.Error())
		}
		if r.Next() {
			r.Scan(&question_id)
			return question_id
		}
	} else {
		r, err := db.Query(`SELECT question_id FROM ANSWER WHERE second_answer = "not-written"`)
		if err != nil {
			fmt.Println("ERROR #53 : ", err.Error())
		}
		if r.Next() {
			r.Scan(&question_id)
			return question_id
		}
	}
	return question_id
}

func GetQuestionByQuestionID(questionID int) (string, string){
	r, err := db.Query(`SELECT target_word, question_contents FROM question WHERE question_id = `+ strconv.Itoa(questionID))
	if err != nil {
		fmt.Println("ERROR #54 : ", err.Error())
	}
	r.Next()
	var questionData QuestionData
	r.Scan(&questionData.Target_word, &questionData.Question_contents)
	return questionData.Target_word, questionData.Question_contents
}

func GetRecentAnswerByConnID(connection_id, num int) []AnswerData {
	r, err := db.Query("SELECT * FROM answer ORDER BY answer_id DESC LIMIT " + strconv.Itoa(num))
	if err != nil {
		fmt.Println("ERROR #55 : ", err.Error())
	}
	var answerData AnswerData
	var answerDatas []AnswerData
	for r.Next() {
		r.Scan(&answerData.Answer_id, &answerData.Connection_id, &answerData.FirstAnswer, &answerData.SecondAnswer, &answerData.AnswerDate, &answerData.Question_id)
		answerDatas = append(answerDatas, answerData)
	}
	return answerDatas
}

func GetFrequentFiveWord(uuid string){

	r, err := db.Query(`SELECT text_body FROM chat WHERE writer_id = "`+uuid+`"`)
	if err != nil {
		fmt.Println("ERROR #56 : ", err.Error())
	}
	var recentChat string
	var recentChats []string
	for r.Next() {
		r.Scan(&recentChat)
		fmt.Println(recentChat)
	fmt.Println(recentChat)
	fmt.Println(recentChat)
		recentChats = append(recentChats, recentChat)
	}
	chatSum := strings.Join(recentChats, " ")
	regexpKorean := regexp.MustCompile("[^가-힣]+")
	onlyKorean := regexpKorean.ReplaceAllString(chatSum, " ")
	fmt.Println(onlyKorean)
}

// type AnswerData struct {
// 	QuestionContents string `json:"question_contents"`
// 	FirstAnswer string `json:"first_answer"`
// 	SecondAnswer string `json:"second_answer"`
// 	AnswerDate string `json:"answer_date"`
// }
// TEST
// TEST
// TEST

func TestUsrs() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM usrs")
	return r,err
}

func TestChat() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM chat")
	return r,err
}

func TestRequest() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM request")
	return r,err
}

func TestConnection() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM connection")
	return r,err
}

func TestQuestion() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM question")
	return r,err
}

func TestAnswer() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM answer")
	return r,err
}

func DeleteAll(){
	_, _ = db.Query("DELETE FROM usrs")
	_, _ = db.Query("DELETE FROM chat")
	_, _ = db.Query("DELETE FROM request")
	_, _ = db.Query("DELETE FROM connection")
	_, _ = db.Query("DELETE FROM answer")
	// _,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("강아지", "강아지와 고양이 중 뭐가 더 좋아?")`)
	// _,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("운동", "운동하는 거 좋아해?")`)
	// _,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("남사친", "남사친/여사친 어디까지 허용 가능하다!")`)
	// _,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("엄마", "부모님께 존댓말 써?")`)
	// _,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("결혼", "결혼은 언제쯤 하고싶어?")`)

}