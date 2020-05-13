import React from "react";
import Card from 'react-bootstrap/Card'
import Col from 'react-bootstrap/Col'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import TileDeck from "./TileDeck";

export default class Search extends React.Component {

    constructor(props) {
        super(props);
        this.onSubmit = this.onSubmit.bind(this);
        this.state = {
            cards: [],
            orgCards: []
        };
    }

    async componentDidMount() {

        const response = await fetch('/listings');
        const data = await response.json();
        console.log(data)

        var listings = [];
        for (var it=0; it<data.length; it++) {
            const card = data[it];
            if (card['deleted'] === "false") {
                listings.push(card);
            }
        }

        this.setState({ 
            cards: listings,
            orgCards: listings
        })

    }

    async onSubmit(event) {
        console.log("submitted")
        console.log(event)
        console.log(event.currentTarget)
        console.log(event.currentTarget.elements.ZipCode.value)
        console.log(event.target.elements.ZipCode.value)
        event.preventDefault()

        var newCards = [], fieldsA = [], fieldsB = [];
        const elements = event.currentTarget.elements;

        if (!(elements.City.value === "")) {
            fieldsA.push("City");
            fieldsB.push("City");
        }
        if (!(elements.State.value === "")) {
            fieldsA.push("State");
            fieldsB.push("State");
        }
        if (!(elements.ZipCode.value === "")) {
            fieldsA.push("ZipCode");
            fieldsB.push("Zip Code");
        }
        if (!(elements.Bedrooms.value === "")) {
            fieldsA.push("Bedrooms");
            fieldsB.push("Bedrooms");
        }
        if (!(elements.Bathrooms.value === "")) {
            fieldsA.push("Bathrooms");
            fieldsB.push("Bathrooms");
        }
        if (!(elements.MLS.value === "")) {
            fieldsA.push("MLS");
            fieldsB.push("MLS");
        }
        if (!(elements.SquareFeet.value === "")) {
            fieldsA.push("SquareFeet");
            fieldsB.push("SquareFeet");
        }

        console.log(fieldsA)
        for (var it=0; it < this.state.orgCards.length; it++) {

            var match = true;
            for (var innerIt=0; innerIt < fieldsA.length; innerIt++) {

                if (this.state.orgCards[it][fieldsB[innerIt]] !== elements[fieldsA[innerIt]].value) {
                    match = false;
                }

            }

            if (match) {
                newCards.push(this.state.orgCards[it]);
            }

        }

        this.setState({ cards: newCards })

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

                        <Form.Row>
                            <Form.Group as={Col} controlId="formGridCity">
                            <Form.Label>City</Form.Label>
                            <Form.Control type="text" name="City" ref="City"/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridState">
                            <Form.Label>State</Form.Label>
                            <Form.Control type="text" name="State" ref="State"/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridZipCode">
                            <Form.Label>Zip Code</Form.Label>
                            <Form.Control type="text" name="ZipCode" ref="ZipCode"/>
                            </Form.Group>
                        </Form.Row>

                        <Form.Row>
                            <Form.Group as={Col} controlId="formGridBedrooms">
                            <Form.Label>Bedrooms</Form.Label>
                            <Form.Control type="text" name="Bedrooms" ref="Bedrooms"/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridBathrooms">
                            <Form.Label>Bathrooms</Form.Label>
                            <Form.Control type="text" name="Bathrooms" ref="Bathrooms"/>
                            </Form.Group>
                        </Form.Row>

                        <Form.Row>
                            <Form.Group as={Col} controlId="formGridMLS">
                            <Form.Label>MLS</Form.Label>
                            <Form.Control type="text" name="MLS" ref="MLS"/>
                            </Form.Group>

                            <Form.Group as={Col} controlId="formGridSquareFeet">
                            <Form.Label>Square Feet</Form.Label>
                            <Form.Control type="text" name="SquareFeet" ref="SquareFeet"/>
                            </Form.Group>
                        </Form.Row>
                        
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
