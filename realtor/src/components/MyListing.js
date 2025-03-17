import React from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom"; // <-- For React Router v6
import Card from "react-bootstrap/Card";
import Carousel from "react-bootstrap/Carousel";
import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import "react-dropzone-uploader/dist/styles.css";
import Dropzone from "react-dropzone-uploader";
import uuid from "react-uuid";
import { NotificationContainer, NotificationManager } from "react-notifications";

// A simple custom HOC to inject location, navigate, params into a class component:
function withRouter(Component) {
  return function Wrapper(props) {
    const location = useLocation();
    const navigate = useNavigate();
    const params = useParams();
    return (
      <Component
        {...props}
        location={location}
        navigate={navigate}
        params={params}
      />
    );
  };
}

class MyListing extends React.Component {
  constructor(props) {
    super(props);
    this.onSubmit = this.onSubmit.bind(this);
    this.onListChange = this.onListChange.bind(this);
    this.onArrayChange = this.onArrayChange.bind(this);
    this.onRemove = this.onRemove.bind(this);

    this.state = {
      loggedIn: this.props == null ? null : this.props.loggedIn,
      user: this.props == null ? null : this.props.user,
      card: null
    };
  }

  async componentDidMount() {
    // Ensure location is available
    if (!this.props.location) {
      console.warn("location prop not found; check your router setup.");
      return;
    }

    // Google auth check
    if (window.gapi) {
      window.gapi.load("auth2", () => {
        this.auth2 = window.gapi.auth2.init({
          client_id: "ThisIsSupposedToBeAnId"
        });
        this.auth2.then(() => {
          const loggedIn = this.auth2.isSignedIn.get();
          let email = null;
          if (loggedIn) {
            email = this.auth2.currentUser.get().getBasicProfile().getEmail();
          }
          this.setState({
            loggedIn,
            user: email,
            loaded: true
          });
        });
      });
    }

    // Extract listing info from query param
    const search = this.props.location.search;
    const regex = /(?:\x3d)([^\x26]*)/i;
    const found = search.match(regex);
    if (found && found.length > 0) {
      const response = await fetch("/listing/" + found[1]);
      const data = await response.json();
      if (data.length > 0) {
        this.setState({ card: data[0] });
      }
    }
  }

