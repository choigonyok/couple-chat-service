import { useState } from "react";
import "./Inputbox.css";
import axios from "axios";

const Inputbox = (props) => {
  const [chat, setChat] = useState("");
  let chatData = {
    chat_body: "",
  };

  const chatHandler = (e) => {
    setChat(e.target.value);
    chatData = {
      chat_body: e.target.value,
    };
  };

  const enterHandler = (e) => {
    if (e.key === "Enter") {
      if (chat !== "") {
        if (e.nativeEvent.isComposing) {
          return;
        } else {
          props.onSendMessage(e.target.value);
          setChat("");
        }
        // 한글 두 번 출력되는 문제
        // 출처 : https://velog.io/@euji42/solved-%ED%95%9C%EA%B8%80-%EC%9E%85%EB%A0%A5%EC%8B%9C-2%EB%B2%88-%EC%9E%85%EB%A0%A5%EC%9D%B4-%EB%90%98%EB%8A%94-%EA%B2%BD%EC%9A%B0
      }
    }
  };

  return (
    <div>
      <input
        type="text"
        placeholder="ENTER로 채팅을 입력할 수 있습니다"
        className="inputbox"
        value={chat}
        onChange={chatHandler}
        onKeyDownCapture={enterHandler}
      ></input>
    </div>
  );
};

export default Inputbox;
