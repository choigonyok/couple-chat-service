import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import Chatpage from "./Chatpage";
import Login from "./Login";
import Logout from "./Logout";
import Signup from "./Signup";

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <BrowserRouter>
          <Routes>
            <Route path="/chat" element={<Chatpage />}></Route>
            <Route path="/logout" element={<Logout />}></Route>
            <Route path="/" element={<Login />}></Route>
            <Route path="/signup" element={<Signup />}></Route>
          </Routes>
        </BrowserRouter>
      </header>
    </div>
  );
}

export default App;
