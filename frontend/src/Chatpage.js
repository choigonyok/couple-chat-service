import REACT, { useEffect, useState } from "react";
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

const Chatpage = () => {
  const navigator = useNavigate();

  const [newSocket, setNewSocket] = useState(null);
  const [recievedMessage, setRecievedMessage] = useState([]);
  const [myUUID, setMyUUID] = useState("");
  const [inputAnswer, setInputAnswer] = useState("");
  const [hideInputBox, setHideInputBox] = useState(false);
  const [seeAnswerBox, setSeeAnswerBox] = useState(false);
  const [answers, setAnswers] = useState([]);
  const [wordNum, setWordNum] = useState("3");
  const [otherWords, setOtherWords] = useState([]);
  const [myWords, setMyWords] = useState([]);
  
  const [chatID, setChatID] = useState(0);

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

  useEffect (()=>{
    const deletedArray = recievedMessage.filter(m => m.chat_id !== chatID);
    setRecievedMessage(deletedArray);
  },[chatID])

  const resetDeleted = () => {
    setdete
  }

  const sendMessageHandler = (data) => {
    if (newSocket !== null) {
      let now = new Date();
      let nowformat =
        now.getFullYear() +
        "/" +
        (now.getMonth() + 1) +
        "/" +
        now.getDate() +
        " " +
        now.getHours() +
        ":" +
        now.getMinutes();
      const sendData = [
        {
          text_body: String(data),
          write_time: nowformat,
          writer_id: myUUID,
          is_answer: 0,
          is_deleted: 0,
        },
      ];
      newSocket.send(JSON.stringify(sendData));
    } else alert("상대방에게 메세지를 보낼 수 없는 상태입니다.");
  };

  const inputAnswerHandler = (e) => {
    setInputAnswer(e.target.value);
  };

  const deleteChatHandler = (value) => {
    const sendData = [
      {
        chat_id: value,
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

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/answer")
      .then((response) => {
        setAnswers([...response.data]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/rank/" + wordNum)
      .then((response) => {
        setMyWords([...response.data.mywords]);
        setOtherWords([...response.data.otherwords]);
      })
      .catch((error) => {
        if (error.response.status === 411) {
          alert("순위를 매기기엔 채팅의 수가 부족합니다.");
        } else {
          console.log(error);
        }
      });
  }, [wordNum]);

  const threeWordsHandler = () => {
    setWordNum("3");
  };
  const fiveWordsHandler = () => {
    setWordNum("5");
  };
  const tenWordsHandler = () => {
    setWordNum("10");
  };

  return (
    <div className="page-container">
      <div className="chat-container">
        {recievedMessage &&
          recievedMessage.map((item, index) => (
            <div>
              {item.is_answer === 0 && item.writer_id === myUUID && (
                <div className="chat-container__chat__usr">
                  <input
                    type="button"
                    value="X"
                    onClick={() => deleteChatHandler(item.chat_id)}
                  />
                  <div className="chat-container__chat">{item.text_body}</div>
                  <div className="chat-container__time">{item.write_time}</div>
                </div>
              )}
              {item.is_answer === 0 && item.writer_id !== myUUID && (
                <div className="chat-container__chat__other">
                  <div className="chat-container__chat">{item.text_body}</div>
                  <div className="chat-container__time">{item.write_time}</div>
                </div>
              )}
              {seeAnswerBox && item.is_answer === 1 && (
                <div>
                  <div className="chat-container__chat__question">
                    <div className="chat-container__question">
                      {item.text_body}
                    </div>
                  </div>
                  <div>
                    <input
                      type="text"
                      placeholder="답을 작성해주세요."
                      autofocus
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
      {!hideInputBox && (
        <Inputbox
          onSendMessage={(messageData) => sendMessageHandler(messageData)}
        />
      )}
      <div>
        ANSWER
        {answers.length > 0 &&
          answers.map((item, index) => (
            <div>
              <div>
                질문 {index + 1} : {item.question_contents}
              </div>
              <div>첫 번째 대답 : {item.first_answer}</div>
              <div>두 번째 대답 : {item.second_answer}</div>
              <div>대답한 날짜 : {item.answer_date}</div>
              <br />
            </div>
          ))}
      </div>
      <div>
        <Exceptword />
      </div>
      <div>많이 쓴 단어 상위 {wordNum}개</div>
      <input type="button" value="3개" onClick={threeWordsHandler} />
      <input type="button" value="5개" onClick={fiveWordsHandler} />
      <input type="button" value="10개" onClick={tenWordsHandler} />
      <div>
        내가 쓴 단어
        <br />
        {myWords.map((item, index) => (
          <div>
            {index + 1}위 : {item}
          </div>
        ))}
      </div>
      <div>
        상대방이 쓴 단어
        <br />
        {otherWords.map((item, index) => (
          <div>
            {index + 1}위 : {item}
          </div>
        ))}
      </div>
      <Changepw />
      <Logout />
      <Withdrawal />
      <Cutconn />
    </div>
  );
};

export default Chatpage;
