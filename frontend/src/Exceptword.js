import axios from "axios";
import { useEffect, useState } from "react";
import "./Exceptword.css";

const Exceptword = () => {
  const [inputWord, setInputWord] = useState("");
  const [exceptWords, setExceptWords] = useState([]);
  const [unLock, setUnLock] = useState(false);
  const [wordNum, setWordNum] = useState("3");
  const [otherWords, setOtherWords] = useState([]);
  const [myWords, setMyWords] = useState([]);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/rank/" + wordNum)
      .then((response) => {
        setMyWords([...response.data.mywords]);
        setOtherWords([...response.data.otherwords]);
      })
      .catch((error) => {
        if (error.response.status === 411) {
          alert("순위를 매기기엔 채팅의 수가 부족합니다.");
        } else {
          console.log(error);
        }
      });
  }, [wordNum]);

  useEffect(() => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/except")
      .then((response) => {
        if (response.status !== 204) {
          setExceptWords([...response.data]);
        } else {
          setExceptWords([]);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, [unLock]);

  const inputWordHandler = (e) => {
    setInputWord(e.target.value);
  };

  const clickHandler = () => {
    if (inputWord.length === 0) {
      alert("입력된 단어가 없습니다!");
    } else {
      const exceptWord = {
        except_word: inputWord,
      };
      axios
        .post(
          process.env.REACT_APP_HOST_URL + "/api/except",
          JSON.stringify(exceptWord)
        )
        .then((response) => {
          setInputWord("");
          setUnLock(!unLock);
        })
        .catch((error) => {
          if (error.response.status === 400) {
            alert("이미 제외된 단어입니다.");
          } else {
            console.log(error);
          }
        });
    }
  };

  const deleteExceptWordHandler = (item) => {
    const sendData = {
      except_word: item,
    };
    axios
      .delete(process.env.REACT_APP_HOST_URL + "/api/except/" + item)
      .then((response) => {
        setUnLock(!unLock);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const threeWordsHandler = () => {
    setWordNum("3");
  };
  const fiveWordsHandler = () => {
    setWordNum("5");
  };
  const tenWordsHandler = () => {
    setWordNum("10");
  };

  return (
    <div className="exceptword-container">
      <div className="exceptword-container__seperate">
        지난 한 주 우리가 많이 말한 단어 TOP {wordNum}
      </div>
      <div className="exceptword-container__seperate">
        <input type="button" value="3개" onClick={threeWordsHandler} className="exceptword-ranknum"/>
        <input type="button" value="5개" onClick={fiveWordsHandler} className="exceptword-ranknum"/>
        <input type="button" value="10개" onClick={tenWordsHandler} className="exceptword-ranknum"/>
      </div>
      <div className="exceptword-container__lists">
        <div className="exceptword-container__rank">
          <div className="exceptword-container__rank__other">
            <div>상대방이 쓴 단어</div>

            {otherWords.map((item, index) => (
              <div>
                {index + 1}위 : {item}
              </div>
            ))}
          </div>
        </div>
        <div className="exceptword-container__rank">
          <div className="exceptword-container__rank__mine">
            <div>내가 쓴 단어</div>
            {myWords.map((item, index) => (
              <div>
                {index + 1}위 : {item}
              </div>
            ))}
          </div>
        </div>
      </div>
      <div className="exceptword-container__seperate">
        <input
          type="text"
          placeholder="제외할 단어롤 입력해주세요"
          value={inputWord}
          onChange={inputWordHandler}
          className="exceptword-input"
        />
        <input
          type="button"
          value="확인"
          onClick={clickHandler}
          className="exceptword-button"
        />
      </div>
      <div className="exceptword-container__seperate">
        <div className="exceptword-font">현재 순위에서 제외된 단어들</div>
      </div>
      <div className="exceptword-container__seperate">
        {exceptWords.length > 1 &&
          exceptWords.map((item, index) => (
            <div className="exceptword-words">
              {item}
              <input
                type="button"
                value="X"
                className="exceptword-words__button"
                onClick={() => deleteExceptWordHandler(item)}
              />
            </div>
          ))}
        {exceptWords.length === 1 && (
          <div className="exceptword-words">
            {exceptWords[0]}
            <input
              type="button"
              value="X"
              className="exceptword-words__button"
              onClick={() => deleteExceptWordHandler(exceptWords[0])}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default Exceptword;
