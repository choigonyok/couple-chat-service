import REACT, { useEffect, useRef, useState } from "react";
// NULL을 사용하려면 REACT를 import 해줘야함
import "./Chatpage.css";
import Inputbox from "./Inputbox";

const Chatpage = () => {
  const [newSocket, setNewSocket] = useState(null);
  const [recievedMessage, setRecievedMessage] = useState([]);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socket.onopen = () => {
      console.log("CONNECTION CONNECTED");
      setNewSocket(socket);
    };

    socket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      setRecievedMessage((prev) => [...prev, String(data)]);
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
      console.log(data);
      newSocket.send(data);
    } else alert("상대방애개 메세지를 보낼 수 없는 상태입니다.");
  };

  return (
    <div className="page-container" id="page-container">
      <div className="chat-container">
        {recievedMessage && recievedMessage.map((item, index) => (
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
