import axios from "axios";
import { useEffect, useState } from "react";

const Dday = () => {
const [dDay, setDDay] = useState([]);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/anniversary/dday")
      .then((response) => {
        if (response.status !== 204) {
            setDDay(response.data);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  return <div>
    {dDay.length === 1 && 
    (
        <div>
            {dDay[0].contents}
            {dDay[0].year}
            {dDay[0].month}
            {dDay[0].date}
        </div>
    )}
  </div>;
};

export default Dday;
