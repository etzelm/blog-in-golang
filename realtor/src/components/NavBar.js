import React from "react";
import { Nav, Navbar, NavDropdown, Container } from 'react-bootstrap';

const NavBar = ({ loggedIn, user, signOut }) => {
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
                            <NavDropdown.Item onClick={signOut}>Sign Out</NavDropdown.Item>
                        </NavDropdown>
                    ) : (
                        <div id="loginButton" />
                    )}
                </Nav>
            </Nav>
        </Navbar.Collapse>
        </Container>
    </Navbar>
  );
};

export default NavBar;