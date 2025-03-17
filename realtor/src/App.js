import React, { useEffect, useState } from 'react';
import NavBar from './components/NavBar';
import Main from './components/Main';

function App() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const initGoogleAuth = async () => {
      try {
        if (!window.gapi) {
          const response = await fetch('https://apis.google.com/js/platform.js');
          if (!response.ok) throw new Error('Failed to load gapi');
        }

        await window.gapi.load('auth2', () => {
          window.gapi.auth2.init({
            client_id: 'ThisIsSupposedToBeAnId',
            scope: 'email',
            prompt: 'select_account'
          }).then(() => {
            const auth2 = window.gapi.auth2.getAuthInstance();
            const signedIn = auth2.isSignedIn.get();

            if (signedIn) {
              const email = auth2.currentUser.get().getBasicProfile().getEmail();
              setUser(email);
              setLoggedIn(true);
            }
            setLoaded(true);
          }).catch(error => {
            console.error('Google Auth initialization failed:', error);
            setLoaded(true);
          });
        });
      } catch (error) {
        console.error('Error loading Google Auth:', error);
        setLoaded(true);
      }
    };

    initGoogleAuth();
  }, []);

  if (!loaded) {
    return <div>Loading...</div>;
  }

  return (
    <div className="App">
      <NavBar 
        loggedIn={loggedIn} 
        user={user} 
        onSignIn={() => {
          window.gapi.auth2.getAuthInstance().signIn();
        }}
        onSignOut={() => {
          window.gapi.auth2.getAuthInstance().signOut();
        }}
      />
      <Main 
        loggedIn={loggedIn} 
        user={user}
      />
    </div>
  );
}

export default App;