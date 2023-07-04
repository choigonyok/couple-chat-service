import { useNavigate } from "react-router-dom";
import Logout from "./Logout";
import axios from "axios";
import { useEffect, useState } from "react";
import Connsend from "./Connsend";
import Connrecieved from "./Connrecieved";

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
        if (error.response.status === 400) {
          navigator("/");
        } else {
          console.log(error);
        }
      });
  }, []);

  const renderHandler = () => {
    setReRender(!reRender);
  };

  return (
    <div>
      <Connrecieved/>
      <Connsend/>
      <div>
        <Logout />
      </div>
    </div>
  );
};

export default Connpage;
