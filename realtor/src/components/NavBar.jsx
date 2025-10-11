import { GoogleLogin } from '@react-oauth/google';
import { Nav, Navbar, NavDropdown, Container } from 'react-bootstrap';

const NavBar = ({ loggedIn, user, onLoginSuccess, onLoginError, onSignOut, authLoading }) => {
  return (
    <Navbar collapseOnSelect expand="lg" bg="dark" variant="dark" fixed="top">
      <Container>
        <Navbar.Brand href="/realtor">realtor webpage.</Navbar.Brand>
        <Navbar.Toggle aria-controls="responsive-navbar-nav" />
        <Navbar.Collapse>
          <Nav>
            <Nav.Link href="/realtor/search">Search Listings</Nav.Link>
            <Nav.Link href="/realtor/new">List Your Property</Nav.Link>
          </Nav>
          <Nav className="ml-auto">
            {loggedIn ? (
              <NavDropdown title={user} id="collapsible-nav-dropdown">
                <NavDropdown.Item href="/realtor/my-listings">View/Modify Listings</NavDropdown.Item>
                <NavDropdown.Item onClick={onSignOut} disabled={authLoading}>
                  {authLoading ? 'Signing Out...' : 'Sign Out'}
                </NavDropdown.Item>
              </NavDropdown>
            ) : (
              <div id="loginButton">
                <GoogleLogin
                  onSuccess={onLoginSuccess}
                  onError={onLoginError}
                  text="signin"
                  shape="rectangular"
                  theme="filled_blue"
                  size="medium"
                />
              </div>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default NavBar;