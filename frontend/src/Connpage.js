import { useNavigate } from "react-router-dom";
import Logout from "./Logout";
import axios from "axios";
import { useEffect } from "react";
import Connsend from "./Connsend";

const Connpage = () => {
  const navigator = useNavigate();

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
        if (error.response.status === 400) {
          navigator("/");
        } else {
          console.log(error);
        }
      });
  }, []);

  return (
    <div>
      <div>받은 커넥션</div>
      <Connsend/>
      <div>
        <Logout />
      </div>
    </div>
  );
};

export default Connpage;
