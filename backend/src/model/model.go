package model

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ChatData struct {
	Text_body string `json:"text_body"`
	Writer_id string `json:"writer_id"`
	Write_time string `json:"write_time"`
	Is_answer int `json:"is_answer"`
	Chat_id int `json:"chat_id"`
	Question_id int `json:"question_id"`
	Is_deleted int `json:"is_deleted"`
	Is_file int `json:"is_file"`
	Is_image int `json:"is_image"`
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
	Order int `json:"order"`
}

type BeAboutToDeleteData struct {
	Delete_Date string
	Connection_id int
}

type AnniversaryData struct {
	Anniversary_id int `json:"anniversary_id"`
	Connection_id int
	Year int `json:"year"`
	Month int `json:"month"`
	Date int `json:"date"`
	Contents string `json:"contents"`
	D_day bool `json:"d_day"`
}

var db *sql.DB

func OpenDB(driverName, dataSourceName string) error {
	database, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	db = database

	// DB와 서버가 연결 되었는지 확인
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	err := db.Close()
	return err
}

// 쿠키가 있는지 확인
func CookieExist(c *gin.Context) (string, error) {
	uuid, err := c.Cookie("uuid")	
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func InsertUsr(id, password, uuid string) error {
	_, err := db.Query(`INSERT INTO usrs (id, password, uuid, conn_id) VALUES ("`+id+`", "`+password+`", "`+uuid+`", 0)`)
	return err
}


func CheckUsrByID(id string) (bool, error) {
	r, err := db.Query(`SELECT * FROM usrs WHERE id = "`+id+`"`)
	defer r.Close()
	if err != nil {
		return false, err
	}

	if r.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func GetUUIDByIDandPW(id, password string) (string, error) {
	r, err := db.Query(`SELECT uuid FROM usrs WHERE id = "`+id+`" and password = "`+password+`"`)
	defer r.Close()
	if err != nil {
		return "", err
	}

	if r.Next() {
		var uuid string
		r.Scan(&uuid)
		return uuid, nil
	} else {
		return "", nil
	}
}

func SelectConnIDFromUsrsByUUID(uuid string) (int, error) {
	r, err := db.Query(`SELECT conn_id FROM usrs WHERE uuid = "`+uuid+`"`)
	defer r.Close()
	if err != nil {
		return 0, err
	}

	var conn_id int
	for r.Next() {
		r.Scan(&conn_id)
	}
	return conn_id, nil
}

func CheckRequestByRequesterUUID(uuid string) (bool, error) {
	r, err := db.Query(`SELECT * FROM request WHERE requester_uuid = "`+uuid+`"`)
	defer r.Close()
	if err != nil {
		return false, err
	}
	
	if r.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func SelectIDFromUsrsByUUID(uuid string) (string, error) {
	r, err1 := db.Query(`SELECT id FROM usrs WHERE uuid = "`+uuid+`"`)
	defer r.Close()
	if err1 != nil {
		return "", err1
	}

	var id string
	r.Next()
	err2 := r.Scan(&id)
	if err2 != nil {
		return "", err2
	}

	return id, nil
}

func SelectConnIDandUUIDFromUsrsByID(id string) (bool, int, string, error) {
	targetConnID := 0
	targetUUID := ""

	r, err1 := db.Query(`SELECT conn_id, uuid FROM usrs WHERE id = "`+id+`"`)
	defer r.Close()
	if err1 != nil {
		return false, targetConnID, targetUUID, err1
	}
	
	// ID가 존재하는 ID면 이미 연결되어있진 않은지 conn_id를 확인
	if r.Next() {	
		err2 := r.Scan(&targetConnID, &targetUUID)
		if err2 != nil {
			return false, targetConnID, targetUUID, err2
		} else {
			return true, targetConnID, targetUUID, nil
		}
	}
	// ID가 존재하지 않는 ID면
	return false, targetConnID, targetUUID, nil
}

func InsertRequest(requester_uuid, target_uuid, request_time, requester_id, target_id  string) error {
	_, err := db.Query(`INSERT INTO request (requester_uuid, target_uuid, request_time, requester_id, target_id) VALUES ("`+requester_uuid+`", "`+target_uuid+`", "`+request_time+`", "`+requester_id+`", "`+target_id+`")`)
	return err
}

func SelectRecieveRequestByTargetUUID(uuid string) ([]RequestData, error) {
	requestedData := RequestData{}
	requestedDatas := []RequestData{}

	r, err := db.Query(`SELECT requester_id, requester_uuid, request_time, request_id FROM request WHERE target_uuid = "`+uuid+`"`)
	defer r.Close()
	if err != nil {
		return nil, err
	}

	for r.Next() {
		r.Scan(&requestedData.Requester_id, &requestedData.Requester_uuid, &requestedData.Request_time, &requestedData.Request_id)
		requestedDatas = append(requestedDatas, requestedData)
	}
	return requestedDatas, nil
}

func SelectSendRequestByTargetUUID(uuid string) (RequestData, error) {
	requestingData := RequestData{}

	r, err := db.Query(`SELECT target_uuid, request_time, target_id FROM request WHERE requester_uuid = "`+uuid+`"`)
	defer r.Close()
	if err != nil {
		return requestingData, err
	}
	
	for r.Next(){
		r.Scan(&requestingData.Target_uuid, &requestingData.Request_time, &requestingData.Target_id)
	}
	return requestingData, nil
}

func InsertConnection(first_usr, second_usr, start_date string) error {
	_, err := db.Query(`INSERT INTO connection (first_usr, second_usr, start_date) VALUES ("`+first_usr+`", "`+second_usr+`", "`+start_date+`")`)
	return err
}

func SelectConnectionIDByUsrsUUID(first_usr, second_usr string) (int, error) {
	r, err1 := db.Query(`SELECT connection_id FROM connection WHERE first_usr = "`+first_usr+`" and second_usr = "`+second_usr+`"`)
	defer r.Close()
	if err1 != nil {
		return 0, err1
	}

	var connID int
	r.Next()
	err2 := r.Scan(&connID)
	if err2 != nil {
		return 0, err2
	}

	return connID, nil
}

func UpdateUsrsConnID(conn_id int, targetUUID string) error {
	_, err := db.Query(`UPDATE usrs SET order_usr = 1, conn_id = `+strconv.Itoa(conn_id)+` WHERE uuid = "`+targetUUID+`"`)
	return err
}

func UpdateUsrsOrder(conn_id int, myUUID string) error {
	_, err := db.Query(`UPDATE usrs SET order_usr = 2, conn_id = `+strconv.Itoa(conn_id)+` WHERE uuid = "`+myUUID+`"`)
	return err
}

func DeleteRestRequest(requester_uuid, target_uuid string) error {
	_, err := db.Query(`DELETE FROM request WHERE requester_uuid = "`+requester_uuid+`" or target_uuid = "`+requester_uuid+`" or requester_uuid = "`+target_uuid+`" or target_uuid = "`+target_uuid+`"`)
	return err
}

func DeleteRequestByRequestID(request_id string) error {
	_, err := db.Query(`DELETE FROM request WHERE request_id = `+request_id)
	return err
}

func SelectConnIDByUUID(uuid string) (int, error) {
	r, err1 := db.Query(`SELECT conn_id FROM usrs WHERE uuid = "`+uuid+`"`)
	defer r.Close()
	if err1 != nil {
		return 0, err1
	}

	var connID int
	r.Next()
	err2 := r.Scan(&connID)
	if err2 != nil {
		return 0, err2
	}
	return connID, nil
}

func GetAnswerandQuestionContentsByConnIDWithOrder(connection_id, order int) ([]AnswerData, error) {
	r, err1 := db.Query(`SELECT first_answer, second_answer, answer_date, question_id FROM answer WHERE connection_id = "`+strconv.Itoa(connection_id)+`"`)
	defer r.Close()
	if err1 != nil {
		return nil, err1
	}

	answerData := AnswerData{}
	answerDatas := []AnswerData{}
	
	for r.Next() {
		err2 := r.Scan(&answerData.FirstAnswer, &answerData.SecondAnswer, &answerData.AnswerDate, &answerData.Question_id)
		if err2 != nil {
			return nil, err2
		}
		questionContents, err3 := selectQuestionContentsByQuestionID(answerData.Question_id)
		if err3 != nil {
			return nil, err3
		}
		answerData.QuestionContents = questionContents

		if answerData.FirstAnswer == "not-written" || answerData.SecondAnswer == "not-written" {
			continue
		}
		answerData.Order = order
		answerDatas = append(answerDatas, answerData)
	}
	return answerDatas, nil
}

func selectQuestionContentsByQuestionID(question_id int) (string, error) {
	r, err := db.Query(`SELECT question_contents FROM question WHERE question_id = `+strconv.Itoa(question_id))
	defer r.Close()
	if err != nil {
		return "", err
	}

	var question_contents string

	r.Next()
	r.Scan(&question_contents)
	
	return question_contents, nil
}

func GetConnectionByUsrsUUID(uuid string) (string, string, int, error) {
	r, err1 := db.Query(`SELECT first_usr, second_usr, connection_id FROM connection WHERE first_usr = "`+uuid+`" or second_usr = "`+uuid+`"`)
	defer r.Close()
	if err1 != nil {
		return "", "", 0, err1
	}

	var first_uuid, second_uuid string
	var conn_id int

	r.Next()
	err2 := r.Scan(&first_uuid, &second_uuid, &conn_id)
	if err2 != nil {
		return "", "", 0, err2
	}
	return first_uuid, second_uuid, conn_id, nil
}

func SelectChatByUsrsUUID(first_uuid, second_uuid string) ([]ChatData, error) {
	initialChat := ChatData{}
	initialChats := []ChatData{}

	r, err := db.Query(`SELECT chat_id, writer_id, write_time, text_body, is_file, is_image FROM chat WHERE writer_id = "`+first_uuid+`" or writer_id = "`+second_uuid+`" ORDER BY chat_id ASC`)
	defer r.Close()
	if err != nil {
		return nil, err
	}

	for r.Next() {
		r.Scan(&initialChat.Chat_id, &initialChat.Writer_id, &initialChat.Write_time, &initialChat.Text_body, &initialChat.Is_file, &initialChat.Is_image)
		initialChat.Is_deleted = 0
		initialChat.Is_answer = 0
		initialChats = append(initialChats, initialChat)		
	}
	return initialChats, nil
}

func InsertChatAndGetChatID(text_body, writer_id, write_time string, is_file, is_image int) (int, error) {
	_, err1 := db.Query(`INSERT INTO chat (text_body, writer_id, write_time, is_file, is_image) VALUES ("`+text_body+`", "`+writer_id+`", "`+write_time+`", "`+strconv.Itoa(is_file)+`", "`+strconv.Itoa(is_image)+`")`)
	if err1 != nil {
		return 0, err1
	}

	r, err2 := db.Query("SELECT chat_id FROM chat ORDER BY chat_id DESC LIMIT 1")
	defer r.Close()
	if err2 != nil {
		return 0, err2
	}

	var chat_id int
	r.Next()
	r.Scan(&chat_id)
	
	return chat_id, nil
}

func CheckAnswerByConnIDandQuestionID(connection_id, question_id int) (bool, error) {
	r, err := db.Query(`SELECT * FROM answer WHERE connection_id = `+strconv.Itoa(connection_id)+` and question_id = `+strconv.Itoa(question_id))
	defer r.Close()
	if err != nil {
		return false, err
	}

	if r.Next() {
		return true, nil
	}
	return false, nil
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

func GetUsrOrderByUUID(uuid string) (int, error) {
	r, err := db.Query(`SELECT order_usr FROM usrs WHERE uuid = "`+uuid+`"`)
	defer r.Close()
	if err != nil {
		return 0, err
	}

	var order_usr int
	r.Next()
	r.Scan(&order_usr)
	return order_usr, nil
}

func QuestionIDOfEmptyAnswerByOrder(order, connection_id int) (int, error) {
	var question_id int

	if order == 1 {
		r, err1 := db.Query(`SELECT question_id FROM ANSWER WHERE first_answer = "not-written" and connection_id = `+strconv.Itoa(connection_id))
		defer r.Close()
		if err1 != nil {
			return 0, err1
		}
		if r.Next() {
			r.Scan(&question_id)
		} else {
			return 0, nil
		}
	} else {
		r, err2 := db.Query(`SELECT question_id FROM ANSWER WHERE second_answer = "not-written" and connection_id = `+strconv.Itoa(connection_id))
		defer r.Close()

		if err2 != nil {
			return 0, err2
		}
		if r.Next() {
			r.Scan(&question_id)
		} else {
			return 0, nil
		}
	}
	return question_id, nil
}

func GetQuestionByQuestionID(questionID int) (string, string, error){
	var questionData QuestionData

	r, err1 := db.Query(`SELECT target_word, question_contents FROM question WHERE question_id = `+ strconv.Itoa(questionID))
	defer r.Close()
	if err1 != nil {
		return "", "", err1
	}
	
	r.Next()
	err2 := r.Scan(&questionData.Target_word, &questionData.Question_contents)
	if err2 != nil {
		return "", "", err2
	}

	return questionData.Target_word, questionData.Question_contents, nil
}

func GetRecentAnswerByConnID(connection_id, num int) []AnswerData {
	r, err := db.Query("SELECT * FROM answer ORDER BY answer_id DESC LIMIT " + strconv.Itoa(num))
	defer r.Close()
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

func GetFrequentWords(uuid string, rankNum int) ([]string, error) {
	r, err := db.Query(`SELECT text_body FROM chat WHERE writer_id = "`+uuid+`" and DATE_ADD(NOW(), INTERVAL -7 DAY) < write_time`)
	defer r.Close()
	if err != nil {
		fmt.Println("ERROR #56 : ", err.Error())
	}

	var recentChat string
	var recentChats []string
	for r.Next() {
		r.Scan(&recentChat)
		recentChats = append(recentChats, recentChat)
	}

	chatSum := strings.Join(recentChats, " ")
	regexpKorean := regexp.MustCompile("[^가-힣]+")
	onlyKorean := regexpKorean.ReplaceAllString(chatSum, " ")

	addWordTimes := make(map[string]int)

	wordSlice := strings.Split(onlyKorean, " ")

	for i := 0; i < len(wordSlice); i++ {
		addWordTimes[wordSlice[i]] = addWordTimes[wordSlice[i]] + 1
	}

	var withOutRepeat string

	for i := 0; i < len(wordSlice); i++ {
		if !strings.Contains(withOutRepeat, wordSlice[i]) {
			withOutRepeat += " "+wordSlice[i]
		}
	}

	conn_id, err := SelectConnIDByUUID(uuid)

	exceptWordsSlice, err := GetExceptWords(conn_id)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(exceptWordsSlice); i++ {
		withOutRepeat = strings.ReplaceAll(withOutRepeat, exceptWordsSlice[i], "")
	}

	withOutRepeatSlice := strings.Fields(withOutRepeat)

	if len(withOutRepeatSlice) < rankNum {
		return nil, err
	}

	// bubble sort
	for i := 0; i < len(withOutRepeatSlice)-1; i++ {
		for i := len(withOutRepeatSlice)-1; i > 0; i-- {
			if addWordTimes[withOutRepeatSlice[i]] > addWordTimes[withOutRepeatSlice[i-1]] {
				temp := withOutRepeatSlice[i-1]
				withOutRepeatSlice[i-1] = withOutRepeatSlice[i]
				withOutRepeatSlice[i] = temp
			}
		}
	}
	return withOutRepeatSlice[:rankNum], nil
}

func InsertExceptWord(connection_id int, except_word string) error {
	_, err := db.Query(`INSERT INTO exceptionword (connection_id, except_word) VALUES (`+strconv.Itoa(connection_id)+`, "`+except_word+`")`)
	return err
}

func CheckWordAlreadyExcepted(connection_id int, except_word string) (bool, error) {
	r, err := db.Query(`SELECT * FROM exceptionword WHERE connection_id = `+strconv.Itoa(connection_id)+` and except_word = "`+except_word+`"`)
	defer r.Close()

	if r.Next() {
		return true, nil
	} else {
		return false, err
	}
}

func CancleExceptWord(connection_id int, except_word string) error {
	_, err := db.Query(`DELETE FROM exceptionword WHERE connection_id = `+strconv.Itoa(connection_id)+` and except_word = "`+except_word+`"`)
	return err
}

func GetExceptWords(connection_id int) ([]string, error) {
	r, err := db.Query("SELECT except_word FROM exceptionword WHERE connection_id = "+strconv.Itoa(connection_id))
	defer r.Close()
	if err != nil {
		return nil, err
	}

	var exceptWord string
	var exceptWords []string
	for r.Next() {
		r.Scan(&exceptWord)
		exceptWords = append(exceptWords, exceptWord)
	}
	return exceptWords, nil
}

func DeleteUsrByUUID(uuid string) error {
	_, err := db.Query(`DELETE FROM usrs WHERE uuid = "`+uuid+`"`)
	return err
}

func InsertBeAboutToDelete(connection_id int) error {
	_, err := db.Query("INSERT INTO beabouttodelete (connection_id) VALUES ("+strconv.Itoa(connection_id)+")")
	return err
}

func DeleteConnectionByConnID(first_uuid, second_uuid string, conn_id int) error {
	_, err := db.Query(`DELETE FROM chat WHERE writer_id = "`+first_uuid+`" or writer_id = "`+second_uuid+`"`)
	_, err = db.Query("DELETE FROM connection WHERE connection_id = "+strconv.Itoa(conn_id))
	_, err = db.Query("DELETE FROM beaboutdelete WHERE connection_id = "+strconv.Itoa(conn_id))
	_, err = db.Query("DELETE FROM answer WHERE connection_id = "+strconv.Itoa(conn_id))
	_, err = db.Query("DELETE FROM exceptionword WHERE connection_id = "+strconv.Itoa(conn_id))
	_, err = db.Query("DELETE FROM anniversary WHERE connection_id = "+strconv.Itoa(conn_id))
	_, err = db.Query(`UPDATE usrs SET conn_id = 0  WHERE uuid = "`+first_uuid+`" or uuid = "`+second_uuid+`"`)
	_, err = db.Query(`UPDATE usrs SET order_usr = 0 WHERE uuid = "`+first_uuid+`" or uuid = "`+second_uuid+`"`)

	return err
}

func ChangePassword(password, uuid string) error {
	_, err := db.Query(`UPDATE usrs SET password = "`+password+`" WHERE uuid = "`+uuid+`"`)
	return err
}

func DeleteChatByChatID(chat_id int) error {
	_, err := db.Query("DELETE FROM chat WHERE chat_id = "+ strconv.Itoa(chat_id))
	return err
}

func InsertAnniversaryByConnID(data AnniversaryData) error {
	_, err := db.Query(`INSERT INTO anniversary (connection_id, year, month, date, contents, d_day) Values (`+strconv.Itoa(data.Connection_id)+`, `+strconv.Itoa(data.Year)+`, `+strconv.Itoa(data.Month)+`, `+strconv.Itoa(data.Date)+`, "`+data.Contents+`", `+strconv.FormatBool(data.D_day)+`)`)
	return err
}

func GetAnniversaryByConnIDAndMonthAndYear(connection_id int, target_month, target_year string) ([]AnniversaryData, error) {
	r, err := db.Query(`SELECT * FROM anniversary WHERE connection_id = `+strconv.Itoa(connection_id)+` and month = "`+target_month+`" and year = "`+target_year+`"`)
	defer r.Close()
	if err != nil {
		return nil, err
	}

	var anniversaryData AnniversaryData
	var anniversaryDatas []AnniversaryData

	for r.Next() {
		r.Scan(&anniversaryData.Anniversary_id, &anniversaryData.Connection_id, &anniversaryData.Year, &anniversaryData.Month, &anniversaryData.Date, &anniversaryData.Contents, &anniversaryData.D_day)
		anniversaryDatas = append(anniversaryDatas, anniversaryData)
	}

	return anniversaryDatas, nil
}

func DeleteAnniversaryByAnniversaryID(anniversary_id string) error {
	_, err := db.Query("DELETE FROM anniversary WHERE anniversary_id = "+anniversary_id)
	return err
}

func GetDDayAnniversaryIDByConnID(connection_id int) (int, error) {
r, err := db.Query("SELECT anniversary_id FROM anniversary WHERE d_day = 1 and connection_id = "+strconv.Itoa(connection_id))
	if err != nil {
		return 0, err
	}

	var anniversary_id int
	if r.Next() {
		r.Scan(&anniversary_id)
		return anniversary_id, nil
	}

	return 0, nil
}

func ChangeDDayZeroByAnniversaryID(anniversary_id int) error {
	_, err := db.Query("UPDATE anniversary SET d_day = 0 WHERE anniversary_id = "+strconv.Itoa(anniversary_id))
	return err
}

func GetDDayByConnID(connection_id int) ([]AnniversaryData, error){
	r, err := db.Query("SELECT * FROM anniversary WHERE d_day = 1 and connection_id = "+strconv.Itoa(connection_id))
	if err != nil {
		return nil, err
	}

	var anniversaryData []AnniversaryData
	var tempData AnniversaryData
	if r.Next(){
		r.Scan(&tempData.Anniversary_id, &tempData.Connection_id, &tempData.Year, &tempData.Month, &tempData.Date, &tempData.Contents, &tempData.D_day)
		anniversaryData = append(anniversaryData, tempData)
		return anniversaryData, nil
	}

	return nil, nil
}

func GetChatIDFromRecentFileChatByUUID(uuid string) (int, error) {
	r, err := db.Query(`SELECT chat_id FROM chat WHERE writer_id = "`+uuid+`" and is_file = 1 ORDER BY chat_id DESC LIMIT 1`)

	var chatID int

	if r.Next() {
		r.Scan(&chatID)
		return chatID, nil
	}
	
	return 0, err	
}

func GetTextBodyByChatID(chat_id int) (string, error){
	r, err := db.Query("SELECT text_body FROM chat WHERE chat_id = "+strconv.Itoa(chat_id))
	if err != nil {		return "", err
	}
	var text_data string
	r.Next()
	r.Scan(&text_data)
	return text_data, nil
}

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

func TestExceptionWord() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM exceptionword")
	return r,err
}

func TestBeAboutToDelete() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM beabouttodelete")
	return r,err
}

func TestAnniversary() (*sql.Rows, error) {
	r, err := db.Query("SELECT * FROM Anniversary")
	return r,err
}

func DeleteAll(){
	// _, _ = db.Query("DELETE FROM usrs")
	// _, _ = db.Query("DELETE FROM chat")
	// _, _ = db.Query("DELETE FROM request")
	// _, _ = db.Query("DELETE FROM connection")
	// _, _ = db.Query("DELETE FROM answer")
	// _, _ = db.Query("DELETE FROM exceptionword")
	// _, _ = db.Query("DELETE FROM anniversary")
	// _, _ = db.Query("DELETE FROM question")
	_,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("강아지", "강아지와 고양이 중 뭐가 더 좋아?")`)
	_,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("운동", "운동하는 거 좋아해?")`)
	_,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("남사친", "남사친/여사친 어디까지 허용 가능하다!")`)
	_,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("엄마", "부모님께 존댓말 써?")`)
	_,_=db.Query(`INSERT INTO QUESTION (target_word, question_contents) VALUES ("결혼", "결혼은 언제쯤 하고싶어?")`)
}