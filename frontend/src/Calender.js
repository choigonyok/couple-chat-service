import { useEffect, useState } from "react";

const Calender = () => {
  const date1 = new Date(); // 06/16
  const thisMonth = date1.getMonth() + 1; // thisMonth = 7
  const [month, setMonth] = useState(thisMonth); // 기본값 month = 7
  date1.setDate(1); // date1 = 06/01
  date1.setMonth(month); // date1 = 07/01

  console.log("DATE1 : ", date1.getMonth());

  let firstWeeksLastDate =  (7 - date1.getDay()) % 7;
  let lastDateOfThisMonth = date1.getDate(date1.setDate(date1.getDate() - 1));

  console.log(lastDateOfThisMonth);

  let weeksOfThisMonth;
  for (let i = 0; firstWeeksLastDate + 7 * i < lastDateOfThisMonth; i++) {
    weeksOfThisMonth = i;
  }
  weeksOfThisMonth += 2;

  const prevMonthHandler = () => {
    setMonth(month - 1);
  };

  const nextMonthHandler = () => {
    setMonth(month + 1);
  };

  return (
    <div>
      <div>
        <div>{month}월</div>
        <div>
          <input type="button" value="prev" onClick={prevMonthHandler} />
          <input type="button" value="next" onClick={nextMonthHandler} />
        </div>
        <div>{weeksOfThisMonth}개 주차 존재</div>
        <div>{lastDateOfThisMonth}일까지 존재</div>
        
      </div>
    </div>
  );
};

export default Calender;
