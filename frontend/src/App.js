import logo from './logo.svg';
import './App.css';
import Example from './Example';
import Chatpage from './Chatpage';
import Signup from './Signup';


function App() {
  return (
    <div className="App">
      <header className="App-header">
        <Chatpage/>
        <Example/>
        <Signup/>
      </header>
    </div>
  );
}

export default App;
