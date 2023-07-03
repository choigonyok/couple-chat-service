import axios from "axios";

const Logout = () => {

  const logoutHandler = () => {
        axios
        .delete(process.env.REACT_APP_HOST_URL+"/api/log")
        .then((response)=>{
                alert(response.data);
                //navigator("/");
        })
        .catch((error)=>{
                alert("로그아웃에 실패했습니다.");
        })
  };

  return (
    <div>
      <div>
        <input type="button" value="LOG OUT" onClick={logoutHandler} />
      </div>
    </div>
  );
};

export default Logout;
