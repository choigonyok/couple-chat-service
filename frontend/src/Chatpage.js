import REACT, { useEffect, useRef, useState } from "react";
// NULL을 사용하려면 REACT를 import 해줘야함
import "./Chatpage.css";
import Inputbox from "./Inputbox";

const Chatpage = () => {
  const [newSocket, setNewSocket] = useState(null);
  const [recievedMessage, setRecievedMessage] = useState([]);
  const [usrID, setUsrID] = useState([]);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socket.onopen = () => {
      console.log("CONNECTION CONNECTED");
      setNewSocket(socket);
    };

    socket.onmessage = (e) => {
      const parsedData = JSON.parse(e.data);
      if (parsedData.created_id) {
        setUsrID(parsedData.created_id);
        console.log("CREATED ID : ", parsedData.created_id);
      } else {
        setRecievedMessage((prev) => [...prev, String(parsedData)]);
      }
    };
    socket.onclose = () => {
      console.log("CONNECTION CLOSED");
    };

    return () => {
      socket.close();
    };
  }, []);

  const sendMessageHandler = (data) => {
    if (newSocket !== null) {
      let now = Date("2023-06-06 12:00:00");
      const sendData = {
        text_body: data,
        writer_id: usrID,
        write_time: now,
      };
      newSocket.send(sendData);
    } else alert("상대방에게 메세지를 보낼 수 없는 상태입니다.");
  };

  return (
    <div className="page-container">
      <div className="chat-container">
        {recievedMessage &&
          recievedMessage.map((item, index) => (
            <div className="chat-container__chat__usr">
              <p className="chat-container__usr">{item}</p>
            </div>
          ))}
      </div>
      <Inputbox
        onSendMessage={(messageData) => sendMessageHandler(messageData)}
      />
    </div>
  );
};

export default Chatpage;
