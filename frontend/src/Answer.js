import axios from "axios";
import { useEffect, useState } from "react";

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
      ANSWER
      {answers.length > 0 &&
        answers.map((item, index) => (
          <div>
            <div>
              질문 {index + 1} : {item.question_contents}
            </div>
            <div>첫 번째 대답 : {item.first_answer}</div>
            <div>두 번째 대답 : {item.second_answer}</div>
            <div>대답한 날짜 : {item.answer_date}</div>
            <br />
          </div>
        ))}
    </div>
  );
};

export default Answer;
