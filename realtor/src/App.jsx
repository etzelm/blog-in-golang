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
  const [hasInteracted, setHasInteracted] = useState(false);
  const [justSignedIn, setJustSignedIn] = useState(false);

  // Function to update auth state
  const updateAuthState = (isSignedIn, isSignInAction = false) => {
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
      if (isSignInAction) {
        setJustSignedIn(true);
        toast.success('Successfully signed in!', { autoClose: 3000 });
      }
    } else {
      setUser(null);
      setLoggedIn(false);
      log('User signed out');
      if (isSignInAction) {
        toast.success('Successfully signed out.', { autoClose: 3000 });
      }
    }
  };

  // Clear justSignedIn flag after showing the toast
  useEffect(() => {
    if (justSignedIn) {
      const timer = setTimeout(() => {
        setJustSignedIn(false);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [justSignedIn]);

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
              if (hasInteracted) {
                toast.error('Authentication setup failed. Please try again later.', { autoClose: 5000 });
              }
              setLoaded(true);
              return;
            }
            // Do not automatically restore session on page load
            setLoaded(true);

            // Listen for sign-in state changes
            const listener = auth2.isSignedIn.listen((isSignedIn) => {
              updateAuthState(isSignedIn, true);
            });
            log('Sign-in state listener added');

            // Clean up listener on unmount
            return () => {
              if (listener) {
                log('Cleaning up auth listener');
              }
            };
          }).catch(error => {
            log('Google Auth initialization failed', { error: error.message || 'Unknown error', details: error });
            if (hasInteracted) {
              toast.error('Failed to initialize authentication. Please try again later.', { autoClose: 5000 });
            }
            setLoaded(true);
          });
        });
      } catch (error) {
        log('Error loading Google Auth', { error: error.message || 'Unknown error' });
        if (hasInteracted) {
          toast.error('Failed to load authentication library. Please check your network and try again.', { autoClose: 5000 });
        }
        setLoaded(true);
      }
    };

    initGoogleAuth();
  }, [hasInteracted]);

  // Clear cookies on sign-out
  const clearCookies = () => {
    document.cookie = 'user=; Max-Age=0; path=/; domain=mitchelletzel.com';
    document.cookie = 'userToken=; Max-Age=0; path=/; domain=mitchelletzel.com';
    log('Cookies cleared', { cookies: ['user', 'userToken'] });
  };

  // Handle sign-in with debounce and error handling
  const handleSignIn = async () => {
    if (authLoading) {
      log('Sign-in attempt ignored: authLoading is true');
      return;
    }
    setAuthLoading(true);
    setHasInteracted(true);
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
      updateAuthState(true, true);
    } catch (error) {
      const errorMessage = error.message || (error.error ? error.error : 'Unknown error');
      log('Sign-in failed', { error: errorMessage, details: error });
      let displayMessage = 'Failed to sign in. Please try again.';
      if (error.error === 'popup_closed_by_user') {
        displayMessage = 'Sign-in canceled. Please complete the sign-in process.';
      } else if (error.error === 'access_denied') {
        displayMessage = 'Permission denied. Please grant the required permissions to sign in.';
      }
      toast.error(displayMessage, { autoClose: 5000 });
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
    setHasInteracted(true);
    try {
      if (!window.gapi || !window.gapi.auth2) {
        throw new Error('Google Auth library not loaded');
      }
      const auth2 = window.gapi.auth2.getAuthInstance();
      if (!auth2) {
        throw new Error('Authentication instance not available');
      }
      log('Initiating sign-out');
      await auth2.disconnect(); // Use disconnect to revoke access
      clearCookies(); // Clear application cookies
      updateAuthState(false, true);
    } catch (error) {
      const errorMessage = error.message || (error.error ? error.error : 'Unknown error');
      log('Sign-out failed', { error: errorMessage, details: error });
      toast.error('Failed to sign out. Please try again.', { autoClose: 5000 });
    } finally {
      setAuthLoading(false);
    }
  };

  // Log state changes
  useEffect(() => {
    log('App state changed', { loggedIn, user, loaded, authLoading, hasInteracted, justSignedIn });
  }, [loggedIn, user, loaded, authLoading, hasInteracted, justSignedIn]);

  log('App rendering', { loggedIn, user, loaded, authLoading, hasInteracted, justSignedIn });

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