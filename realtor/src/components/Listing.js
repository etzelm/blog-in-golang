import React from "react";
import Card from 'react-bootstrap/Card'
import Carousel from 'react-bootstrap/Carousel'
import ListGroup from 'react-bootstrap/ListGroup'
import ListGroupItem from 'react-bootstrap/ListGroupItem'
import Tile from "./Tile";

export default class Listing extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            MLS: ""
        };
    }

    async componentDidMount() {

        const search = this.props.location.search;
        const regex = /(?:\x3d)([^\x26]*)/i;
        const found = search.match(regex);

        const rootDomain = `${window.location.protocol}//${window.location.host}`.replace('/realtor', '');
        const fullPath = `${rootDomain}/listing/${found[1]}`;
        console.log("Full Listing Domain: ", fullPath);
        const response = await fetch(fullPath);
        const data = await response.json();

        if (data.length > 0 ) {
            this.setState({ card: data[0] })
        }

    }

    render() {
        const listingStyle = {
            backgroundColor: 'Gray',
            margin: "0px",
            padding: "0px",
            height: "240vh"
        }

        const cardStyle = {
            width: '90vw',
            paddingTop: '3vw',
            paddingLeft: '2vw',
            paddingRight: '2vw',
            paddingBottom: '3vw',
            margin: 'auto',
            backgroundColor: 'LightGray'
        };

        const carouselStyle = {
            width: '70vw',
            height: '25vw',
            margin: 'auto',
            paddingBottom: '12vw',
            borderStyle: "solid",
            borderWidth: "8px"
        };

        const itemStyle = {
            backgroundSize: 'auto',
            objectFit: 'cover',
            width: '100%',
            height: '24vw',
            overflow: 'hidden',
            alignItems: 'center'
        };

        const card = this.state.card
        const photos = !(card == null) ? card['Photo Array'] : null;
        var addr = null, desc1 = null, desc2 = null, price = null, 
            ago = null, listed = null, desc3 = null, d = new Date(0);

        if (!(card == null)) {
            const street = (card['Street2'] !== "*") ?
                        `${card['Street1']}, ${card['Street2']} | ` :
                        `${card['Street1']} | `;
            addr = street + 
                        `${card['City']}, ${card['State']} ` + 
                        `${card['Zip Code']}`; 
            desc1 = `Square Feet: ${card['Square Feet']} | ` + 
                        `Lot Size: ${card['Lot Size']}`;
            desc2 = `Beds: ${card['Bedrooms']} | ` + 
                        `Baths: ${card['Bathrooms']}`;
            const time = new Date().getTime();
            ago = new Tile().timeDifference(time, card['Last Modified']);
            price = `Price: $${card['Sales Price']}`;
            d.setUTCSeconds(card['Date Listed']/1000);
            listed = `First Listed: ${d.toString()}`;
            desc3 = `Garage Size: ${card['Garage Size']} | ` + 
                    `Neighborhood: ${card['Neighborhood']}`;
        }

        return (
            <div style={listingStyle}>
        
                <br/><br/><br/>
                <Card style={cardStyle}>

                    <p style={{ whiteSpace: 'pre-wrap' }}>
                        <Carousel style={carouselStyle}>
                            {
                                photos && 
                                photos.map((photo) => (
                                    <Carousel.Item style={itemStyle}>
                                        <img
                                        className="d-block w-100"
                                        src={photo}
                                        alt={""}
                                        />
                                    </Carousel.Item>
                                ))
                            }  
                        </Carousel>
                        {'\u00A0'}{'\u000A'}
                    </p>

                    <br/>
                    <Card>
                        <Card.Body>
                        <Card.Title>{addr}</Card.Title>
                            <Card.Text>
                            { card && card['Description']}
                            </Card.Text>
                        </Card.Body>
                        <ListGroup className="list-group-flush">
                            <ListGroupItem>{price}</ListGroupItem>
                            <ListGroupItem>{desc1}</ListGroupItem>
                            <ListGroupItem>{desc2}</ListGroupItem>
                            <ListGroupItem>{desc3}</ListGroupItem>
                            <ListGroupItem>{listed}</ListGroupItem>
                        </ListGroup>
                        <Card.Footer>
                            <small className="text-muted">Last updated: {ago}</small>
                        </Card.Footer>
                    </Card>
                    
                </Card>
                <br/><br/><br/>

            </div>
        );
  }

}