  async onSubmit(event) {
    event.preventDefault();
    const { card } = this.state;
    const elements = event.currentTarget.elements;
    const time = new Date().getTime();
    const firstTime =
      !card || !card["Date Listed"] ? time : `${card["Date Listed"]}`;
    const newUuid = !card || !card["MLS"] ? uuid() : `${card["MLS"]}`;
    const status = !card || !card["deleted"] ? "true" : `${card["deleted"]}`;
    const list = !card || !card["List Photo"] ? "" : `${card["List Photo"]}`;
    const array = !card || !card["Photo Array"] ? [] : card["Photo Array"];

    const json = {
      Bathrooms: elements.Bathrooms.value,
      Bedrooms: elements.Bedrooms.value,
      City: elements.City.value,
      "Date Listed": `${firstTime}`,
      deleted: status,
      Description: elements.Description.value,
      "Garage Size": elements.GarageSize.value,
      "Last Modified": `${time}`,
      "List Photo": list,
      "Lot Size": elements.LotSize.value,
      MLS: `${newUuid}`,
      Neighborhood: elements.Neighborhood.value,
      "Photo Array": array,
      "Sales Price": elements.Price.value,
      "Square Feet": elements.SquareFeet.value,
      State: elements.State.value,
      Street1: elements.Address.value,
      Street2: elements.Address2.value === "" ? "*" : elements.Address2.value,
      User: this.state.user,
      "Zip Code": elements.ZipCode.value
    };

    const rawResponse = await fetch("/listings/add/HowMuchDoesSecurityCost", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json"
      },
      body: JSON.stringify(json)
    });

    if (rawResponse.status === 200) {
      NotificationManager.success("Success", "Success", 3000);
    } else {
      NotificationManager.warning("Failure", "Failure", 3000);
    }
  }

  onListChange({ meta }, status) {
    const sml = "https://files.mitchelletzel.com/media/";
    const path = `${sml}${this.state.user}/${meta.name}`;
    let newCard = this.state.card || {};

    if (status === "done") {
      newCard["List Photo"] = path;
    }
    if (status === "removed" && newCard["List Photo"] === path) {
      newCard["List Photo"] = "";
    }
    this.setState({ card: newCard });
  }

  onArrayChange({ meta }, status) {
    const sml = "https://files.mitchelletzel.com/media/";
    const path = `${sml}${this.state.user}/${meta.name}`;
    let newCard = this.state.card || { "Photo Array": [] };
    let photoArr = newCard["Photo Array"] || [];

    if (status === "done") {
      photoArr.push(path);
      newCard["Photo Array"] = photoArr;
      this.setState({ card: newCard });
    }
    if (status === "removed") {
      let filtered = photoArr.filter((p) => p !== path);
      newCard["Photo Array"] = filtered;
      this.setState({ card: newCard });
    }
  }

  onRemove(photo) {
    let newCard = this.state.card || { "Photo Array": [] };
    let filtered = (newCard["Photo Array"] || []).filter((p) => p !== photo);
    newCard["Photo Array"] = filtered;
    this.setState({ card: newCard });
  }

  render() {
    const h3Style = { textAlign: "center" };
    const listingStyle = {
      backgroundColor: "Gray",
      margin: "0px",
      padding: "0px",
      height: "400vh"
    };
    const cardStyle = {
      width: "90vw",
      paddingTop: "3vw",
      paddingLeft: "2vw",
      paddingRight: "2vw",
      paddingBottom: "3vw",
      margin: "auto",
      backgroundColor: "LightGray"
    };
    const card2Style = {
      width: "82vw",
      paddingTop: "3vw",
      paddingLeft: "2vw",
      paddingRight: "2vw",
      paddingBottom: "3vw",
      margin: "auto",
      backgroundColor: "White"
    };
    const carouselStyle = {
      width: "70vw",
      height: "25vw",
      margin: "auto",
      paddingBottom: "7vw",
      borderStyle: "solid",
      borderWidth: "8px"
    };
    const itemStyle = {
      backgroundSize: "auto",
      objectFit: "cover",
      width: "100%",
      height: "24vw",
      overflow: "hidden",
      alignItems: "center"
    };
    const buttonStyle = { margin: "0", position: "absolute", left: "50%" };

    let photos = [];
    if (this.state.card && this.state.card["Photo Array"]) {
      photos = this.state.card["Photo Array"];
    }

    return (
      <div>
        {!this.state.loggedIn && (
          <div>
            <br />
            <br />
            <br />
            <br />
            <br />
            <h3 style={h3Style}>Please sign in above to list your property.</h3>
            <br />
            <br />
            <h3 style={h3Style}>Thank you.</h3>
          </div>
        )}
        {this.state.loggedIn && (
          <div style={listingStyle}>
            <br />
            <br />
            <br />
            <Card style={cardStyle}>
              <h3 style={h3Style}>
                {!this.state.card && <div>List your property with us.</div>}
                {this.state.card && <div>Edit your listing</div>}
              </h3>
              <br />
              <br />
              <p style={{ whiteSpace: "pre-wrap" }}>
                <Carousel style={carouselStyle}>
                  {photos.map((photo) => (
                    <Carousel.Item style={itemStyle} key={photo}>
                      <img className="d-block w-100" src={photo} alt="" />
                      <Carousel.Caption>
                        <Button variant="primary" onClick={() => this.onRemove(photo)}>
                          Remove
                        </Button>
                      </Carousel.Caption>
                    </Carousel.Item>
                  ))}
                </Carousel>
                {"\u00A0"}
                {"\u000A"}
              </p>
              <br />
              <Card style={card2Style}>
                <Form onSubmit={this.onSubmit}>
                  <Form.Group controlId="formGridAddress1">
                    <Form.Label>Address</Form.Label>
                    <Form.Control
                      type="text"
                      name="Address"
                      required
                      defaultValue={
                        this.state.card ? this.state.card["Street1"] : ""
                      }
                    />
                  </Form.Group>

                  <Form.Group controlId="formGridAddress2">
                    <Form.Label>Address 2</Form.Label>
                    <Form.Control
                      type="text"
                      name="Address2"
                      defaultValue={
                        this.state.card
                          ? this.state.card["Street2"] === "*"
                            ? ""
                            : this.state.card["Street2"]
                          : ""
                      }
                    />
                  </Form.Group>

                  <Row>
                    <Form.Group as={Col} controlId="formGridCity">
                      <Form.Label>City</Form.Label>
                      <Form.Control
                        type="text"
                        name="City"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["City"] : ""
                        }
                      />
                    </Form.Group>

                    <Form.Group as={Col} controlId="formGridState">
                      <Form.Label>State</Form.Label>
                      <Form.Control
                        type="text"
                        name="State"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["State"] : ""
                        }
                      />
                    </Form.Group>

                    <Form.Group as={Col} controlId="formGridZipCode">
                      <Form.Label>Zip Code</Form.Label>
                      <Form.Control
                        type="text"
                        name="ZipCode"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["Zip Code"] : ""
                        }
                      />
                    </Form.Group>
                  </Row>

                  <Row>
                    <Form.Group as={Col} controlId="formGridPrice">
                      <Form.Label>Sales Price</Form.Label>
                      <Form.Control
                        type="text"
                        name="Price"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["Sales Price"] : ""
                        }
                      />
                    </Form.Group>

                    <Form.Group as={Col} controlId="formGridNeighborhood">
                      <Form.Label>Neighborhood</Form.Label>
                      <Form.Control
                        type="text"
                        name="Neighborhood"
                        required
                        defaultValue={
                          this.state.card
                            ? this.state.card["Neighborhood"]
                            : ""
                        }
                      />
                    </Form.Group>
                  </Row>

                  <Row>
                    <Form.Group as={Col} controlId="formGridBedrooms">
                      <Form.Label>Bedrooms</Form.Label>
                      <Form.Control
                        type="text"
                        name="Bedrooms"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["Bedrooms"] : ""
                        }
                      />
                    </Form.Group>

                    <Form.Group as={Col} controlId="formGridBathrooms">
                      <Form.Label>Bathrooms</Form.Label>
                      <Form.Control
                        type="text"
                        name="Bathrooms"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["Bathrooms"] : ""
                        }
                      />
                    </Form.Group>
                  </Row>

                  <Row>
                    <Form.Group as={Col} controlId="formGridSquareFeet">
                      <Form.Label>Square Feet</Form.Label>
                      <Form.Control
                        type="text"
                        name="SquareFeet"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["Square Feet"] : ""
                        }
                      />
                    </Form.Group>

                    <Form.Group as={Col} controlId="formGridLotSize">
                      <Form.Label>Lot Size</Form.Label>
                      <Form.Control
                        type="text"
                        name="LotSize"
                        required
                        defaultValue={
                          this.state.card ? this.state.card["Lot Size"] : ""
                        }
                      />
                    </Form.Group>

                    <Form.Group as={Col} controlId="formGridGarageSize">
                      <Form.Label>Garage Size</Form.Label>
                      <Form.Control
                        type="text"
                        name="GarageSize"
                        required
                        defaultValue={
                          this.state.card
                            ? this.state.card["Garage Size"]
                            : ""
                        }
                      />
                    </Form.Group>
                  </Row>

                  <Form.Group controlId="formGridDescription">
                    <Form.Label>Description</Form.Label>
                    <Form.Control
                      as="textarea"
                      rows="3"
                      name="Description"
                      required
                      defaultValue={
                        this.state.card ? this.state.card["Description"] : ""
                      }
                    />
                  </Form.Group>

                  <div>List Photo (Only One Image Please)</div>
                  <br />
                  <Dropzone
                    getUploadParams={() => ({
                      url: `/upload/image/${this.state.user}`
                    })}
                    onChangeStatus={this.onListChange}
                    accept="image/*"
                  />

                  <br />
                  Photo Array
                  <br />
                  <br />
                  <Dropzone
                    getUploadParams={() => ({
                      url: `/upload/image/${this.state.user}`
                    })}
                    onChangeStatus={this.onArrayChange}
                    accept="image/*"
                  />

                  <br />
                  <br />
                  <Button
                    style={buttonStyle}
                    variant="primary"
                    type="submit"
                  >
                    Submit
                  </Button>
                  <br />
                  <br />
                  <NotificationContainer />
                  <br />
                  <br />
                </Form>
              </Card>
            </Card>
            <br />
            <br />
            <br />
          </div>
        )}
      </div>
    );
  }
}

export default withRouter(MyListing);