import axios from "axios";
import { useEffect, useState } from "react";
import "./Answer.css";

const Answer = () => {
  const [answers, setAnswers] = useState([]);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/answer")
      .then((response) => {
        setAnswers([...response.data]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  return (
    <div>
      {answers.length > 0 &&
        answers.map((item, index) => (
          <div className="answer-container">
            <div className="answer-container__question">
              {item.question_contents}
            </div>

            <div className="answer-container__seperate">
              <div className="answer__second">
                {item.order === 1 ? item.second_answer : item.first_answer}
              </div>
              <div className="answer__first">
                {item.order === 1 ? item.first_answer : item.second_answer}
              </div>
            </div>
            <div className="answer__date">{item.answer_date}</div>
            <br />
          </div>
        ))}
    </div>
  );
};

export default Answer;
