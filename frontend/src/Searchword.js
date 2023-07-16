import axios from "axios";
import { useState } from "react";

const Searchword = () => {
  const [inputWord, setInputWord] = useState();

  const inputWordHandler = (e) => {
    setInputWord(e.target.value);
  };

  const sendHandler = () => {
    axios
        .get(process.env.REACT_APP_HOST_URL+"/api/chat")
        .then((response)=>{
            console.log(response);
        })
        .catch((error)=>{
            console.log(error);
        })
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
    </div>
  );
};

export default Searchword;
