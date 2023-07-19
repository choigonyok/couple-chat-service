import axios from "axios";
import { useState } from "react";
import "./Changepw.css";

const Changepw = () => {
  const [inputBox, setInputBox] = useState(false);
  const [inputPW, setInputPW] = useState("");

  const inputPWHandler = (e) => {
    if (e.target.value.length > 20) {
      alert("PASSWORD는 최대 20자까지만 입력 가능합니다.");
    } else {
      setInputPW(e.target.value);
    }
  };

  const sendPWHandler = () => {
    if (inputPW === "") {
      alert("ID 혹은 PASSWORD가 작성되지 않았습니다.");
      return;
    }
    const sendData = {
      usr_pw: inputPW,
    };
    axios
      .put(process.env.REACT_APP_HOST_URL + "/api/usr", sendData)
      .then((response) => {
        alert("PASSWORD가 성공적으로 변경되었습니다.");
        setInputPW("");
        setInputBox(false);
      })
      .catch((error) => {
        if (error.response.status === 400) {
          alert("PASSWORD는 영어 소문자와 숫자 조합만 가능합니다.");
        }
        console.log(error);
      });
  };

  const clickHandler = () => {
    setInputBox(true);
  };

  return (
    <div>
      {!inputBox && (
        <input
          type="button"
          value="비밀번호 변경"
          onClick={clickHandler}
          className="buttons"
        />
      )}
      {inputBox && (
        <div>
          <input
            type="password"
            placeholder="PASSWORD"
            value={inputPW}
            onChange={inputPWHandler}
          />
          <input type="button" value="변경" onClick={sendPWHandler} />
        </div>
      )}
    </div>
  );
};

export default Changepw;
