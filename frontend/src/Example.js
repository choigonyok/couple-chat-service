import axios from "axios";
import { useState } from "react";

const Example=()=> {
        const [test, setTest] = useState();

        axios
                .get("http://localhost/api/test")
                // api호출은 go port num인 8080이 아니라 container port num인 1000으로 요청해야 통신이 됨
                // localhost:8080으로 요청하면 통신 안됨
                .then((response)=>{
                        console.log("SUCCESS");
                        setTest(response.data);
                })
                .catch((error)=>{
                        console.log("FAILED");
                })
        return <div>
                {test}
                <br/>
                {test}
        </div>
}

export default Example;