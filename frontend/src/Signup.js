import axios from "axios";
import { useState } from "react";

const Signup = () => {
  const [inputID, setInputID] = useState("");
  const [inputPW, setInputPW] = useState("");
  const [isIDChecked, setIsIDChecked] = useState(false);
  const [isAlready, setIsAlready] = useState(false);

  const inputIDHandler = (e) => {
    setInputID(e.target.value);
    setIsIDChecked(false);
    setIsAlready(false);
    // 아이디 중복 확인을 위해서 작성값이 바뀌면 다시 중복 검사를 하도록 함
  };

  const inputPWHandler = (e) => {
    setInputPW(e.target.value);
  };

  const checkIDHandler = () => {
    const data = {
      input_id: inputID,
    };
    axios
      .post(process.env.REACT_APP_HOST_URL + "/api/id", JSON.stringify(data))
      .then((response) => {
        setIsIDChecked(true);
      })
      .catch((error) => {
        setIsAlready(true);
      });
  };

  const signUpHandler = () => {
    if (inputID === "" || inputPW === "") {
      alert("ID 혹은 PASSWORD가 작성되지 않았습니다.");
    } else if (isIDChecked === false && isAlready === false) {
      alert("ID 중복 확인을 해주세요.");
    } else if (isIDChecked === false && isAlready) {
      alert("이미 사용중인 아이디입니다.");
    } else {
      const usrData = {
        usr_id: inputID,
        usr_pw: inputPW,
      };
      axios
        .post(process.env.REACT_APP_HOST_URL + "/api/usr", usrData)
        .then((response) => {
          alert("회원가입이 완료되었습니다.");
          setInputID("");
          setInputPW("");
          setIsIDChecked(false);
        })
        .catch((error) => {
          alert(error.response.data);
          setInputID("");
          setInputPW("");
          setIsIDChecked(false);
        });
    }
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
        <input type="button" value="중복 확인" onClick={checkIDHandler} />
        <div>
                <h6>(ID는 알파벳 소문자로 시작하는, 알파벳 소문자와 숫자의 최대 20자 조합)</h6>
        </div>
      </div>
      {isIDChecked && (
        <div>
          <h5>현재 사용중이지 않은 아이디 입니다.</h5>
        </div>
      )}
      {isAlready && (
        <div>
          <h5>사용중인 아이디 입니다.</h5>
        </div>
      )}
      <div>
        <input
          type="password"
          placeholder="PASSWORD"
          value={inputPW}
          onChange={inputPWHandler}
        />
      </div>
      <div>
        <input type="button" value="가입하기" onClick={signUpHandler} />
      </div>
    </div>
  );
};

export default Signup;
