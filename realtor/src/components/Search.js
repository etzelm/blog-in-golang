import React from "react";
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Button from 'react-bootstrap/Button';
import TileDeck from "./TileDeck";

export default class Search extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            loggedIn: props.loggedIn || null,
            user: props.user || null,
            cards: [],
            orgCards: []
        };
        
        // Create refs
        this.cityRef = React.createRef();
        this.stateRef = React.createRef();
        this.zipCodeRef = React.createRef();
        this.bedroomsRef = React.createRef();
        this.bathroomsRef = React.createRef();
        this.mlsRef = React.createRef();
        this.squareFeetRef = React.createRef();
    }

    async componentDidMount() {
        try {
            const response = await fetch('/listings');
            const data = await response.json();
            // Fix: The filter condition was inverted (checking !== "false" instead of === true)
            const listings = data.filter(card => card.deleted !== true);
            this.setState({ 
                cards: listings,
                orgCards: listings
            });
        } catch (error) {
            console.error("Error fetching listings:", error);
        }
    }

    // Remove redundant onSubmit binding in constructor since it's defined as arrow function
    onSubmit = async (event) => {
        event.preventDefault();
        const { orgCards } = this.state;
        const fieldsToCheck = [];
        
        // Add filters only if values exist
        if (this.cityRef.current?.value) {
            fieldsToCheck.push({ field: "City", value: this.cityRef.current.value });
        }
        if (this.stateRef.current?.value) {
            fieldsToCheck.push({ field: "State", value: this.stateRef.current.value });
        }
        if (this.zipCodeRef.current?.value) {
            fieldsToCheck.push({ field: "ZipCode", value: this.zipCodeRef.current.value });
        }
        if (this.bedroomsRef.current?.value) {
            fieldsToCheck.push({ field: "Bedrooms", value: Number(this.bedroomsRef.current.value) });
        }
        if (this.bathroomsRef.current?.value) {
            fieldsToCheck.push({ field: "Bathrooms", value: Number(this.bathroomsRef.current.value) });
        }
        if (this.mlsRef.current?.value) {
            fieldsToCheck.push({ field: "MLS", value: this.mlsRef.current.value });
        }
        if (this.squareFeetRef.current?.value) {
            fieldsToCheck.push({ field: "SquareFeet", value: Number(this.squareFeetRef.current.value) });
        }

        const filteredCards = orgCards.filter(card => 
            fieldsToCheck.every(filter => card[filter.field] === filter.value)
        );
        this.setState({ cards: filteredCards });
    }

    render() {
        const homeStyle = {
            backgroundColor: 'LightGray',
            margin: "0px",
            padding: "0px",
            minHeight: "100vh" // Changed from fixed height to minimum height
        };

        const cardStyle = {
            width: '82vw',
            paddingTop: '1vw',
            paddingLeft: '2vw',
            paddingRight: '2vw',
            paddingBottom: '5vw',
            margin: 'auto',
            backgroundColor: 'White'
        };

        const buttonStyle = {
            margin: "0",
            position: "relative", // Changed from absolute to relative for better form layout
            left: "50%",
            transform: "translateX(-50%)" // Added to center the button
        };

        return (
            <div style={homeStyle}>
                <br/><br/><br/>
                <Card style={cardStyle}>
                    <Form onSubmit={this.onSubmit}>
                        <Row>
                            <Form.Group as={Col} controlId="formGridCity">
                                <Form.Label>City</Form.Label>
                                <Form.Control type="text" ref={this.cityRef} />
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridState">
                                <Form.Label>State</Form.Label>
                                <Form.Control type="text" ref={this.stateRef} />
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridZip">
                                <Form.Label>Zip Code</Form.Label>
                                <Form.Control type="text" ref={this.zipCodeRef} />
                            </Form.Group>
                        </Row>

                        <Row className="mb-3">
                            <Form.Group as={Col} controlId="formGridBedrooms">
                                <Form.Label>Bedrooms</Form.Label>
                                <Form.Control type="number" ref={this.bedroomsRef} />
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridBathrooms">
                                <Form.Label>Bathrooms</Form.Label>
                                <Form.Control type="number" ref={this.bathroomsRef} />
                            </Form.Group>
                        </Row>

                        <Row className="mb-3">
                            <Form.Group as={Col} controlId="formGridMLS">
                                <Form.Label>MLS</Form.Label>
                                <Form.Control type="text" ref={this.mlsRef} />
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridSquareFeet">
                                <Form.Label>Square Feet</Form.Label>
                                <Form.Control type="number" ref={this.squareFeetRef} />
                            </Form.Group>
                        </Row><br/>
                
                        <Button 
                            style={buttonStyle} 
                            variant="primary" 
                            type="submit"
                        >
                            Submit
                        </Button>
                    </Form>
                </Card>
                <TileDeck cards={this.state.cards} />
            </div>
        );
    }
}