import axios from "axios";
import { useState } from "react";

const Searchword = () => {
  const now = new Date();

  const [inputWord, setInputWord] = useState("");
  const [searchedWords, setSearchedWords] = useState([]);
  const [searchedDate, setSearchedDate] = useState([]);
  const [wordsIndex, setWordsIndex] = useState(0);
  const [year, setYear] = useState(now.getFullYear().toString());
  const [month, setMonth] = useState((now.getMonth() + 1).toString());
  const [date, setDate] = useState(now.getDate().toString());

  const inputWordHandler = (e) => {
    setInputWord(e.target.value);
  };

  const sendWordHandler = () => {
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

  const sendDateHandler = () => {
    axios
      .get(process.env.REACT_APP_HOST_URL + "/api/chat/date", {
        params: {
          year: year,
          month: month,
          date: date,
        },
      })
      .then((response) => {
        console.log(response.data);
        if (response.status === 204) {
          alert("해당 날짜의 채팅 기록이 없습니다.");
        } else {
          setSearchedDate(response.data);
        }
      })
      .catch((error) => {
          console.log(error);
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

  const yearHandler = (e) => {
    setYear(e.target.value);
  };
  const monthHandler = (e) => {
    setMonth(e.target.value);
  };
  const dateHandler = (e) => {
    setDate(e.target.value);
  };

  return (
    <div>
      <input
        type="text"
        placeholder="검색할 단어를 입력하세요"
        value={inputWord}
        onChange={inputWordHandler}
      />
      <input type="button" value="검색" onClick={sendWordHandler} />
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
      <div>
        <select name="select" id="year" onChange={yearHandler}>
          <option value="2023">2023</option>
          <option value="2024">2024</option>
          <option value="2025">2025</option>
        </select>
        <select name="select" id="month" onChange={monthHandler}>
          <option value="1">1</option>
          <option value="2">2</option>
          <option value="3">3</option>
          <option value="4">4</option>
          <option value="5">5</option>
          <option value="6">6</option>
          <option value="7">7</option>
          <option value="8">8</option>
          <option value="9">9</option>
          <option value="10">10</option>
          <option value="11">11</option>
          <option value="12">12</option>
        </select>
        <select name="select" id="date" onChange={dateHandler}>
          <option value="1">1</option>
          <option value="2">2</option>
          <option value="3">3</option>
          <option value="4">4</option>
          <option value="5">5</option>
          <option value="6">6</option>
          <option value="7">7</option>
          <option value="8">8</option>
          <option value="9">9</option>
          <option value="10">10</option>
          <option value="11">11</option>
          <option value="12">12</option>
          <option value="13">13</option>
          <option value="14">14</option>
          <option value="15">15</option>
          <option value="16">16</option>
          <option value="17">17</option>
          <option value="18">18</option>
          <option value="19">19</option>
          <option value="20">20</option>
          <option value="21">21</option>
          <option value="22">22</option>
          <option value="23">23</option>
          <option value="24">24</option>
          <option value="25">25</option>
          <option value="26">26</option>
          <option value="27">27</option>
          <option value="28">28</option>
          <option value="29">29</option>
          <option value="30">30</option>
          <option value="31">31</option>
        </select>
      </div>
      <div>
        <input type="button" value="검색" onClick={sendDateHandler} />
      </div>
      <div>
        {searchedDate.length && <div>
          {searchedDate[0].text_body}
          {searchedDate[0].write_time}
          </div>}
      </div>
    </div>
  );
};

export default Searchword;
