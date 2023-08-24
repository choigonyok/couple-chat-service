import { useNavigate } from "react-router-dom";
import Logout from "./Logout";
import axios from "axios";
import { useEffect, useState } from "react";
import Connrecieved from "./Connrecieved";
import Withdrawal from "./Withdrawal";
import Changepw from "./Changepw";
import "./Connpage.css";

const Connpage = () => {
  const [reRender, setReRender] = useState(false);
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
        if (error.response.status === 401) {
          navigator("/");
        } else {
          console.log(error);
        }
      });
  }, []);



  return (
    <div>
      <Connrecieved/>
      <div className="buttons-bottom">
        <Changepw/>
        <Logout />
        <Withdrawal />
      </div>
    </div>
  );
};

export default Connpage;
