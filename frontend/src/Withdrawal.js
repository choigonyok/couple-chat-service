import axios from "axios";
import { useNavigate } from "react-router-dom";

const Withdrawal = () => {
    const navigator = useNavigate();

    const withdrawalHandler = () => {
        axios
            .delete(process.env.REACT_APP_HOST_URL+"/api/usr")
            .then((response)=>{
                alert("성공적으로 회원탈퇴가 완료되었습니다.");
            })
            .catch((error)=>{
                if (error.response.status === 400) {
                    alert("상대방과의 연결을 끊은 후에 회원탈퇴가 가능합니다.");
                } else {
                    console.log(error);
                }
            })
    }

    return <div>
        <input type="button" value="회원탈퇴하기" onClick={withdrawalHandler}/>
    </div>
}

export default Withdrawal;