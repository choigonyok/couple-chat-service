import axios from "axios";
import { useState } from "react";

const Login = () => {

  const [inputID, setInputID] = useState("");
  const [inputPW, setInputPW] = useState("");

  const inputIDHandler = (e) => {
    setInputID(e.target.value);
  };

  const inputPWHandler = (e) => {
    setInputPW(e.target.value);
  };

  const loginHandler = () => {
        const usrData = {
                usr_id: inputID,
                usr_pw: inputPW,
              };
        axios
        .post(process.env.REACT_APP_HOST_URL+"/api/login", usrData)
        .then((response)=>{
                alert("로그인에 성공했습니다.");
                setInputPW("");
                setInputID("");
                //navigator("/");
        })
        .catch((error)=>{
                alert("ID 혹은 PASSWORD가 틀렸습니다.");
        })
  };

  return (
    <div>
      <div>
        <input
          type="text"
          placeholder="ID"
          value={inputID}
          onChange={inputIDHandler}
        />
        <div>
          <input
            type="password"
            placeholder="PASSWORD"
            value={inputPW}
            onChange={inputPWHandler}
          />
        </div>
      </div>
      <div>
        <input type="button" value="로그인하기" onClick={loginHandler} />
      </div>
    </div>
  );
};

export default Login;
