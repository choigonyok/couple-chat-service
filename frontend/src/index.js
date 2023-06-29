import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <App />
  // 왜 자꾸 api가 두 번씩 호출되나 했더니 <React.StrictMode/> 얘 떄문이었음.
  // 개발모드에서 문제를 알아내기 위해 검사하면서 호출을 한 번 더 해버림
  // 우선 지우고 나중에 필요하면 다시 작성하기로 함
);