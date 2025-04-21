import React, { useEffect, useState } from 'react';
import NavBar from './components/NavBar';
import Main from './components/Main';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

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

  // Initialize Google Auth
  useEffect(() => {
    const initGoogleAuth = async () => {
      try {
        log('Starting Google Auth initialization');
        if (!window.gapi) {
          log('Loading gapi library');
          const response = await fetch('https://apis.google.com/js/platform.js');
          if (!response.ok) throw new Error('Failed to load Google Auth library');
        }

        await window.gapi.load('auth2', () => {
          log('gapi.auth2 loaded');
          const clientId = 'ThisIsSupposedToBeAnId'; // Replace with actual client ID
          log('Initializing Google Auth with client ID', { clientId });

          window.gapi.auth2.init({
            client_id: clientId,
            scope: 'email',
            prompt: 'select_account',
            fetch_basic_profile: true,
          }).then(() => {
            log('Google Auth initialized');
            const auth2 = window.gapi.auth2.getAuthInstance();
            if (!auth2) {
              log('Failed to get auth2 instance');
              toast.error('Authentication setup failed.', { autoClose: 5000 });
              setLoaded(true);
              return;
            }

            // Check if user manually signed out
            const isSignedOut = localStorage.getItem('signedOut') === 'true';
            if (!isSignedOut && auth2.isSignedIn.get()) {
              const email = auth2.currentUser.get().getBasicProfile().getEmail();
              setUser(email);
              setLoggedIn(true);
              log('Session restored', { email });
            }

            // Listen for sign-in state changes
            auth2.isSignedIn.listen((isSignedIn) => {
              if (isSignedIn) {
                const email = auth2.currentUser.get().getBasicProfile().getEmail();
                setUser(email);
                setLoggedIn(true);
                localStorage.removeItem('signedOut');
                log('User signed in via listener', { email });
              } else {
                setUser(null);
                setLoggedIn(false);
                localStorage.setItem('signedOut', 'true');
                log('User signed out via listener');
              }
            });

            setLoaded(true);
          }).catch(error => {
            log('Google Auth initialization failed', { error: error.message || 'Unknown error' });
            toast.error('Failed to initialize authentication.', { autoClose: 5000 });
            setLoaded(true);
          });
        });
      } catch (error) {
        log('Error loading Google Auth', { error: error.message || 'Unknown error' });
        toast.error('Failed to load authentication library.', { autoClose: 5000 });
        setLoaded(true);
      }
    };

    initGoogleAuth();
  }, []);

  // Handle sign-in
  const handleSignIn = async () => {
    if (authLoading) {
      log('Sign-in attempt ignored: authLoading is true');
      return;
    }
    setAuthLoading(true);
    try {
      const auth2 = window.gapi.auth2.getAuthInstance();
      if (!auth2) throw new Error('Authentication instance not available');
      log('Initiating sign-in');
      await auth2.signIn();
      const email = auth2.currentUser.get().getBasicProfile().getEmail();
      setUser(email);
      setLoggedIn(true);
      localStorage.removeItem('signedOut');
      log('User signed in', { email });
      toast.success('Successfully signed in!', { autoClose: 3000, toastId: 'sign-in' });
    } catch (error) {
      log('Sign-in failed', { error: error.message || (error.error ? error.error : 'Unknown error') });
      const displayMessage = error.error === 'popup_closed_by_user'
        ? 'Sign-in canceled.'
        : error.error === 'access_denied'
        ? 'Permission denied.'
        : 'Failed to sign in.';
      toast.error(displayMessage, { autoClose: 5000 });
    } finally {
      setAuthLoading(false);
    }
  };

  // Handle sign-out
  const handleSignOut = async () => {
    if (authLoading) {
      log('Sign-out attempt ignored: authLoading is true');
      return;
    }
    setAuthLoading(true);
    try {
      const auth2 = window.gapi.auth2.getAuthInstance();
      if (!auth2) throw new Error('Authentication instance not available');
      log('Initiating sign-out');
      await auth2.signOut();
      setUser(null);
      setLoggedIn(false);
      localStorage.setItem('signedOut', 'true');
      log('Sign-out successful');
      toast.success('Successfully signed out.', { autoClose: 3000, toastId: 'sign-out' });
    } catch (error) {
      log('Sign-out failed', { error: error.message || (error.error ? error.error : 'Unknown error') });
      toast.error('Failed to sign out.', { autoClose: 5000 });
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
      <ToastContainer />
    </div>
  );
}

export default App;