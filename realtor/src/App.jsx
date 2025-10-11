import { useEffect, useState } from 'react';
import { GoogleOAuthProvider } from '@react-oauth/google';
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

// Replace with your actual Google OAuth client ID
const GOOGLE_CLIENT_ID = 'ThisIsSupposedToBeAnId';

function App() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  const [loaded, setLoaded] = useState(true); // No async loading needed with new library
  const [authLoading, setAuthLoading] = useState(false);

  // Check for existing session on component mount
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    const isSignedOut = localStorage.getItem('signedOut') === 'true';
    
    if (storedUser && !isSignedOut) {
      setUser(storedUser);
      setLoggedIn(true);
      log('Session restored', { email: storedUser });
    }
  }, []);

  // Handle successful login
  const handleLoginSuccess = (credentialResponse) => {
    try {
      log('Login successful, processing credential response');
      
      // Decode the JWT token to get user info
      const payload = JSON.parse(atob(credentialResponse.credential.split('.')[1]));
      const email = payload.email;
      
      setUser(email);
      setLoggedIn(true);
      localStorage.setItem('user', email);
      localStorage.removeItem('signedOut');
      
      log('User signed in', { email });
      toast.success('Successfully signed in!', { autoClose: 3000, toastId: 'sign-in' });
    } catch (error) {
      log('Login processing failed', { error: error.message || 'Unknown error' });
      toast.error('Failed to process login.', { autoClose: 5000 });
    }
  };

  // Handle login error
  const handleLoginError = (error) => {
    log('Login failed', { error });
    toast.error('Failed to sign in.', { autoClose: 5000 });
  };

  // Handle sign-out
  const handleSignOut = () => {
    if (authLoading) {
      log('Sign-out attempt ignored: authLoading is true');
      return;
    }
    
    setAuthLoading(true);
    try {
      setUser(null);
      setLoggedIn(false);
      localStorage.removeItem('user');
      localStorage.setItem('signedOut', 'true');
      
      log('Sign-out successful');
      toast.success('Successfully signed out.', { autoClose: 3000, toastId: 'sign-out' });
    } catch (error) {
      log('Sign-out failed', { error: error.message || 'Unknown error' });
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
    <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
      <div className="App">
        <NavBar 
          loggedIn={loggedIn} 
          user={user} 
          onLoginSuccess={handleLoginSuccess}
          onLoginError={handleLoginError}
          onSignOut={handleSignOut}
          authLoading={authLoading}
        />
        <Main 
          loggedIn={loggedIn} 
          user={user}
        />
        <ToastContainer />
      </div>
    </GoogleOAuthProvider>
  );
}

export default App;