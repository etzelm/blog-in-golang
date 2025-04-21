import React from "react";
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';

export default class Tile extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      card: this.props.card || null
    };
    this.timeDifference = this.timeDifference.bind(this);
    this.onStateChange = this.onStateChange.bind(this);
  }

  timeDifference(current, previous) {
    var msPerMinute = 60 * 1000;
    var msPerHour = msPerMinute * 60;
    var msPerDay = msPerHour * 24;
    var msPerMonth = msPerDay * 30;
    var msPerYear = msPerDay * 365;

    var elapsed = current - previous;

    if (elapsed < msPerMinute) {
      return Math.round(elapsed/1000) + ' seconds ago';   
    } else if (elapsed < msPerHour) {
      return Math.round(elapsed/msPerMinute) + ' minutes ago';   
    } else if (elapsed < msPerDay ) {
      return Math.round(elapsed/msPerHour ) + ' hours ago';   
    } else if (elapsed < msPerMonth) {
      return Math.round(elapsed/msPerDay) + ' days ago';   
    } else if (elapsed < msPerYear) {
      return Math.round(elapsed/msPerMonth) + ' months ago';   
    } else {
      return Math.round(elapsed/msPerYear ) + ' years ago';   
    }
  }

  async onStateChange(event) {
    event.preventDefault();
    console.log("submitted");

    var newCard = this.props.card;
    if (this.props.card['deleted'] === "true") {
      newCard['deleted'] = "false";
    } else {
      newCard['deleted'] = "true";
    }

    const rawResponse = await fetch('/listings/add/HowMuchDoesSecurityCost', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(newCard)
    });
    const content = await rawResponse.json();

    console.log(content);

    this.setState({
      card: newCard,
    });
  }

  render() {
    const imgStyle = {
      width: '100%', 
      height: '15vw',
      whiteSpace: 'pre-line'
    };
    const linkStyle = {
      cursor: 'pointer',
      color: 'black'
    };
    const street = (this.props.card['Street2'] !== "*") ?
      `${this.props.card['Street1']}, ${this.props.card['Street2']} | ` :
      `${this.props.card['Street1']} | `;
    const addr = street + 
      `${this.props.card['City']}, ${this.props.card['State']} ` + 
      `${this.props.card['Zip Code']}`;
    const desc1 = `Beds: ${this.props.card['Bedrooms']} | ` + 
      `Baths: ${this.props.card['Bathrooms']}`; 
    const desc2 = `Square Feet: ${this.props.card['Square Feet']} | ` + 
      `Lot Size: ${this.props.card['Lot Size']}`;
    const time = new Date().getTime();
    const ago = this.timeDifference(time, this.props.card['Last Modified']);
    const price = `Price: $${this.props.card['Sales Price']}`;
    
    return (
      <a style={linkStyle} href={"/realtor/listing?MLS=" + this.props.card['MLS']} data-testid={`tile-${this.props.card['MLS']}`}>
        <Card>
          <Card.Img style={imgStyle} variant="top" src={this.props.card['List Photo']} />
          <Card.Body>
            <Card.Title>{addr}</Card.Title>
            <Card.Text>{price}</Card.Text>
            <Card.Text>{desc2}</Card.Text>
            <Card.Text>{desc1}</Card.Text>
            {this.props.user && (
              <div>
                <a style={linkStyle} href={"/realtor/my-listing?MLS=" + this.props.card['MLS']}>
                  <Button variant="primary">Edit Listing</Button>
                </a>
                {'       '}
                {this.state.card["deleted"] === "true" && (
                  <Button variant="primary" onClick={this.onStateChange}>Publish Listing</Button>
                )}
                {this.state.card["deleted"] === "false" && (
                  <Button variant="primary" onClick={this.onStateChange}>Remove Listing</Button>
                )}
              </div>
            )}
          </Card.Body>
          <Card.Footer>
            <small className="text-muted">Last updated: {ago}</small>
          </Card.Footer>
        </Card>
      </a>
    );
  }
}