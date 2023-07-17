import axios from "axios";
import { useEffect, useState } from "react";

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
    <div>
      <div>
        <input
          type="text"
          placeholder="제외할 단어롤 입력해주세요"
          value={inputWord}
          onChange={inputWordHandler}
        />
        <input type="button" value="확인" onClick={clickHandler} />
      </div>
      <div>현재 순위에서 제외된 단어들</div>
      <div>
        {exceptWords.length > 1 &&
          exceptWords.map((item, index) => (
            <div>
              {item}
              <input
                type="button"
                value="X"
                onClick={() => deleteExceptWordHandler(item)}
              />
            </div>
          ))}
        {exceptWords.length === 1 && (
          <div>
            {exceptWords[0]}
            <input
              type="button"
              value="X"
              onClick={() => deleteExceptWordHandler(exceptWords[0])}
            />
          </div>
        )}
      </div>
      <div>많이 쓴 단어 상위 {wordNum}개</div>
      <input type="button" value="3개" onClick={threeWordsHandler} />
      <input type="button" value="5개" onClick={fiveWordsHandler} />
      <input type="button" value="10개" onClick={tenWordsHandler} />
      <div>
        내가 쓴 단어
        <br />
        {myWords.map((item, index) => (
          <div>
            {index + 1}위 : {item}
          </div>
        ))}
      </div>
      <div>
        상대방이 쓴 단어
        <br />
        {otherWords.map((item, index) => (
          <div>
            {index + 1}위 : {item}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Exceptword;
