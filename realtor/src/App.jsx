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

  // Function to update auth state
  const updateAuthState = (isSignedIn) => {
    if (!window.gapi || !window.gapi.auth2) {
      log('Cannot update auth state: gapi.auth2 not available');
      return;
    }
    const auth2 = window.gapi.auth2.getAuthInstance();
    if (!auth2) {
      log('Cannot update auth state: auth2 instance is null');
      return;
    }
    if (isSignedIn) {
      const email = auth2.currentUser.get().getBasicProfile().getEmail();
      setUser(email);
      setLoggedIn(true);
      log('User signed in', { email });
      toast.success('Successfully signed in!', { autoClose: 3000 });
    } else {
      setUser(null);
      setLoggedIn(false);
      log('User signed out');
      toast.success('Successfully signed out.', { autoClose: 3000 });
    }
  };

  // Auth initialization (runs once)
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
          const clientId = process.env.REACT_APP_GOOGLE_CLIENT_ID || 'ThisIsSupposedToBeAnId';
          log('Initializing Google Auth with client ID', { clientId });

          window.gapi.auth2.init({
            client_id: clientId,
            scope: 'email',
            prompt: 'select_account'
          }).then(() => {
            log('Google Auth initialized');
            const auth2 = window.gapi.auth2.getAuthInstance();
            if (!auth2) {
              log('Failed to get auth2 instance');
              toast.error('Authentication setup failed. Please try again later.', { autoClose: 5000 });
              setLoaded(true);
              return;
            }
            const signedIn = auth2.isSignedIn.get();
            updateAuthState(signedIn);

            // Listen for sign-in state changes
            const listener = auth2.isSignedIn.listen(updateAuthState);
            log('Sign-in state listener added');
            setLoaded(true);

            // Clean up listener on unmount
            return () => {
              if (listener) {
                log('Cleaning up auth listener');
              }
            };
          }).catch(error => {
            log('Google Auth initialization failed', { error: error.message, details: error });
            toast.error('Failed to initialize authentication. Please try again later.', { autoClose: 5000 });
            setLoaded(true);
          });
        });
      } catch (error) {
        log('Error loading Google Auth', { error: error.message });
        toast.error('Failed to load authentication library. Please check your network and try again.', { autoClose: 5000 });
        setLoaded(true);
      }
    };

    initGoogleAuth();
  }, []);

  // Handle sign-in with debounce and error handling
  const handleSignIn = async () => {
    if (authLoading) {
      log('Sign-in attempt ignored: authLoading is true');
      return;
    }
    setAuthLoading(true);
    try {
      if (!window.gapi || !window.gapi.auth2) {
        throw new Error('Google Auth library not loaded');
      }
      const auth2 = window.gapi.auth2.getAuthInstance();
      if (!auth2) {
        throw new Error('Authentication instance not available');
      }
      log('Initiating sign-in');
      await auth2.signIn();
      updateAuthState(true);
    } catch (error) {
      log('Sign-in failed', { error: error.message, details: error });
      let errorMessage = 'Failed to sign in. Please try again.';
      if (error.error === 'popup_closed_by_user') {
        errorMessage = 'Sign-in canceled. Please complete the sign-in process.';
      } else if (error.error === 'access_denied') {
        errorMessage = 'Permission denied. Please grant the required permissions to sign in.';
      } else if (error.error === 'invalid_client') {
        errorMessage = 'Invalid authentication configuration. Please contact support.';
      } else if (error.error === 'idpiframe_initialization_failed') {
        errorMessage = 'Authentication configuration error. Please contact support.';
      } else if (error.message === 'Google Auth library not loaded' || error.message === 'Authentication instance not available') {
        errorMessage = 'Authentication service unavailable. Please try again later.';
      } else if (error.message.includes('null')) {
        errorMessage = 'Authentication error: Service not initialized. Please try again.';
      }
      toast.error(errorMessage, { autoClose: 5000 });
    } finally {
      setAuthLoading(false);
    }
  };

  // Handle sign-out with debounce and error handling
  const handleSignOut = async () => {
    if (authLoading) {
      log('Sign-out attempt ignored: authLoading is true');
      return;
    }
    setAuthLoading(true);
    try {
      if (!window.gapi || !window.gapi.auth2) {
        throw new Error('Google Auth library not loaded');
      }
      const auth2 = window.gapi.auth2.getAuthInstance();
      if (!auth2) {
        throw new Error('Authentication instance not available');
      }
      log('Initiating sign-out');
      await auth2.signOut();
      updateAuthState(false);
    } catch (error) {
      log('Sign-out failed', { error: error.message, details: error });
      toast.error('Failed to sign out. Please try again.', { autoClose: 5000 });
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