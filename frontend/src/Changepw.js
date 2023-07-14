import axios from "axios";
import { useState } from "react";
import Inputbox from "./Inputbox";

const Changepw = () => {
  const [inputBox, setInputBox] = useState(false);
  const [inputPW, setInputPW] = useState("");

  const inputPWHandler = (e) => {
    setInputPW(e.target.value);
  };

  const sendPWHandler = () => {
    const sendData = {
      usr_pw: inputPW,
    };
    axios
      .put(process.env.REACT_APP_HOST_URL + "/api/usr", sendData)
      .then((response) => {
        console.log(response);
        setInputPW("");
        setInputBox(false);
      })
      .catch((error) => {
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
          value="비밀번호 변경하기"
          onClick={clickHandler}
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
