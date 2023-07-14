import axios from "axios";
import { useNavigate } from "react-router-dom";

const Cutconn = () => {
  const navigator = useNavigate();

  const cutConnHandler = () => {
    axios
      .delete(process.env.REACT_APP_HOST_URL + "/api/conn")
      .then((response) => {
        alert("3분 뒤에 상대방과의 커넥션이 끊어질 예정입니다.");
        navigator("/conn");
      })
      .catch((error) => {
        if (error.response.status === 400) {
          alert("이미 커넥션 끊기가 예정되어 있습니다.");
        } else {
          console.log(error);
        }
      });
  };

  const rollBackConnHandler = () => {
    axios
      .put(process.env.REACT_APP_HOST_URL + "/api/conn")
      .then((response) => {
        alert("이전에 예정되어있던 커넥션 끊기가 정상적으로 취소되었습니다.")
      })
      .catch((error) => {
        if (error.response.status === 400) {
          alert("커넥션 끊기가 신청되지 않은 상태입니다.")
        } else {
          console.log(error);
        }
      });
  };

  return (
    <div>
      <input type="button" value="커넥션 끊기" onClick={cutConnHandler} />
      <input type="button" value="커넥션 끊기 취소" onClick={rollBackConnHandler} />
    </div>
  );
};

export default Cutconn;
