
import axios from "axios";
import { useEffect, useState } from "react";

const Example = () => {
  const [test, setTest] = useState();
  const [data, setData] = useState([]);

  useEffect(() => {
    axios
      .get("http://localhost/api/test")
      // api호출은 go port num인 8080이 아니라 container port num인 1000으로 요청해야 통신이 됨
      // localhost:8080으로 요청하면 통신 안됨
      .then((response) => {
        console.log("SUCCESS1");
        setTest(response.data);
      })
      .catch((error) => {
        console.log("FAILED");
      });

    axios
      .get("http://localhost/api/usr")
      .then((response) => {
        setData(response.data);
      })
      .catch((error) => {
        console.log("ERROR");
      });
  }, []);

  return (
    <div>
      <div>
        {data.map((item, index) => (
          <div>{item.id}</div>
        ))}
      </div>
      <div>{test}</div>
    </div>
  );
};

export default Example;
