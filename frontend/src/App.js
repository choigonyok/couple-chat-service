import './App.css';
import Chatpage from './Chatpage';
import Login from './Login';
import Logout from './Logout';
import Signup from './Signup';


function App() {
  return (
    <div className="App">
      <header className="App-header">
        <Chatpage/>
        <Signup/>
        <Login/>
        <Logout/>
      </header>
    </div>
  );
}

export default App;
