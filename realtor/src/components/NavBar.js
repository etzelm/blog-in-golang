import React from "react";
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import NavDropdown from 'react-bootstrap/NavDropdown'

export default class NavBar extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            loggedIn: this.props == null ? null : this.props.loggedIn,
            user: this.props == null ? null : this.props.user,
            loggedOut: false
        };
        this.signOut = this.signOut.bind(this);
    }

    async signOut() {
        await this.props.signOut();
        console.log("got back")
        window.gapi.load('auth2', () => {
            this.auth2 = window.gapi.auth2.init({
              client_id: 'ThisIsSupposedToBeAnId',
            })
      
            this.auth2.then(() => {
              console.log('on init');
              this.setState({
                loggedIn: this.auth2.isSignedIn.get(),
              });
            });
        });
        window.gapi.load('signin2', function() {
            // Method 3: render a sign in button
            // using this method will show Signed In if the user is already signed in
            var opts = {
              width: 100,
              height: 25,
              client_id: 'ThisIsSupposedToBeAnId'
            }
            window.gapi.signin2.render('loginButton', opts)
        })
    }

    render() {
        
        const dropStyle = {
            paddingRight: '5vw'
        };

        return (
            <Navbar collapseOnSelect expand="lg" bg="dark" variant="dark" fixed="top" >
                <Navbar.Brand href="/realtor">realtor webpage.</Navbar.Brand>
                <Navbar.Toggle aria-controls="responsive-navbar-nav" />
                <Navbar.Collapse id="responsive-navbar-nav">
                    <Nav className="mr-auto">
                    <Nav.Link href="/realtor/search">Search Listings</Nav.Link>
                    <Nav.Link href="/realtor/new">List Your Property</Nav.Link>
                    </Nav>
                    <Nav>

                        {
                            !this.props.loggedIn && !this.state.loggedOut && 
                            <Nav.Link>
                                <button id="loginButton">Login with Google</button>
                            </Nav.Link>
                        }

                        {
                            this.props.loggedIn && !this.state.loggedOut && 
                            <NavDropdown style={dropStyle} title={`${this.props.user}`} id="collasible-nav-dropdown">
                                <NavDropdown.Item href="/realtor/my-listings">View/Modify Your Listings</NavDropdown.Item>
                                <NavDropdown.Divider />
                                <NavDropdown.Item onClick={this.signOut}>Sign out</NavDropdown.Item>
                            </NavDropdown>
                        }
                    
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
        );
    }

}
