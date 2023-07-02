import axios from "axios";
import { useState } from "react";


const Signup = () => {
        const [inputID, setInputID] = useState("");
        const [inputPW, setInputPW] = useState("");
        const [isIDChecked, setIsIDChecked] = useState(false);


        const inputIDHandler = (e) => {
                setInputID(e.target.value);
                setIsIDChecked(false);
                // 아이디 중복 확인을 위해서 작성값이 바뀌면 다시 중복 검사를 하도록 함
        }

        const inputPWHandler = (e) => {
                setInputPW(e.target.value);
        }

        const checkIDHandler = () => {
                const data = {
                        input_id : inputID,
                }
                axios
                        .post(process.env.REACT_APP_HOST_URL+"/api/id", JSON.stringify(data))
                        .then((response)=>{
                                console.log("중복된 아이디 없음. 사용가능");
                                setIsIDChecked(true);
                        })
                        .catch((error)=>{
                                console.log("아이디 중복 확인 에러 발생");
                        })
        }

        const signUpHandler = () => {
                if (inputID === "" || inputPW === "") {
                        alert("ID 혹은 PASSWORD가 작성되지 않았습니다.");
                } else if (isIDChecked === false) {
                        alert("ID 중복 확인을 해주세요.");
                } else {
                        const usrData = {
                                usr_id : inputID,
                                usr_pw : inputPW,
                        }
                        axios
                        .post(process.env.REACT_APP_HOST_URL+"/api/usr", usrData)
                        .then((response)=> {
                                alert("회원가입이 완료되었습니다.");
                        })
                        .catch((error)=>{
                                alert("회원가입에 실패했습니다.")
                        })
                };
        }

        return <div>
                <input type="text" placeholder="ID" value={inputID} onChange={inputIDHandler}/>
                <input type="button" value="중복 확인" onClick={checkIDHandler}/>
                <input type="password" placeholder="PASSWORD" value={inputPW} onChange={inputPWHandler}/>
                <input type="button" value="가입하기" onClick={signUpHandler}/>
        </div>
};

export default Signup;
