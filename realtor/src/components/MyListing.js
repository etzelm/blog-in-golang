import React from "react";
import Card from 'react-bootstrap/Card'
import Carousel from 'react-bootstrap/Carousel'
import Col from 'react-bootstrap/Col'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import 'react-dropzone-uploader/dist/styles.css'
import Dropzone from 'react-dropzone-uploader'
import uuid from 'react-uuid'
import {NotificationContainer, NotificationManager} from 'react-notifications';

export default class MyListing extends React.Component {

    constructor(props) {
        super(props);
        this.onSubmit = this.onSubmit.bind(this);
        this.onListChange = this.onListChange.bind(this);
        this.onArrayChange = this.onArrayChange.bind(this);
        this.onRemove = this.onRemove.bind(this);

        const token = localStorage.getItem('aToken');
        const loggedIn = token !== null;
        const user = token !== null ? token : null;
        
        this.state = {
            loggedIn: loggedIn,
            user: user,
            loaded: true,
            reload: true,
            card: null 
        };
    }

    async componentDidMount() {

        if (!(window.gapi == null)) {

            window.gapi.load('auth2', () => {
                this.auth2 = window.gapi.auth2.init({
                  client_id: 'ThisIsSupposedToBeAnId',
                })
          
                this.auth2.then(() => {
                  console.log('on init');
          
                  const loggedIn = this.auth2.isSignedIn.get();
                  var email = null;
                  if (loggedIn) {
                    email = this.auth2.currentUser.get().getBasicProfile().getEmail();
                    console.log(email);
                  }
                  
                  this.setState({
                    loggedIn: loggedIn,
                    user: email,
                    loaded: true
                  });
                });
            });

        }

        const search = this.props.location.search;
        const regex = /(?:\x3d)([^\x26]*)/i;
        const found = search.match(regex);

        console.log(search)
        console.log(this.props)
        console.log(this.props.location)
        if (!(found == null) && found.length > 0) {

            const response = await fetch('/listing/'+found[1]);
            const data = await response.json();

            if (data.length > 0 ) {
                this.setState({ card: data[0] })
            }

        }

    }

    async onSubmit(event) {
        console.log("submitted")
        console.log(event)
        console.log(event.currentTarget)
        console.log(event.currentTarget.elements.Address.value)
        console.log(event.target.elements.Address.value)
        event.preventDefault()

        const card = this.state.card;
        const elements = event.currentTarget.elements;
        const time = new Date().getTime();
        const firstTime = card == null || card['Date Listed'] == null ? time : `${card['Date Listed']}`;
        const newUuid = card == null || card['MLS'] == null ? uuid() : `${card['MLS']}`;
        const status = card == null || card['deleted'] == null ? "true" : `${card['deleted']}`;
        const list = card == null || card['List Photo'] == null ? "" : `${card['List Photo']}`;
        const array = card == null || card['Photo Array'] == null ? [] : card['Photo Array'];
        const json = {
            "Bathrooms": elements.Bathrooms.value,
            "Bedrooms": elements.Bedrooms.value,
            "City": elements.City.value,
            "Date Listed": `${firstTime}`,
            "deleted": status,
            "Description": elements.Description.value,
            "Garage Size": elements.GarageSize.value,
            "Last Modified": `${time}`,
            "List Photo": list,
            "Lot Size": elements.LotSize.value,
            "MLS": `${newUuid}`,
            "Neighborhood": elements.Neighborhood.value,
            "Photo Array": array,
            "Sales Price": elements.Price.value,
            "Square Feet": elements.SquareFeet.value,
            "State": elements.State.value,
            "Street1": elements.Address.value,
            "Street2": elements.Address2.value === "" ? "*" : elements.Address2.value,
            "User": this.state.user,
            "Zip Code": elements.ZipCode.value
        }

        const rawResponse = await fetch('/listings/add/HowMuchDoesSecurityCost', {
            method: 'POST',
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
            },
            body: JSON.stringify(json)
        });
        const content = await rawResponse.json();

