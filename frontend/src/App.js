import logo from './logo.svg';
import './App.css';
import Example from './Example';
import Chatpage from './Chatpage';


function App() {
  return (
    <div className="App">
      <header className="App-header">
        <Chatpage/>
        <Example/>
      </header>
    </div>
  );
}

export default App;
