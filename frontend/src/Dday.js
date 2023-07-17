import axios from "axios";
import { useEffect, useState } from "react";

const Dday = () => {
  const now = new Date();
  const dday = new Date();
  const [dDay, setDDay] = useState([]);
  const [subDate, setSubDate] = useState();

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/anniversary/dday")
      .then((response) => {
        if (response.status !== 204) {
          setDDay(response.data);
          dday.setFullYear(response.data[0].year);
          dday.setMonth(response.data[0].month - 1);
          dday.setDate(response.data[0].date);
          const subTime = dday.getTime() - now.getTime();
          console.log("now : ", now.toDateString());
          console.log("dday : ", dday.toDateString());
          setSubDate(subTime / 1000 / 3600 / 24);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  return (
    <div>
      {dDay.length === 1 && (
        <div>
          <div>
            {subDate > 0 && <div>D - {subDate}</div>}
            {subDate === 0 && <div>D - DAY</div>}
            {subDate < 0 && <div>D + {-subDate}</div>}
          </div>
          <div>
            {dDay[0].year}년 {dDay[0].month}월 {dDay[0].date}일
          </div>
          <div>{dDay[0].contents}</div>
        </div>
      )}
    </div>
  );
};

export default Dday;
