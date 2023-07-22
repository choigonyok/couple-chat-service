import { useState } from "react";
import "./Inputbox.css";
import axios from "axios";

const Inputbox = (props) => {
  const [chat, setChat] = useState("");
  const [img, setIMG] = useState();

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

  const fileHandler = (e) => {
    const file = new FormData();
    file.set("file", e.target.files[0]);
    axios
      .post(process.env.REACT_APP_HOST_URL + "/api/file", file, {
        withCredentials: true,
        "Content-Type": "multipart/form-data",
      })
      .then((response) => {
        alert("파일전송 성공");
        if (e.target.files[0].type.includes("image/")) {
          props.onSendMessage(0);
        } else {
          props.onSendMessage(1);
        }
        
        setChat("");
      })
      .catch((error) => {
        if (error.response.status === 400) {
          alert("전송할 수 없는 파일 형식입니다.");
        } else {
          alert("파일전송 실패");
          console.log(error);
        }
      });
  };

  return (
    <div className="inputbox-container">
      <input
        type="text"
        placeholder="ENTER로 채팅을 입력할 수 있습니다"
        className="inputbox"
        value={chat}
        onChange={chatHandler}
        onKeyDownCapture={enterHandler}
      ></input>
      <input
        type="file"
        id="imgfile"
        name="imgfile"
        className="file-button"
        onChange={fileHandler}
      />
    </div>
  );
};

export default Inputbox;
