import React from "react";
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import NavDropdown from 'react-bootstrap/NavDropdown'
import { GoogleAPI } from 'react-google-oauth'

export default class NavBar extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            loggedIn: false,
            user: null
        };
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
                    <Nav.Link href="/realtor">Current Listings</Nav.Link>
                    <Nav.Link href="/realtor/new">List Your Property</Nav.Link>
                    </Nav>
                    <Nav>

                        {
                            !this.props.loggedIn && 
                            <Nav.Link href="/realtor/login">Login</Nav.Link>
                        }

                        {
                            this.props.loggedIn &&
                            <NavDropdown style={dropStyle} title={`${this.props.user}`} id="collasible-nav-dropdown">
                                <NavDropdown.Item href="/realtor/my-listings">View/Modify Your Listings</NavDropdown.Item>
                                <NavDropdown.Divider />
                                <NavDropdown.Item href="/realtor/support">Contact Seller Support</NavDropdown.Item>
                                <NavDropdown.Divider />
                                <NavDropdown.Item href="/realtor/logout">Logout</NavDropdown.Item>
                            </NavDropdown>
                        }
                    
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
        );
    }

}
