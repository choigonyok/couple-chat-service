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
      {answers.length > 1 &&
        answers.map((item, index) => <div>{item.question_contents}</div>)}
    </div>
  );
};

export default Answer;
