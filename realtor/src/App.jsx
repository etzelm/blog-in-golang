import React, { useEffect, useState } from 'react';
import NavBar from './components/NavBar';
import Main from './components/Main';
import { NotificationManager } from 'react-notifications';
import 'react-notifications/lib/notifications.css';

const log = (message, data = {}) => {
  console.log(JSON.stringify({
    message,
    timestamp: new Date().toISOString(),
    ...data,
  }, null, 2));
};

function App() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  const [loaded, setLoaded] = useState(false);
  const [authLoading, setAuthLoading] = useState(false);

  // Function to update auth state
  const updateAuthState = (isSignedIn) => {
    const auth2 = window.gapi.auth2.getAuthInstance();
    if (isSignedIn) {
      const email = auth2.currentUser.get().getBasicProfile().getEmail();
      setUser(email);
      setLoggedIn(true);
      log('User signed in', { email });
    } else {
      setUser(null);
      setLoggedIn(false);
      log('User signed out');
    }
  };

  // Auth initialization (runs once)
  useEffect(() => {
    const initGoogleAuth = async () => {
      try {
        if (!window.gapi) {
          const response = await fetch('https://apis.google.com/js/platform.js');
          if (!response.ok) throw new Error('Failed to load Google Auth library');
        }

        await window.gapi.load('auth2', () => {
          // Validate client_id
          const clientId = 'ThisIsSupposedToBeAnId';
          if (clientId === 'ThisIsSupposedToBeAnId') {
            console.error('Invalid Google OAuth Client ID');
            NotificationManager.error('Authentication configuration error. Please contact support.', 'Error', 5000);
            setLoaded(true);
            return;
          }

          window.gapi.auth2.init({
            client_id: clientId,
            scope: 'email',
            prompt: 'select_account'
          }).then(() => {
            const auth2 = window.gapi.auth2.getAuthInstance();
            const signedIn = auth2.isSignedIn.get();
            updateAuthState(signedIn);

            // Listen for sign-in state changes
            const listener = auth2.isSignedIn.listen(updateAuthState);
            setLoaded(true);

            // Clean up listener on unmount
            return () => {
              if (listener) {
                log('Cleaning up auth listener');
              }
            };
          }).catch(error => {
            console.error('Google Auth initialization failed:', error);
            NotificationManager.error('Failed to initialize authentication. Please try again later.', 'Error', 5000);
            setLoaded(true);
          });
        });
      } catch (error) {
        console.error('Error loading Google Auth:', error);
        NotificationManager.error('Failed to load authentication library. Please check your network and try again.', 'Error', 5000);
        setLoaded(true);
      }
    };

    initGoogleAuth();
  }, []);

  // Handle sign-in with debounce and error handling
  const handleSignIn = async () => {
    if (authLoading) return; // Prevent multiple clicks
    setAuthLoading(true);
    try {
      const auth2 = window.gapi.auth2.getAuthInstance();
      await auth2.signIn();
      updateAuthState(true);
      NotificationManager.success('Successfully signed in!', 'Welcome', 3000);
    } catch (error) {
      console.error('Sign-in failed:', error);
      let errorMessage = 'Failed to sign in. Please try again.';
      if (error.error === 'popup_closed_by_user') {
        errorMessage = 'Sign-in canceled. Please complete the sign-in process.';
      } else if (error.error === 'access_denied') {
        errorMessage = 'Permission denied. Please grant the required permissions to sign in.';
      }
      NotificationManager.error(errorMessage, 'Sign-In Error', 5000);
    } finally {
      setAuthLoading(false);
    }
  };

  // Handle sign-out with debounce and error handling
  const handleSignOut = async () => {
    if (authLoading) return; // Prevent multiple clicks
    setAuthLoading(true);
    try {
      const auth2 = window.gapi.auth2.getAuthInstance();
      await auth2.signOut();
      updateAuthState(false);
      NotificationManager.success('Successfully signed out.', 'Goodbye', 3000);
    } catch (error) {
      console.error('Sign-out failed:', error);
      NotificationManager.error('Failed to sign out. Please try again.', 'Sign-Out Error', 5000);
    } finally {
      setAuthLoading(false);
    }
  };

  // Log state changes
  useEffect(() => {
    log('App state changed', { loggedIn, user, loaded, authLoading });
  }, [loggedIn, user, loaded, authLoading]);

  log('App rendering', { loggedIn, user, loaded, authLoading });

  if (!loaded) {
    return <div>Loading...</div>;
  }

  return (
    <div className="App">
      <NavBar 
        loggedIn={loggedIn} 
        user={user} 
        onSignIn={handleSignIn}
        onSignOut={handleSignOut}
        authLoading={authLoading}
      />
      <Main 
        loggedIn={loggedIn} 
        user={user}
      />
    </div>
  );
}

export default App;