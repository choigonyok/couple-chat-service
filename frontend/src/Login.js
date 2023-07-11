import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Signup from "./Signup";

const Login = () => {
  const navigator = useNavigate();

  const [isLogined, setIsLogined] = useState(false);
  const [inputID, setInputID] = useState("");
  const [inputPW, setInputPW] = useState("");

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/log")
      .then((response) => {
        if (response.data === "CONNECTED") {
          navigator("/chat");
        } else if (response.data === "NOT_CONNECTED") {
          navigator("/conn");
        }
      })
      .catch((error) => {
          console.log(error);
      });
  }, [isLogined]);

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
      .post(process.env.REACT_APP_HOST_URL + "/api/log", usrData)
      .then((response) => {
        alert("로그인에 성공했습니다.");
        setInputPW("");
        setInputID("");
        setIsLogined(!isLogined);
      })
      .catch((error) => {
        alert("ID 혹은 PASSWORD가 틀렸습니다.");
      });
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
      <br />
      <Signup />
    </div>
  );
};

export default Login;
