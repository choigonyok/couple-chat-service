import REACT, { useEffect, useRef, useState } from "react";
// NULL을 사용하려면 REACT를 import 해줘야함
import "./Chatpage.css";
import Inputbox from "./Inputbox";
import Logout from "./Logout";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import Exceptword from "./Exceptword";
import Withdrawal from "./Withdrawal";
import Cutconn from "./Cutconn";
import Changepw from "./Changepw";
import Searchword from "./Searchword";
import Calender from "./Calender";
import Answer from "./Answer";
import Canclecutconn from "./Canclecutconn";

const Chatpage = () => {
  const flexboxRef = useRef({});

  const navigator = useNavigate();

  const [newSocket, setNewSocket] = useState(null);
  const [chatDate, setChatDate] = useState([]);
  const [recievedMessage, setRecievedMessage] = useState([]);
  const [myUUID, setMyUUID] = useState("");
  const [inputAnswer, setInputAnswer] = useState("");
  const [hideInputBox, setHideInputBox] = useState(false);
  const [seeAnswerBox, setSeeAnswerBox] = useState(false);
  const [searchButton, setSearchButton] = useState(false);
  const [answerButton, setAnswerButton] = useState(false);
  const [rankingButton, setRankingButton] = useState(false);
  const [calenderButton, setCalenderButton] = useState(false);
  const [fileClick, setFileClick] = useState(0);
  const [highLight, setHighLight] = useState(0);
  const [latestChatID, setLatestChatID] = useState(0);
  

  useEffect (()=>{
    setHighLight(0);
  },[searchButton])

  const [chatID, setChatID] = useState(0);

  const onChangeScroll = (id) => {
    const tags = flexboxRef.current.querySelectorAll('[id^="tag"]');
    const scrollPositions = {};
    let sum = 18;
    tags.forEach((tag) => {
      scrollPositions[tag.id] = sum;
      sum = sum + tag.scrollHeight;
    });
    if (flexboxRef.current) {
      flexboxRef.current.scrollTop = scrollPositions['tag'+id];
      setHighLight(id);
    }    
  }

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/log")
      .then((response) => {
        if (response.data === "CONNECTED") {
        } else if (response.data === "NOT_CONNECTED") {
          navigator("/conn");
        }
      })
      .catch((error) => {
        if (error.response.status === 401) {
          navigator("/");
        } else {
          console.log(error);
        }
      });
  }, []);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socket.onopen = () => {
      console.log("CONNECTION CONNECTED");
      setNewSocket(socket);
    };

    socket.onmessage = (e) => {
      const parsedData = JSON.parse(e.data);
      if (parsedData.uuid) {
        setMyUUID(parsedData.uuid);
      } else {
        if (parsedData[0].is_deleted === 1) {
          setChatID(parsedData[0].chat_id);
        } else {
          setRecievedMessage((prev) => [...prev, ...parsedData]);
          setLatestChatID(parsedData[(parsedData.length)-1].chat_id);
          if (parsedData.length !== 0) {
            if (parsedData[0].is_answer === 1) {
              setHideInputBox(true);
              setSeeAnswerBox(true);
            }
          }
        }
      }
    };
    socket.onclose = () => {
      console.log("CONNECTION CLOSED");
    };

    return () => {
      socket.close();
    };
  }, []);

  useEffect(() => {
    setChatDate([]);
    for (let i = 0; i < recievedMessage.length; i++) {
      if (recievedMessage[i].chat_id !== 0) {
        const writeTime = recievedMessage[i].write_time;
        const [dateString, timeString] = writeTime.split(" ");
        const [year, month, date] = dateString.split("-");
        const [hour, minute, second] = timeString.split(":");
        setChatDate((prev) => [
          ...prev,
          year + "/" + month + "/" + date + " " + hour + ":" + minute,
        ]);
      }
    }
  }, [recievedMessage]);

  useEffect(()=>{
    onChangeScroll(latestChatID);
    setHighLight(0);
  },[chatDate])

  useEffect(() => {
    const deletedArray = recievedMessage.filter((m) => m.chat_id !== chatID);
    setRecievedMessage(deletedArray);
  }, [chatID]);

  const sendMessageHandler = (data) => {
    if (newSocket !== null) {
      const now = new Date();
      const nowMonth =
        now.getMonth() + 1 < 10
          ? "0" + (now.getMonth() + 1)
          : now.getMonth() + 1;
      const nowDate = now.getDate() < 10 ? "0" + now.getDate() : now.getDate();
      const nowHour =
        now.getHours() < 10 ? "0" + now.getHours() : now.getHours();
      const nowMinute =
        now.getMinutes() < 10 ? "0" + now.getMinutes() : now.getMinutes();
      const nowSecond =
        now.getSeconds() < 10 ? "0" + now.getSeconds() : now.getSeconds();

      let nowformat =
        now.getFullYear() +
        "-" +
        nowMonth +
        "-" +
        nowDate +
        " " +
        nowHour +
        ":" +
        nowMinute +
        ":" +
        nowSecond;

      // 파일 전송시
      // 0 : 이미지
      if (data === 0) {
        const sendData = [
          {
            text_body: "",
            write_time: nowformat,
            writer_id: myUUID,
            is_answer: 0,
            is_deleted: 0,
            is_file: 1,
            is_image: 1,
          },
        ];
        newSocket.send(JSON.stringify(sendData));
      } else if (data === 1) {
        const sendData = [
          {
            text_body: "",
            write_time: nowformat,
            writer_id: myUUID,
            is_answer: 0,
            is_deleted: 0,
            is_file: 1,
            is_image: 0,
          },
        ];
        newSocket.send(JSON.stringify(sendData));
      } else {
        const sendData = [
          {
            text_body: String(data),
            write_time: nowformat,
            writer_id: myUUID,
            is_answer: 0,
            is_deleted: 0,
            is_file: 0,
          },
        ];
        newSocket.send(JSON.stringify(sendData));
      }
    } else alert("상대방에게 메세지를 보낼 수 없는 상태입니다.");
  };

  const inputAnswerHandler = (e) => {
    setInputAnswer(e.target.value);
  };

  const deleteChatHandler = (value) => {
    const sendData = [
      {
        is_file: value.is_file,
        chat_id: value.chat_id,
        is_deleted: 1,
      },
    ];
    newSocket.send(JSON.stringify(sendData));
  };

  const enterHandler = (item, e) => {
    if (e.key === "Enter") {
      if (inputAnswer !== "") {
        if (e.nativeEvent.isComposing) {
          return;
        } else {
          if (newSocket !== null) {
            let now = new Date();
            let nowformat =
              now.getFullYear() +
              "-" +
              now.getMonth() +
              2 +
              "-" +
              now.getDate() +
              " " +
              now.getHours() +
              ":" +
              now.getMinutes();
            const sendData = [
              {
                text_body: inputAnswer,
                write_time: nowformat,
                writer_id: myUUID,
                is_answer: item.is_answer,
                question_id: item.question_id,
              },
            ];
            newSocket.send(JSON.stringify(sendData));
            setInputAnswer("");
            setHideInputBox(false);
            setSeeAnswerBox(false);
          } else alert("상대방에게 메세지를 보낼 수 없는 상태입니다.");
        }
      }
    }
  };

  const searchClickHandler = () => {
    setSearchButton(!searchButton);
    setAnswerButton(false);
    setRankingButton(false);
    setCalenderButton(false);
  };
  const answerClickHandler = () => {
    setAnswerButton(!answerButton);
    setSearchButton(false);
    setRankingButton(false);
    setCalenderButton(false);
  };
  const rankingClickHandler = () => {
    setRankingButton(!rankingButton);
    setSearchButton(false);
    setAnswerButton(false);
    setCalenderButton(false);
  };
  const calenderClickHandler = () => {
    setCalenderButton(!calenderButton);
    setSearchButton(false);
    setAnswerButton(false);
    setRankingButton(false);
  };

  const fileClickHandler = (value) => {
    setFileClick(value.chat_id);
  };

  const resetFileClickHandler = () => {
    setFileClick(0);
  };

  const saveFileHandler = (value) => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/file/name/" + value)
      .then((response) => {
        const link = document.createElement("a");
        link.href =
          process.env.REACT_APP_HOST_URL + "/api/file/" + value;
        link.download = response.data.filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      })
      .catch((error) => {
        alert("저장 오류 발생");
        return;
      });
  };

  const saveImageHandler = () => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/file/name/" + fileClick)
      .then((response) => {
        const link = document.createElement("a");
        link.href =
          process.env.REACT_APP_HOST_URL + "/api/file/" + fileClick;
        link.download = response.data.filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      })
      .catch((error) => {
        alert("저장 오류 발생");
        return;
      });
  };

  return (
    <div className="page-container">
      {fileClick !== 0 && (
        <div>
          <div>
            <img
              src={
                process.env.REACT_APP_HOST_URL + "/api/file/" + fileClick
              }
            />
          </div>
          <div>
            <input type="button" value="X" onClick={resetFileClickHandler} />
            <input type="button" value="저장하기" onClick={saveImageHandler} />
          </div>
        </div>
      )}
      <div className="button-container">
        <div>
          <input
            type="button"
            value="SEARCH"
            className="buttons"
            onClick={searchClickHandler}
          />
        </div>

        <div>
          <input
            type="button"
            value="ANSWER"
            className="buttons"
            onClick={answerClickHandler}
          />
        </div>

        <div>
          <input
            type="button"
            value="RANK"
            className="buttons"
            onClick={rankingClickHandler}
          />
        </div>
        <div>
          <input
            type="button"
            value="CALENDER"
            className="buttons"
            onClick={calenderClickHandler}
          />
        </div>
      </div>
      <div>
        <div>{searchButton && <Searchword onWordSearch={onChangeScroll}/>}</div>
        <div>{answerButton && <Answer />}</div>
        <div>{rankingButton && <Exceptword />}</div>
        <div>{calenderButton && <Calender />}</div>
      </div>
      <div className="chat-container" ref={flexboxRef}>
        {recievedMessage &&
          recievedMessage.map((item, index) => (
            <div className="chat-box" id={"tag"+item.chat_id}>
              {item.is_answer === 0 && item.writer_id === myUUID && (
                <div className="chat-container__chat__usr" >
                  <div className="chat-container__chatandbutton">
                    <div className={item.is_file === 0 && highLight === item.chat_id ? "chat__usr__highlight" : "chat__usr"}>
                      {item.is_file === 0 ? (
                        item.text_body
                      ) : item.is_image === 1 ? (
                        <img
                          src={
                            process.env.REACT_APP_HOST_URL +
                            "/api/file/" +
                            item.chat_id
                          }
                          className="chat__image"
                          onClick={() => fileClickHandler(item)}
                        />
                      ) : (
                        <div className="chat__file" onClick={()=>saveFileHandler(item.chat_id)}>
                          <div>{item.text_body}</div>
                        </div>
                      )}
                    </div>
                    <div>
                      <input
                        type="button"
                        value="X"
                        className="chat__button"
                        onClick={() => deleteChatHandler(item)}
                      />
                    </div>
                  </div>
                  <div className="chat__time">{chatDate[index]}</div>
                </div>
              )}
              {item.is_answer === 0 && item.writer_id !== myUUID && (
                <div className="chat-container__chat__other">
                  <div>
                    <div className={item.is_file === 0 && highLight === item.chat_id ? "chat__other__highlight" : "chat__other"}>
                      {item.is_file === 0 ? (
                        item.text_body
                      ) : item.is_image === 1 ? (
                        <img
                          src={
                            process.env.REACT_APP_HOST_URL +
                            "/api/file/" +
                            item.chat_id
                          }
                          className="chat__image"
                          onClick={() => fileClickHandler(item)}
                        />
                      ) : (
                        <div className="chat__file" onClick={()=>saveFileHandler(item.chat_id)}>
                          <div>{item.text_body}</div>
                        </div>
                      )}
                    </div>
                  </div>
                  <div className="chat__time">{chatDate[index]}</div>
                </div>
              )}
              {seeAnswerBox && item.is_answer === 1 && (
                <div>
                  <div className="chat-container__question">
                    <div className="chat__question">{item.text_body}</div>
                  </div>
                  <div className="chat-container__answer">
                    <input
                      type="text"
                      placeholder="답을 작성해주세요."
                      autofocus
                      className="chat__answer"
                      value={inputAnswer}
                      onChange={inputAnswerHandler}
                      onKeyDownCapture={(e) => enterHandler(item, e)}
                    />
                  </div>
                </div>
              )}
            </div>
          ))}
      </div>
      <div>
        {!hideInputBox && (
          <Inputbox
            onSendMessage={(messageData) => sendMessageHandler(messageData)}
          />
        )}
      </div>
      <div className="usr-button">
        <Changepw />
        <Logout />
        <Withdrawal />
        <Cutconn />
        <Canclecutconn />
      </div>
    </div>
  );
};

export default Chatpage;
