import './App.css';
import Chatpage from './Chatpage';
import Login from './Login';
import Signup from './Signup';


function App() {
  return (
    <div className="App">
      <header className="App-header">
        <Chatpage/>
        <Signup/>
        <Login/>
      </header>
    </div>
  );
}

export default App;
