import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import Main from './components/Main';
import NavBar from './components/NavBar';

const App = () => {
  const [loggedIn, setLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  const [loaded, setLoaded] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const signOut = () => {
    if (window.gapi) {
      window.gapi.auth2.getAuthInstance().signOut().then(() => {
        window.gapi.auth2.getAuthInstance().disconnect();
        localStorage.removeItem('aToken');
        setLoggedIn(false);
        setUser(null);

        if (location.pathname === "/realtor/my-listings" || location.pathname === "/realtor/new") {
          navigate("/realtor");
        }
      });
    }
  };

  useEffect(() => {
    const token = localStorage.getItem('aToken');
    if (token) {
      setLoggedIn(true);
      setUser(token);
    }

    const initializeGapi = async () => {
      if (window.gapi) {
        await window.gapi.load('auth2', () => {
          window.gapi.auth2.init({
            client_id: 'ThisIsSupposedToBeAnId',
          }).then(() => {
            const auth2 = window.gapi.auth2.getAuthInstance();
            const signedIn = auth2.isSignedIn.get();

            if (signedIn) {
              const email = auth2.currentUser.get().getBasicProfile().getEmail();
              setUser(email);
              setLoggedIn(true);
            }
            setLoaded(true);
          });
        });

        window.gapi.load('signin2', () => {
          window.gapi.signin2.render('loginButton', {
            width: 100,
            height: 25,
            client_id: 'ThisIsSupposedToBeAnId',
            onsuccess: () => {
              const auth2 = window.gapi.auth2.getAuthInstance();
              const email = auth2.currentUser.get().getBasicProfile().getEmail();
              localStorage.setItem('aToken', email);
              setLoggedIn(true);
              setUser(email);

              if (location.pathname === "/realtor/my-listings" || location.pathname === "/realtor/new") {
                navigate(0); // Force reload for state consistency
              }
            }
          });
        });
      }
    };

    initializeGapi();
  }, [navigate, location]);

  if (!loaded) return null; // Add loading spinner if needed

  return (
    <div className="App">
      <NavBar 
        loggedIn={loggedIn} 
        user={user} 
        signOut={signOut}
      />
      <Main loggedIn={loggedIn} user={user} />
    </div>
  );
};

export default App;