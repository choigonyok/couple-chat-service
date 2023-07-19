import axios from "axios";
import "./Canclecutconn.css";

const Canclecutconn = () => {

  const rollBackConnHandler = () => {
    axios
      .put(process.env.REACT_APP_HOST_URL + "/api/conn")
      .then((response) => {
        if (response.status === 204) {
          alert("커넥션 끊기가 신청되지 않은 상태입니다.");
        } else {
          alert("이전에 예정되어있던 커넥션 끊기가 정상적으로 취소되었습니다.");
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <div>
      <input
        type="button"
        value="연결끊기 취소"
        className="buttons"
        onClick={rollBackConnHandler}
      />
    </div>
  );
};

export default Canclecutconn;
