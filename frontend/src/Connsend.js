import axios from "axios";
import { useState } from "react";

const Connsend = () => {
  const [inputId, setInputId] = useState("");

  const sendConnHandler = () => {
    if (inputId !== "") {
      const data = {
        input_id: inputId,
      };

      axios
        .post(
          process.env.REACT_APP_HOST_URL + "/api/request",
          JSON.stringify(data)
        )
        .then((response) => {
          alert(
            "요청이 완료되었습니다. 상대방이 승인하면 채팅 서비스를 이용할 수 있습니다."
          );
        })
        .catch((error) => {
          if (error.response.data === "ALREADY_REQUEST") {
            alert("이미 진행중인 요청이 있습니다. 삭제 후 재시도해주세요.");
          } else if (error.response.data === "NOT_YOURSELF") {
            alert("자기 자신과는 연결할 수 없습니다.");
          } else if (error.response.data === "ALREADY_CONNECTED") {
            alert("해당 사용자는 이미 다른 사람과 연결되어 있습니다.");
          } else if (error.response.data === "NOT_EXIST") {
            alert("해당 ID를 가진 사용자가 존재하지 않습니다.");
          } else {
            console.log(error);
          }
        });
    } else {
      alert("요청할 ID가 입력되지 않았습니다.");
    }
  };

  const inputIdHandler = (e) => {
    setInputId(e.target.value);
  };

  return (
    <div>
      <div>보낸 커넥션</div>

      <div>
        <input
          type="text"
          placeholder="연결할 상대방의 ID를 입력해주세요."
          value={inputId}
          onChange={inputIdHandler}
        />
        <div>
          <input
            type="button"
            value="연결 요청 보내기"
            onClick={sendConnHandler}
          />
        </div>
      </div>
    </div>
  );
};

export default Connsend;
