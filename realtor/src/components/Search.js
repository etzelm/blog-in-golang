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
            orgCards: [],
            noResults: false
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
            console.log('Raw API response:', data);
            const listings = data.filter(card => card.deleted === "false" || card.deleted === false);
            const excludedListings = data.filter(card => !(card.deleted === "false" || card.deleted === false));
            console.log('Filtered listings (deleted=false):', listings);
            console.log('Excluded listings (deleted=true):', excludedListings);
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
            City: this.cityRef.current?.value?.trim().toLowerCase() || null,
            State: this.stateRef.current?.value?.trim().toLowerCase() || null,
            "Zip Code": this.zipCodeRef.current?.value?.trim() || null,
            Bedrooms: this.bedroomsRef.current?.value ? Number(this.bedroomsRef.current.value) : null,
            Bathrooms: this.bathroomsRef.current?.value ? Number(this.bathroomsRef.current.value) : null,
            MLS: this.cityRef.current?.value?.trim() || null,
            "Square Feet": this.squareFeetRef.current?.value ? Number(this.squareFeetRef.current.value) : null,
        };

        console.log('Filter values:', filters);

        // Filter cards based on provided values
        const filteredCards = orgCards.filter((card, index) => {
            const result = Object.keys(filters).every((field) => {
                // Skip if filter value is empty or null
                if (!filters[field]) return true;

                // Normalize card data for comparison
                let cardValue = card[field];
                let filterValue = filters[field];

                // Handle string fields (City, State, MLS)
                if (['City', 'State', 'MLS'].includes(field)) {
                    cardValue = cardValue ? String(cardValue).toLowerCase() : '';
                    filterValue = String(filterValue).toLowerCase();
                }

                // Handle Zip Code (strip extra characters like - for ZIP+4)
                if (field === 'Zip Code') {
                    cardValue = cardValue ? String(cardValue).replace(/[^0-9]/g, '') : '';
                    filterValue = String(filterValue).replace(/[^0-9]/g, '');
                }

                // Handle numeric fields (Bedrooms, Bathrooms, Square Feet)
                if (['Bedrooms', 'Bathrooms', 'Square Feet'].includes(field)) {
                    cardValue = cardValue ? Number(cardValue) : null;
                    filterValue = Number(filterValue);
                }

                const match = cardValue === filterValue;
                if (!match) {
                    console.log(`Card ${index} failed filter: ${field} (card: ${cardValue}, filter: ${filterValue})`);
                }
                return match;
            });

            return result;
        });

        console.log('Filtered cards:', filteredCards);
        this.setState({ 
            cards: filteredCards,
            noResults: filteredCards.length === 0
        });
    }

    resetForm = () => {
        // Clear form inputs
        this.cityRef.current.value = '';
        this.stateRef.current.value = '';
        this.zipCodeRef.current.value = '';
        this.bedroomsRef.current.value = '';
        this.bathroomsRef.current.value = '';
        this.mlsRef.current.value = '';
        this.squareFeetRef.current.value = '';
        // Restore original cards
        this.setState({ 
            cards: this.state.orgCards,
            noResults: false
        });
        console.log('Form reset, restored original cards');
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
            margin: "0 10px",
            position: "relative",
            display: "inline-block"
        };

        const buttonContainerStyle = {
            textAlign: "center",
            marginTop: "20px"
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
                        </Row>

                        <div style={buttonContainerStyle}>
                            <Button 
                                style={buttonStyle} 
                                variant="primary" 
                                type="submit"
                            >
                                Submit
                            </Button>
                            <Button 
                                style={buttonStyle} 
                                variant="secondary" 
                                onClick={this.resetForm}
                            >
                                Reset
                            </Button>
                        </div>
                    </Form>
                </Card>
                {this.state.noResults && (
                    <div style={{ textAlign: 'center', marginTop: '20px', color: 'red' }}>
                        No listings match your criteria.
                    </div>
                )}
                <TileDeck cards={this.state.cards} />
            </div>
        );
    }
}