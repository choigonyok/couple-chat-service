import { useEffect, useState } from "react";
import "./Calender.css";
import axios from "axios";
import Dday from "./Dday";

const Calender = () => {
  const date1 = new Date(); // 06/16, 이달 1일 계산
  const date2 = new Date(); // 06/16, 이달 말일 계산
  const thisYear = date1.getFullYear();
  const thisDate = date1.getDate();
  const thisMonth = date1.getMonth(); // thisMonth = 7
  const [month, setMonth] = useState(thisMonth); // 기본값 month = 7
  const [year, setYear] = useState(thisYear); // 기본값 month = 7
  date1.setDate(1); // date1 = 06/01
  date1.setMonth(month); // date1 = 07/01
  date2.setDate(1); // date1 = 06/01
  date2.setMonth(month + 1); // date1 = 07/01
  const [dateArray, setDateArray] = useState([]);
  const [weeksArray, setWeeksArray] = useState([]);
  const [dateInfo, setDateInfo] = useState(0);
  const [inputAnniversary, setInputAnniversary] = useState("");
  const [render, setRender] = useState(false);
  const [anniversaries, setAnniversaries] = useState([]);
  const [inputDDay, setInputDDay] = useState(false);

  let firstWeeksLastDate = 7 - date1.getDay();
  let lastDateOfThisMonth = date2.getDate(date2.setDate(date2.getDate() - 1));

  let weeksOfThisMonth;
  for (let i = 0; firstWeeksLastDate + 7 * i < lastDateOfThisMonth; i++) {
    weeksOfThisMonth = i;
  }
  weeksOfThisMonth += 2;

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/anniversary/", {
        params: {
          month: month + 1 <= 0 ? ((month + 1) % 12) + 12 : (month % 12) + 1,
          year: year,
        },
      })
      .then((response) => {
        if (response.data.length == 0) {
          setAnniversaries(response.data);
        } else {
          setAnniversaries([...response.data]);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, [month, render]);

  useEffect(() => {
    const array = [];
    let temp_date = 1;
    for (let i = 0; i < 7 * weeksOfThisMonth; i++) {
      if (
        date1.getDay() <= i &&
        i <= lastDateOfThisMonth + date1.getDay() - 1
      ) {
        array[i] = temp_date;
        temp_date += 1;
      } else {
        array[i] = 0;
      }
    }
    setDateArray(array);

    const temp_weeks = [];
    for (let j = 0; j < weeksOfThisMonth; j++) {
      temp_weeks[j] = j;
    }
    setWeeksArray(temp_weeks);
  }, [month]);

  const prevMonthHandler = () => {
    if (month % 12 === 0) {
      setYear(year - 1);
    }
    setAnniversaries([]);
    setMonth(month - 1);
    setDateInfo(0);
    setInputAnniversary("");
  };

  const nextMonthHandler = () => {
    if ((month + 1) % 12 === 0) {
      setYear(year + 1);
    }
    setAnniversaries([]);
    setMonth(month + 1);
    setDateInfo(0);
    setInputAnniversary("");
  };

  const setTodayHanndler = () => {
    setAnniversaries([]);
    setMonth(thisMonth);
    setYear(thisYear);
    setDateInfo(0);
    setInputAnniversary("");
  };

  const dateClickHandler = (value) => {
    setDateInfo(value);
  };

  const deleteBoxHandler = () => {
    setDateInfo(0);
    setInputAnniversary("");
  };

  const sendAnniversaryHandler = () => {
    const monthInfo =
      month + 1 <= 0 ? ((month + 1) % 12) + 12 : (month % 12) + 1;

    const sendData = {
      year: year,
      month: monthInfo,
      date: dateInfo,
      contents: inputAnniversary,
      d_day: inputDDay,
    };

    axios
      .post(
        process.env.REACT_APP_HOST_URL + "/api/anniversary",
        JSON.stringify(sendData)
      )
      .then((response) => {
        setInputAnniversary("");
        setDateInfo(0);
        setRender(!render);
        setInputDDay(false);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const inputAnniversaryHandler = (e) => {
    setInputAnniversary(e.target.value);
  };

  const deleteAniversaryHandler = (value) => {
    axios
      .delete(process.env.REACT_APP_HOST_URL + "/api/anniversary/" + value)
      .then((response) => {
        setRender(!render);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const dDayHandler = () => {
    setInputDDay(!inputDDay);
  };

  return (
    <div>
      <div>
        <div>
          <Dday />
        </div>
        <div>
          {year}년 {month + 1 <= 0 ? ((month + 1) % 12) + 12 : (month % 12) + 1}
          월
        </div>
        {dateInfo !== 0 && (
          <div>
            <div>
              {year}년
              {month + 1 <= 0 ? ((month + 1) % 12) + 12 : (month % 12) + 1}월
              {dateInfo}일
            </div>
            {anniversaries.length !== 0 &&
              anniversaries.map(
                (item, index) =>
                  item.date === dateInfo && (
                    <div>
                      {item.contents}
                      <input
                        type="button"
                        value="X"
                        onClick={() =>
                          deleteAniversaryHandler(item.anniversary_id)
                        }
                      />
                    </div>
                  )
              )}
            <div>
              <input
                type="text"
                value={inputAnniversary}
                onChange={inputAnniversaryHandler}
                placeholder="일정을 입력해주세요"
              />
              <input
                type="button"
                value="일정 저장하기"
                onClick={sendAnniversaryHandler}
              />
              <input type="button" value="X" onClick={deleteBoxHandler} />
            </div>
            <div>
              <input
                type="checkbox"
                name="dday"
                id="dday"
                onChange={dDayHandler}
              />
              <label for="dday">D-DAY 설정</label>
            </div>
          </div>
        )}
        <div>
          <input type="button" value="prev" onClick={prevMonthHandler} />
          <input type="button" value="today" onClick={setTodayHanndler} />
          <input type="button" value="next" onClick={nextMonthHandler} />
        </div>
        <div className="calender">
          <div className="calender-day__container">
            <div className="calender-day">일</div>
            <div className="calender-day">월</div>
            <div className="calender-day">화</div>
            <div className="calender-day">수</div>
            <div className="calender-day">목</div>
            <div className="calender-day">금</div>
            <div className="calender-day">토</div>
          </div>
          {weeksArray.map((item, index) => (
            <div className="calender-container">
              <div
                className={
                  thisMonth === month && dateArray[index * 7 + 0] === thisDate
                    ? "date__today"
                    : "date"
                }
                onClick={() => dateClickHandler(dateArray[index * 7 + 0])}
              >
                <div>
                  <div className="date__date">
                    {dateArray[index * 7 + 0] !== 0 && dateArray[index * 7 + 0]}
                  </div>
                  {anniversaries.length !== 0 &&
                    anniversaries.map((item, i) => (
                      <div>
                        {item.date === dateArray[index * 7 + 0] &&
                          item.contents}
                      </div>
                    ))}
                </div>
              </div>
              <div
                className={
                  thisMonth === month && dateArray[index * 7 + 1] === thisDate
                    ? "date__today"
                    : "date"
                }
                onClick={() => dateClickHandler(dateArray[index * 7 + 1])}
              >
                <div className="date__date">
                  {dateArray[index * 7 + 1] !== 0 && dateArray[index * 7 + 1]}
                </div>
                <div>
                  {anniversaries.length !== 0 &&
                    anniversaries.map((item, i) => (
                      <div>
                        {item.date === dateArray[index * 7 + 1] &&
                          item.contents}
                      </div>
                    ))}
                </div>
              </div>
              <div
                className={
                  thisMonth === month && dateArray[index * 7 + 2] === thisDate
                    ? "date__today"
                    : "date"
                }
                onClick={() => dateClickHandler(dateArray[index * 7 + 2])}
              >
                <div className="date__date">
                  {dateArray[index * 7 + 2] !== 0 && dateArray[index * 7 + 2]}
                </div>
                <div>
                  {anniversaries.length !== 0 &&
                    anniversaries.map((item, i) => (
                      <div>
                        {item.date === dateArray[index * 7 + 2] &&
                          item.contents}
                      </div>
                    ))}
                </div>
              </div>
              <div
                className={
                  thisMonth === month && dateArray[index * 7 + 3] === thisDate
                    ? "date__today"
                    : "date"
                }
                onClick={() => dateClickHandler(dateArray[index * 7 + 3])}
              >
                <div className="date__date">
                  {dateArray[index * 7 + 3] !== 0 && dateArray[index * 7 + 3]}
                </div>
                <div>
                  {anniversaries.length !== 0 &&
                    anniversaries.map((item, i) => (
                      <div>
                        {item.date === dateArray[index * 7 + 3] &&
                          item.contents}
                      </div>
                    ))}
                </div>
              </div>
              <div
                className={
                  thisMonth === month && dateArray[index * 7 + 4] === thisDate
                    ? "date__today"
                    : "date"
                }
                onClick={() => dateClickHandler(dateArray[index * 7 + 4])}
              >
                <div className="date__date">
                  {dateArray[index * 7 + 4] !== 0 && dateArray[index * 7 + 4]}
                </div>
                <div>
                  {anniversaries.length !== 0 &&
                    anniversaries.map((item, i) => (
                      <div>
                        {item.date === dateArray[index * 7 + 4] &&
                          item.contents}
                      </div>
                    ))}
                </div>
              </div>
              <div
                className={
                  thisMonth === month && dateArray[index * 7 + 5] === thisDate
                    ? "date__today"
                    : "date"
                }
                onClick={() => dateClickHandler(dateArray[index * 7 + 5])}
              >
                <div className="date__date">
                  {dateArray[index * 7 + 5] !== 0 && dateArray[index * 7 + 5]}
                </div>
                <div>
                  {anniversaries.length !== 0 &&
                    anniversaries.map((item, i) => (
                      <div>
                        {item.date === dateArray[index * 7 + 5] &&
                          item.contents}
                      </div>
                    ))}
                </div>
              </div>
              <div
                className={
                  thisMonth === month && dateArray[index * 7 + 6] === thisDate
                    ? "date__today"
                    : "date"
                }
                onClick={() => dateClickHandler(dateArray[index * 7 + 6])}
              >
                <div className="date__date">
                  {dateArray[index * 7 + 6] !== 0 && dateArray[index * 7 + 6]}
                </div>
                <div>
                  {anniversaries.length !== 0 &&
                    anniversaries.map((item, i) => (
                      <div>
                        {item.date === dateArray[index * 7 + 6] &&
                          item.contents}
                      </div>
                    ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Calender;
