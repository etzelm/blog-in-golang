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
            loggedIn: props ? props.loggedIn : null,
            user: props ? props.user : null,
            cards: [],
            orgCards: []
        };
        this.onSubmit = this.onSubmit.bind(this);
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
            const listings = data.filter(card => card.deleted !== "false");
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
        const fieldsToCheck = [];
        const { orgCards } = this.state;
        
        if (this.cityRef.current.value) {
            fieldsToCheck.push({field: "City", value: this.cityRef.current.value});
        }
        if (this.stateRef.current.value) {
            fieldsToCheck.push({field: "State", value: this.stateRef.current.value});
        }
        if (this.zipCodeRef.current.value) {
            fieldsToCheck.push({field: "ZipCode", value: this.zipCodeRef.current.value});
        }
        if (this.bedroomsRef.current.value) {
            fieldsToCheck.push({field: "Bedrooms", value: this.bedroomsRef.current.value});
        }
        if (this.squareFeetRef.current.value) {
            fieldsToCheck.push({field: "SquareFeet", value: this.squareFeetRef.current.value});
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
            height: "240vh"
        }

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
            position: "absolute",
            left: "50%"
        };

        return (
            <div style={homeStyle}>
                <br/><br/><br/>
                <Card style={cardStyle}>
                    <Form onSubmit={this.onSubmit}>
                        <Row>
                            <Form.Group as={Col} controlId="formGridCity">
                            <Form.Label>City</Form.Label>
                            <Form.Control type="text" ref={this.cityRef}/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridState">
                            <Form.Label>State</Form.Label>
                            <Form.Control type="text" ref={this.stateRef}/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridZip">
                            <Form.Label>Zip Code</Form.Label>
                            <Form.Control type="text" ref={this.zipCodeRef}/>
                            </Form.Group>
                        </Row>

                        <Row>
                            <Form.Group as={Col} controlId="formGridBedrooms">
                            <Form.Label>Bedrooms</Form.Label>
                            <Form.Control type="text" ref={this.bedroomsRef}/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridBathrooms">
                            <Form.Label>Bathrooms</Form.Label>
                            <Form.Control type="text" ref={this.bathroomsRef}/>
                            </Form.Group>
                        </Row>

                        <Row>
                            <Form.Group as={Col} controlId="formGridMLS">
                            <Form.Label>MLS</Form.Label>
                            <Form.Control type="text" ref={this.mlsRef}/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridSquareFeet">
                            <Form.Label>Square Feet</Form.Label>
                            <Form.Control type="text" ref={this.squareFeetRef}/>
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
                <TileDeck cards={this.state.cards}/>
            </div>
        );
   
    }

}
