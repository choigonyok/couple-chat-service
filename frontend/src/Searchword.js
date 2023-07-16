import axios from "axios";
import { useState } from "react";

const Searchword = () => {
  const [inputWord, setInputWord] = useState("");
  const [searchedWords, setSearchedWords] = useState([]);
  const [wordsIndex, setWordsIndex] = useState(0);

  const inputWordHandler = (e) => {
    setInputWord(e.target.value);
  };

  const sendHandler = () => {
    if (inputWord.length === 0) {
      alert("검색할 단어를 입력해주세요.");
      return;
    }
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/chat/word/" + inputWord)
      .then((response) => {
        setWordsIndex(0);
        if (response.data === 1) {
          setSearchedWords(response.data);
        } else {
          setSearchedWords([...response.data]);
        }
      })
      .catch((error) => {
        if (error.response.status === 404) {
          alert("검색 결과가 없습니다.");
        } else {
          console.log(error);
        }
      });
  };

  const prevClickHandler = () => {
    if (wordsIndex === 0) {
      alert("첫번째 검색 결과입니다.");
      return;
    }
    setWordsIndex(wordsIndex - 1);
  };

  const nextClickHandler = () => {
    if (wordsIndex === searchedWords.length - 1) {
      alert("마지막 검색 결과입니다.");
      return;
    }
    setWordsIndex(wordsIndex + 1);
  };

  return (
    <div>
      <input
        type="text"
        placeholder="검색할 단어를 입력하세요"
        value={inputWord}
        onChange={inputWordHandler}
      />
      <input type="button" value="검색" onClick={sendHandler} />
      <div>
        {searchedWords.length === 1 && (
          <div>
            <div>1/1</div>
            <div>{searchedWords[0].text_body}</div>
            <div>{searchedWords[0].write_time}</div>
          </div>
        )}
      </div>
      <div>
        {searchedWords.length > 1 && (
          <div>
            <div>
              {wordsIndex + 1}/{searchedWords.length}
            </div>
            <div>
              {searchedWords[wordsIndex].text_body}
              {searchedWords[wordsIndex].write_time}
            </div>
            <div>
              <input type="button" value="prev" onClick={prevClickHandler} />
              <input type="button" value="next" onClick={nextClickHandler} />
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Searchword;
