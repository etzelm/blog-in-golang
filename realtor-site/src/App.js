import React from 'react';
import Main from './components/Main';
import NavBar from './components/NavBar';

function App() {

  var state = {
    loggedIn: true,
    user: "etzelm@live.com"
  };

  return (
    <div className="App">
    
      <NavBar loggedIn={state["loggedIn"]} user={state["user"]} />
      <Main loggedIn={state["loggedIn"]} user={state["user"]} />

    </div>
  );
}

export default App;
