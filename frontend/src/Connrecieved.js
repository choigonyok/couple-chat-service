import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./Connrecieved.css";

const Connrecieved = (props) => {
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
        console.log("RR,", response.data);
        setRequested([...response.data]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [render]);

  const connectHandler = (value) => {
    const sendData_conn = {
      uuid_delete: value.Requester_uuid,
    };

    axios
      .put(
        process.env.REACT_APP_HOST_URL + "/api/request",
        JSON.stringify(sendData_conn)
      )
      .then((response) => {
        navigator("/chat");
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const deleteRequestHandler = (value) => {
    axios
      .delete(
        process.env.REACT_APP_HOST_URL + "/api/request/" + value.Request_id
      )
      .then((response) => {
        alert("요쳥이 삭제되었습니다.");
        setRender(!render);
      })
      .catch((error) => {
        console.log(error);
      });
  };

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
          alert("요청이 완료되었습니다. 상대방이 승인하면 채팅 서비스를 이용할 수 있습니다.");
          setRender(!render);
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
      <div className="row-align">
        <div className="col-align">
          <div className="conn-title">보낸 커넥션</div>
          <div className="row-align">
            <div className="subtitle">요청 보낸 대상</div>
            <div className="subtitle">요청 보낸 시간</div>
          </div>
          {requesting && (
            <div>
              <div className="row-align">
                <div className="section">
                  <div>{requesting.Target_id}</div>
                </div>
                <div className="section">
                  <div>{requesting.Request_time}</div>
                </div>
              </div>
            </div>
          )}
        </div>
        <div className="col-align">
          <div className="conn-title">받은 커넥션</div>
          <div className="row-align">
            <div className="subtitle">요청 받은 대상</div>
            <div className="subtitle">요청 받은 시간</div>
          </div>
          {requested.length > 0 &&
            requested.map((item, index) => (
              <div>
                <div className="row-align">
                  <div className="section">
                    <div>{item.Requester_id}</div>
                  </div>
                  <div className="section">
                    <div>{item.Request_time}</div>
                  </div>
                </div>

                <div>
                  <div>
                    <input
                      type="button"
                      value="연결하기"
                      className="connpage_button"
                      onClick={() => connectHandler(item)}
                    />
                    <input
                      type="button"
                      value="요청삭제하기"
                      className="connpage_button"
                      onClick={() => deleteRequestHandler(item)}
                    />
                  </div>
                </div>
              </div>
            ))}
        </div>
      </div>
      <div className="row-align">
        <div className="col-align">
          <input
            type="text"
            placeholder="연결할 상대방의 ID를 입력해주세요."
            value={inputId}
            className="inputtext"
            onChange={inputIdHandler}
          />
            <input
              type="button"
              value="연결 요청 보내기"
              className="connpage_button"
              onClick={sendConnHandler}
            />
        </div>
      </div>
    </div>
  );
};

export default Connrecieved;