        console.log(content);
        console.log(rawResponse);
        if (rawResponse.status === 200) {
            NotificationManager.success('Success', 'Success', 3000);
        } else {
            NotificationManager.warning('Failure', 'Failure', 3000);
        }

    }

    onListChange({ meta, file }, status) {
        console.log(status, meta, file) 
        const sml = "https://files.mitchelletzel.com/media/";
        const path = `${sml}${this.state.user}/${meta.name}`;
        if (this.state.card == null) {
            this.setState({
                card: {},
                user: this.state.user,
                loggedIn: this.state.loggedIn
            });
        }
        if (status === 'done') {
            var newCard = this.state.card;
            newCard['List Photo'] = path;
            this.setState({
                card: newCard,
                user: this.state.user,
                loggedIn: this.state.loggedIn
            });
        }
        if (status === 'removed' && path === this.state.card['List Photo']) {
            var removedCard = this.state.card;
            removedCard['List Photo'] = "";
            this.setState({
                card: removedCard,
                user: this.state.user,
                loggedIn: this.state.loggedIn
            });
        }
    }

    onArrayChange({ meta, file }, status) {
        console.log(status, meta, file) 
        const sml = "https://files.mitchelletzel.com/media/";
        const path = `${sml}${this.state.user}/${meta.name}`;
        if (this.state.card == null) {
            this.setState({
                card: {'Photo Array': []},
                user: this.state.user,
                loggedIn: this.state.loggedIn
            });
        }
        if (status === 'done') {
            this.setState({
                loggedIn: this.state.loggedIn,
                loggedOut: false,
                loaded: false,
                reload: true,
                user: this.state.user,
                card: this.state.card
            });
            if (this.state.card['Photo Array'] == null) {
                this.setState({
                    card: {'Photo Array': []},
                    user: this.state.user,
                    loggedIn: this.state.loggedIn
                });
            }
            var newArr = this.state.card['Photo Array'];
            newArr.push(path);
            var nCard = this.state.card;
            nCard['Photo Array'] = newArr;
            this.setState({
                loggedIn: this.state.loggedIn,
                loggedOut: false,
                loaded: true,
                reload: true,
                user: this.state.user,
                card: nCard
            });
        }
        if (status === 'removed') {
            this.setState({
                loggedIn: this.state.loggedIn,
                loggedOut: false,
                loaded: false,
                reload: true,
                user: this.state.user,
                card: this.state.card
            });
            var newCard = this.state.card;
            var photos = [];
            console.log(this.state.card['Photo Array']);
            for (var it=0; it < this.state.card['Photo Array'].length; it++) {
                console.log(this.state.card['Photo Array'][it])
                console.log(path)
                console.log(!(this.state.card['Photo Array'][it] === path))
                if (!(this.state.card['Photo Array'][it] === path)) {
                    photos.push(this.state.card['Photo Array'][it]);
                }
            }
            newCard['Photo Array'] = photos;
            console.log(photos)
            console.log(newCard)
            this.setState({
                loggedIn: this.state.loggedIn,
                loggedOut: false,
                loaded: true,
                reload: true,
                user: this.state.user,
                card: newCard
            });
        }
    }

    onRemove(photo) {

        this.setState({
            loggedIn: this.state.loggedIn,
            loggedOut: false,
            loaded: false,
            reload: true,
            user: this.state.user,
            card: this.state.card
        });
        var newCard = this.state.card;
        var photos = [];
        console.log(this.state.card['Photo Array']);
        for (var it=0; it < this.state.card['Photo Array'].length; it++) {
            console.log(this.state.card['Photo Array'][it])
            console.log(photo)
            console.log(!(this.state.card['Photo Array'][it] === photo))
            if (!(this.state.card['Photo Array'][it] === photo)) {
                photos.push(this.state.card['Photo Array'][it]);
            }
        }
        newCard['Photo Array'] = photos;
        console.log(photos)
        console.log(newCard)
        this.setState({
            loggedIn: this.state.loggedIn,
            loggedOut: false,
            loaded: true,
            reload: true,
            user: this.state.user,
            card: newCard
        });
        
    }

    render() {
        const h3Style = {
            textAlign: "center"
        };

        const listingStyle = {
            backgroundColor: 'Gray',
            margin: "0px",
            padding: "0px",
            height: "400vh"
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

        const card2Style = {
            width: '82vw',
            paddingTop: '3vw',
            paddingLeft: '2vw',
            paddingRight: '2vw',
            paddingBottom: '3vw',
            margin: 'auto',
            backgroundColor: 'White'
        };

        const carouselStyle = {
            width: '70vw',
            height: '25vw',
            margin: 'auto',
            paddingBottom: '7vw',
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

        const buttonStyle = {
            margin: "0",
            position: "absolute",
            left: "50%"
        };

        var photos = this.state.card == null ? 
                        [] : this.state.card['Photo Array'] == null ?
                            [] : this.state.card['Photo Array'];
        return (

            <div>
        
                {

                    !this.state.loggedIn && 
                    <div>
                        <br/><br/><br/><br/><br/>
                        <h3 style={h3Style}>Please sign in above to list your property.</h3>
                        <br/><br/>
                        <h3 style={h3Style}>Thank you.</h3>
                    </div>

                }

                {

                    this.state.loggedIn && 
                    <div style={listingStyle}>
        
                        <br/><br/><br/>
                        <Card style={cardStyle}>
                            
                            <h3 style={h3Style}>
                                {
                                    this.state.card == null && 
                                    <div>{"List your property with us."}</div>
                                }
                                {
                                    !(this.state.card == null) && 
                                    <div>{"Edit your listing"}</div>
                                }
                            </h3>
                            <br/><br/>

                            <p style={{ whiteSpace: 'pre-wrap' }}>
                                <Carousel style={carouselStyle}>
                                    {
                                        photos && this.state.loaded && 
                                        photos.map((photo) => (
                                            <Carousel.Item style={itemStyle}>
                                                <img
                                                className="d-block w-100"
                                                src={photo}
                                                alt={''}
                                                />
                                                <Carousel.Caption>
                                                    <Button
                                                        variant="primary" 
                                                        onClick={() => this.onRemove(photo)}
                                                    >
                                                        Remove
                                                    </Button>
                                                </Carousel.Caption>
                                            </Carousel.Item>
                                        ))
                                    }  
                                </Carousel>
                                {'\u00A0'}{'\u000A'}
                            </p>

                            <br/>
                            <Card style={card2Style}>
                                <Form onSubmit={this.onSubmit}>

                                    <Form.Group controlId="formGridAddress1">
                                        <Form.Label>Address</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="Address" 
                                            ref="Address" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Street1']
                                            }
                                        />
                                    </Form.Group>

                                    <Form.Group controlId="formGridAddress2">
                                        <Form.Label>Address 2</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="Address2" 
                                            ref="Address2" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Street2'] === "*" ?
                                                        "" : this.state.card['Street2']
                                            }
                                        />
                                    </Form.Group>

                                    <Form.Row>
                                        <Form.Group as={Col} controlId="formGridCity">
                                        <Form.Label>City</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="City" 
                                            ref="City" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['City']
                                            }
                                        />
                                        </Form.Group>

                                        <Form.Group as={Col} controlId="formGridState">
                                        <Form.Label>State</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="State" 
                                            ref="State" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['State']
                                            }
                                        />
                                        </Form.Group>

                                        <Form.Group as={Col} controlId="formGridZipCode">
                                        <Form.Label>Zip Code</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="ZipCode" 
                                            ref="ZipCode" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Zip Code']
                                            }
                                        />
                                        </Form.Group>
                                    </Form.Row>

                                    <Form.Row>
                                        <Form.Group as={Col} controlId="formGridPrice">
                                        <Form.Label>Sales Price</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="Price" 
                                            ref="Price" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Sales Price']
                                            }
                                        />
                                        </Form.Group>

                                        <Form.Group as={Col} controlId="formGridNeighborhood">
                                        <Form.Label>Neighborhood</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="Neighborhood" 
                                            ref="Neighborhood" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Neighborhood']
                                            }
                                        />
                                        </Form.Group>
                                    </Form.Row>

                                    <Form.Row>
                                        <Form.Group as={Col} controlId="formGridBedrooms">
                                        <Form.Label>Bedrooms</Form.Label>
                                        <Form.Control 
                                            type="text"
                                            name="Bedrooms" 
                                            ref="Bedrooms" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Bedrooms']
                                            }
                                        />
                                        </Form.Group>

                                        <Form.Group as={Col} controlId="formGridBathrooms">
                                        <Form.Label>Bathrooms</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="Bathrooms" 
                                            ref="Bathrooms" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Bedrooms']
                                            }
                                        />
                                        </Form.Group>
                                    </Form.Row>

                                    <Form.Row>
                                        <Form.Group as={Col} controlId="formGridSquareFeet">
                                        <Form.Label>Square Feet</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="SquareFeet" 
                                            ref="SquareFeet" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Square Feet']
                                            }
                                        />
                                        </Form.Group>

                                        <Form.Group as={Col} controlId="formGridLotSize">
                                        <Form.Label>Lot Size</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="LotSize" 
                                            ref="LotSize" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Lot Size']
                                            }
                                        />
                                        </Form.Group>

                                        <Form.Group as={Col} controlId="formGridGarageSize">
                                        <Form.Label>Garage Size</Form.Label>
                                        <Form.Control 
                                            type="text" 
                                            name="GarageSize" 
                                            ref="GarageSize" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Garage Size']
                                            }
                                        />
                                        </Form.Group>
                                    </Form.Row>
                                    
                                    <Form.Group controlId="formGridDescription">
                                        <Form.Label>Description</Form.Label>
                                        <Form.Control 
                                            as="textarea" 
                                            rows="3" 
                                            type="text" 
                                            name="Description" 
                                            ref="Description" 
                                            required="true" 
                                            defaultValue={
                                                this.state.card == null ? 
                                                    null : 
                                                    this.state.card['Description']
                                            }
                                        />
                                    </Form.Group>

                                    List Photo(Only One Image Please)<br/><br/>
                                    <Dropzone
                                        getUploadParams={() => ({ url: `/upload/image/${this.state.user}` })}
                                        onChangeStatus={this.onListChange}
                                        accept="image/*"
                                    />

                                    <br/>Photo Array<br/><br/>
                                    <Dropzone
                                        getUploadParams={() => ({ url: `/upload/image/${this.state.user}` })}
                                        onChangeStatus={this.onArrayChange}
                                        accept="image/*"
                                    />
                                    
                                    <br/><br/>
                                    <Button 
                                        style={buttonStyle} 
                                        variant="primary" 
                                        type="submit"
                                    >
                                        Submit
                                    </Button>
                                    <br/><br/>
                                    <NotificationContainer/>
                                    <br/><br/>
                                </Form>
                            </Card>
                            
                        </Card>
                        <br/><br/><br/>

                    </div>

                }

            </div>

        );
    }

}
