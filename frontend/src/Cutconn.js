import axios from "axios";
import "./Cutconn.css";

const Cutconn = () => {

  const cutConnHandler = () => {
    axios
      .delete(process.env.REACT_APP_HOST_URL + "/api/conn")
      .then((response) => {
        alert("3분 뒤에 상대방과의 커넥션이 끊어질 예정입니다.");
      })
      .catch((error) => {
        if (error.response.status === 400) {
          alert("이미 커넥션 끊기가 예정되어 있습니다.");
        } else {
          console.log(error);
        }
      });
  };

  return (
    <div className="buttons__conn">
      <input
        type="button"
        value="연결끊기"
        onClick={cutConnHandler}
        className="buttons"
      />
    </div>
  );
};

export default Cutconn;
