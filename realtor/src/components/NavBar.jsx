import React from "react";
import { Nav, Navbar, NavDropdown, Container } from 'react-bootstrap';

const NavBar = ({ loggedIn, user, onSignIn, onSignOut, authLoading }) => {
  return (
    <Navbar collapseOnSelect expand="lg" bg="dark" variant="dark" fixed="top">
      <Container>
        <Navbar.Brand href="/realtor">realtor webpage.</Navbar.Brand>
        <Navbar.Toggle aria-controls="responsive-navbar-nav" />
        <Navbar.Collapse>
          <Nav className="me-auto">
            <Nav.Link href="/realtor/search">Search Listings</Nav.Link>
            <Nav.Link href="/realtor/new">List Your Property</Nav.Link>
            <Nav>
              {loggedIn ? (
                <NavDropdown title={user} id="collapsible-nav-dropdown">
                  <NavDropdown.Item href="/realtor/my-listings">View/Modify Listings</NavDropdown.Item>
                  <NavDropdown.Item onClick={onSignOut} disabled={authLoading}>
                    {authLoading ? 'Signing Out...' : 'Sign Out'}
                  </NavDropdown.Item>
                </NavDropdown>
              ) : (
                <button 
                  id="loginButton"
                  onClick={onSignIn}
                  className="btn btn-light"
                  disabled={authLoading}
                >
                  {authLoading ? 'Signing In...' : 'Sign In'}
                </button>
              )}
            </Nav>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default NavBar;