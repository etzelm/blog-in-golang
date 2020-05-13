import React from 'react';
import { withRouter } from 'react-router-dom';
import Main from './components/Main';
import NavBar from './components/NavBar';

class App extends React.Component {

  constructor(props) {
    super(props);
    this.signOut = this.signOut.bind(this);

    const token = localStorage.getItem('aToken');
    var loggedIn = false, user = null;
    console.log(token);
    if (!(token == null)) {
      loggedIn = true
      user = token
    }

    this.state = {
      loggedIn,
      loggedOut: false,
      loaded: false,
      reload: true,
      user
    }
    
  }

  componentDidMount() {

    this.setState({
      loaded: false
    });

    const successCallback = this.onSuccess.bind(this);

    if (!(window.gapi == null) && !this.state.loggedIn) {

      window.gapi.load('auth2', () => {
        this.auth2 = window.gapi.auth2.init({
          client_id: 'ThisIsSupposedToBeAnId',
        })
  
        this.auth2.then(() => {
          console.log('on init');
  
          const loggedIn = this.auth2.isSignedIn.get();
          var email = null;
          if (loggedIn) {
            email = this.auth2.currentUser.get().getBasicProfile().getEmail();
            console.log(email);
          }
          
          this.setState({
            loggedIn: loggedIn,
            user: email,
            loaded: true
          });
        });
      });    

      window.gapi.load('signin2', function() {
        // Method 3: render a sign in button
        // using this method will show Signed In if the user is already signed in
        var opts = {
          width: 100,
          height: 25,
          client_id: 'ThisIsSupposedToBeAnId',
          onsuccess: successCallback
        }
        window.gapi.signin2.render('loginButton', opts)
      })

    }

    this.setState({
      loaded: true,
      reload: true
    });

  }

  onSuccess() {
    console.log('on success')
    const path = this.props.location.pathname
    if (path === "/realtor/my-listings" || path === "/realtor/new") {
      this.setState({reload: false})
    }
    const email = this.auth2.currentUser.get().getBasicProfile().getEmail();
    localStorage.setItem('aToken', email)

    this.setState({
      loggedIn: true,
      user: email,
      reload: true,
      err: null
    })

  }

  onLoginFailed(err) {
    this.setState({
      loggedIn: false,
      user: null,
      error: err
    })
  }

  signOut() {
    window.gapi.load('auth2', () => {
      this.auth2 = window.gapi.auth2.init({
        client_id: 'ThisIsSupposedToBeAnId',
      })

      this.auth2.signOut().then(this.auth2.disconnect().then(function () {
        console.log('User signed out.');
      }));
    });
    localStorage.removeItem("aToken");
    this.setState({
      loggedIn: false,
      user: null,
      loggedOut: true
    });
    const path = this.props.location.pathname
    if (path === "/realtor/my-listings" || path === "/realtor/new") {
      this.props.history.go("/realtor");
    }
  }

  render() {
    return (
      <div className="App">
      
        {

          this.state.loggedIn && this.state.loaded && 
          <NavBar 
            loggedIn={this.state["loggedIn"]} 
            loggedOut={this.state["loggedOut"]} 
            user={this.state["user"]} 
            signOut={this.signOut}
          />

        }

        {

          !this.state.loggedIn && this.state.loaded && 
          <NavBar loggedIn={false} user={null} />

        }
        
        {

          this.state.reload &&
          <Main loggedIn={this.state["loggedIn"]} user={this.state["user"]} />

        }

      </div>
    );
  }
}

export default withRouter(App);
