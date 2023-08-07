import axios from "axios";
import { useNavigate } from "react-router-dom";
import "./Logout.css";

const Logout = () => {
  const navigator = useNavigate();

  const logoutHandler = () => {

        axios
        .delete(process.env.REACT_APP_HOST_URL+"/api/log")
        .then((response)=>{
                alert("로그아웃 되었습니다.");
                navigator("/");
        })
        .catch((error)=>{
                alert("로그아웃에 실패했습니다.");
        })
  };

  return (
    <div>
      <div>
        <input type="button" value="로그아웃" onClick={logoutHandler} className="buttons"/>
      </div>
    </div>
  );
};

export default Logout;
