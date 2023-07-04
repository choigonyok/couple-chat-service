import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import Chatpage from "./Chatpage";
import Login from "./Login";
import Connpage from "./Connpage";

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <BrowserRouter>
          <Routes>
            <Route path="/chat" element={<Chatpage />}></Route>
            <Route path="/" element={<Login />}></Route>
            <Route path="/conn" element={<Connpage/>}></Route>
          </Routes>
        </BrowserRouter>
      </header>
    </div>
  );
}

export default App;
