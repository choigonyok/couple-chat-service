import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const Connrecieved = () => {
  const navigator = useNavigate();

  const [requested, setRequested] = useState([]);
  const [requesting, setRequesting] = useState({});
  const [render, setRender] = useState(false);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/request/send")
      .then((response) => {
        setRequesting(response.data);
      })
      .catch((error) => {
        console.log(error);
      });

    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/request/recieved")
      .then((response) => {
        setRequested([...response.data]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const connectHandler = (value) => {
    const sendData_conn = {
      uuid_delete : value.Requester_uuid
    };

    axios
      .put(process.env.REACT_APP_HOST_URL+"/api/request",JSON.stringify(sendData_conn))
      .then((response)=>{
        navigator("/chat");
      })
      .catch((error)=>{
        console.log(error);
      })
  }

  const deleteRequestHandler = (value) => {
    
    axios
      .delete(process.env.REACT_APP_HOST_URL+"/api/request/"+value.Request_id)
      .then((response)=>{
        setRender(!render);
      })
      .catch((error)=>{
        console.log(error);
      })
  }

  return (
    <div>
      보낸 커넥션
      {requesting && (
        <div>
          <div>
            <p>요청 보낸 대상</p>
            {requesting.Target_id}
          </div>
          <div>
            <p>요청 보낸 시간</p>
            {requesting.Request_time}
          </div>
        </div>
      )}
      <p>받은 커넥션</p>
      {requested.length > 0 && 
          requested.map((item, index) =>  (
            <div>
              <div>
                <p>요청 받은 대상</p>
                {item.Requester_id}
              </div>
              <div>
                <p>요청 받은 시간</p>
                {item.Request_time}
              </div>
              <div>
                <input type="button" value="연결하기" onClick={() => connectHandler(item)}/>
                <input type="button" value="요청삭제하기" onClick={() => deleteRequestHandler(item)}/>
              </div>
            </div>
          ))
      }
    </div>
  );
};

export default Connrecieved;
