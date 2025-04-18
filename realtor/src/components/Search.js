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
            const listings = data.filter(card => card.deleted === "false");
            this.setState({ 
                cards: listings,
                orgCards: listings
            });
        } catch (error) {
            console.error("Error fetching listings:", error);
        }
    }

    onSubmit = async (event) => {
        event.preventDefault();
        const { orgCards } = this.state;

        // Get filter values and normalize them
        const filters = {
            City: this.cityRef.current?.value?.trim().toLowerCase(),
            State: this.stateRef.current?.value?.trim().toLowerCase(),
            ZipCode: this.zipCodeRef.current?.value?.trim(),
            Bedrooms: this.bedroomsRef.current?.value ? Number(this.bedroomsRef.current.value) : null,
            Bathrooms: this.bathroomsRef.current?.value ? Number(this.bathroomsRef.current.value) : null,
            MLS: this.mlsRef.current?.value?.trim(),
            SquareFeet: this.squareFeetRef.current?.value ? Number(this.squareFeetRef.current.value) : null,
        };

        // Filter cards based on provided values
        const filteredCards = orgCards.filter((card) => {
            return Object.keys(filters).every((field) => {
                // Skip if filter value is empty or null
                if (!filters[field]) return true;

                // Normalize card data for comparison
                const cardValue = typeof card[field] === 'string' ? card[field].toLowerCase() : card[field];
                const filterValue = typeof filters[field] === 'string' ? filters[field] : filters[field];

                // Special handling for ZipCode to ensure string comparison
                if (field === 'ZipCode') {
                    return String(card[field]) === String(filterValue);
                }

                return cardValue === filterValue;
            });
        });

        this.setState({ cards: filteredCards });
    }

    render() {
        const homeStyle = {
            backgroundColor: 'LightGray',
            margin: "0px",
            padding: "0px",
            minHeight: "100vh"
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
            position: "relative",
            left: "50%",
            transform: "translateX(-50%)"
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