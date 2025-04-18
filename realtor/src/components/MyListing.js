import React from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import Card from "react-bootstrap/Card";
import Carousel from "react-bootstrap/Carousel";
import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import "react-dropzone-uploader/dist/styles.css";
import Dropzone from "react-dropzone-uploader";
import { v4 as uuid } from "uuid";
import { NotificationContainer, NotificationManager } from "react-notifications";

// HOC to inject router props
function withRouter(Component) {
  return function Wrapper(props) {
    const location = useLocation();
    const navigate = useNavigate();
    const params = useParams();
    return <Component {...props} location={location} navigate={navigate} params={params} />;
  };
}

class MyListing extends React.Component {
  constructor(props) {
    super(props);
    this.onSubmit = this.onSubmit.bind(this);
    this.onListChange = this.onListChange.bind(this);
    this.onArrayChange = this.onArrayChange.bind(this);
    this.onRemove = this.onRemove.bind(this);

    this.formRef = React.createRef();

    this.state = {
      loggedIn: props.loggedIn ?? false,
      user: props.user ?? null,
      card: null,
      loaded: false,
    };

    // Debug initial props
    console.log("Initial props:", props);
  }

  async componentDidMount() {
    console.log("componentDidMount: Starting fetch, loggedIn:", this.state.loggedIn);
    let isMounted = true; // Track mount status
    try {
      if (!this.props.location) {
        console.warn("location prop not found; check router setup.");
        if (isMounted) this.setState({ loaded: true });
        return;
      }

      const search = this.props.location.search;
      const params = new URLSearchParams(search);
      const listingId = params.get("id");
      if (listingId) {
        const response = await fetch(`/listing/${listingId}`);
        if (!response.ok) throw new Error(`Failed to fetch listing: ${response.status}`);
        const data = await response.json();
        if (data.length > 0 && isMounted) {
          this.setState({ card: data[0] });
        }
      }
    } catch (error) {
      console.error("Error in componentDidMount:", error);
    } finally {
      if (isMounted) {
        console.log("componentDidMount: Setting loaded to true");
        this.setState({ loaded: true });
      }
    }
  }

  componentWillUnmount() {
    this.isMounted = false;
  }

  async onSubmit(event) {
    event.preventDefault();
    const { card } = this.state;
    const elements = this.formRef.current.elements;
    const time = new Date().getTime();
    const firstTime = card?.["Date Listed"] ?? time;
    const newUuid = card?.["MLS"] ?? uuid();
    const status = card?.["deleted"] ?? "false";

    const json = {
      Bathrooms: elements.Bathrooms.value,
      Bedrooms: elements.Bedrooms.value,
      City: elements.City.value,
      "Date Listed": `${firstTime}`,
      deleted: status,
      Description: elements.Description.value,
      "Garage Size": elements.GarageSize.value,
      "Last Modified": `${time}`,
      "List Photo": card?.["List Photo"] ?? "",
      "Lot Size": elements.LotSize.value,
      MLS: newUuid,
      Neighborhood: elements.Neighborhood.value,
      "Photo Array": card?.["Photo Array"] ?? [],
      "Sales Price": elements.Price.value,
      "Square Feet": elements.SquareFeet.value,
      State: elements.State.value,
      Street1: elements.Address.value,
      Street2: elements.Address2.value === "" ? "*" : elements.Address2.value,
      User: this.state.user,
      "Zip Code": elements.ZipCode.value,
    };

    try {
      const rawResponse = await fetch("/listings/add/HowMuchDoesSecurityCost", {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(json),
      });

      if (rawResponse.ok) {
        NotificationManager.success("Success", "Success", 3000);
      } else {
        throw new Error(`Failed to submit listing: ${rawResponse.status}`);
      }
    } catch (error) {
      console.error("Submission error:", error);
      NotificationManager.warning("Failure", "Failure", 3000);
    }
  }

  onListChange({ meta }, status) {
    const sml = "https://files.mitchelletzel.com/media/";
    const path = `${sml}${this.state.user}/${meta.name}`;
    let newCard = { ...this.state.card } || {};

    if (status === "done") {
      newCard["List Photo"] = path;
    } else if (status === "removed" && newCard["List Photo"] === path) {
      newCard["List Photo"] = "";
    }
    this.setState({ card: newCard });
  }

  onArrayChange({ meta }, status) {
    const sml = "https://files.mitchelletzel.com/media/";
    const path = `${sml}${this.state.user}/${meta.name}`;
    let newCard = { ...this.state.card } || { "Photo Array": [] };
    let photoArr = [...(newCard["Photo Array"] || [])];

    if (status === "done") {
      photoArr.push(path);
    } else if (status === "removed") {
      photoArr = photoArr.filter((p) => p !== path);
    }
    newCard["Photo Array"] = photoArr;
    this.setState({ card: newCard });
  }

  onRemove(photo) {
    let newCard = { ...this.state.card } || { "Photo Array": [] };
    newCard["Photo Array"] = (newCard["Photo Array"] || []).filter((p) => p !== photo);
    this.setState({ card: newCard });
  }

  render() {
    const h3Style = { textAlign: "center" };
    const listingStyle = {
      backgroundColor: "Gray",
      margin: "0px",
      padding: "0px",
      height: "400vh",
    };
    const cardStyle = {
      width: "90vw",
      paddingTop: "3vw",
      paddingLeft: "2vw",
      paddingRight: "2vw",
      paddingBottom: "3vw",
      margin: "auto",
      backgroundColor: "LightGray",
    };
    const card2Style = {
      width: "82vw",
      paddingTop: "3vw",
      paddingLeft: "2vw",
      paddingRight: "2vw",
      paddingBottom: "3vw",
      margin: "auto",
      backgroundColor: "White",
    };
    const carouselStyle = {
      width: "70vw",
      height: "25vw",
      margin: "auto",
      paddingBottom: "7vw",
      borderStyle: "solid",
      borderWidth: "8px",
    };
    const itemStyle = {
      backgroundSize: "auto",
      objectFit: "cover",
      width: "100%",
      height: "24vw",
      overflow: "hidden",
      alignItems: "center",
    };
    const buttonStyle = { margin: "0", position: "absolute", left: "50%", transform: "translateX(-50%)" };

    const photos = Array.isArray(this.state.card?.["Photo Array"]) ? this.state.card["Photo Array"] : [];

    console.log("Render: loggedIn:", this.state.loggedIn, "loaded:", this.state.loaded);

    if (!this.state.loaded) {
      return <div>Loading...</div>;
    }

    if (!this.state.loggedIn) {
      return (
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
      );
    }

    return (
      <div style={listingStyle}>
        <br />
        <br />
        <br />
        <Card style={cardStyle}>
          <h3 style={h3Style}>
            {!this.state.card ? "List your property with us." : "Edit your listing"}
          </h3>
          <br />
          <br />
          <p style={{ whiteSpace: "pre-wrap" }}>
            <Carousel style={carouselStyle}>
              {photos.map((photo) => (
                <Carousel.Item style={itemStyle} key={photo}>
                  <img className="d-block w-100" src={photo} alt="Property" />
                  <Carousel.Caption>
                    <Button variant="primary" onClick={() => this.onRemove(photo)}>
                      Remove
                    </Button>
                  </Carousel.Caption>
                </Carousel.Item>
              ))}
            </Carousel>
          </p>
          <br />
          <Card style={card2Style}>
            <Form ref={this.formRef} onSubmit={this.onSubmit}>
              <Form.Group controlId="formGridAddress1">
                <Form.Label>Address</Form.Label>
                <Form.Control
                  type="text"
                  name="Address"
                  required
                  defaultValue={this.state.card?.["Street1"] ?? ""}
                />
              </Form.Group>

              <Form.Group controlId="formGridAddress2">
                <Form.Label>Address 2</Form.Label>
                <Form.Control
                  type="text"
                  name="Address2"
                  defaultValue={
                    this.state.card?.["Street2"] === "*" ? "" : this.state.card?.["Street2"] ?? ""
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
                    defaultValue={this.state.card?.["City"] ?? ""}
                  />
                </Form.Group>

                <Form.Group as={Col} controlId="formGridState">
                  <Form.Label>State</Form.Label>
                  <Form.Control
                    type="text"
                    name="State"
                    required
                    defaultValue={this.state.card?.["State"] ?? ""}
                  />
                </Form.Group>

                <Form.Group as={Col} controlId="formGridZipCode">
                  <Form.Label>Zip Code</Form.Label>
                  <Form.Control
                    type="text"
                    name="ZipCode"
                    required
                    defaultValue={this.state.card?.["Zip Code"] ?? ""}
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
                    defaultValue={this.state.card?.["Sales Price"] ?? ""}
                  />
                </Form.Group>

                <Form.Group as={Col} controlId="formGridNeighborhood">
                  <Form.Label>Neighborhood</Form.Label>
                  <Form.Control
                    type="text"
                    name="Neighborhood"
                    required
                    defaultValue={this.state.card?.["Neighborhood"] ?? ""}
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
                    defaultValue={this.state.card?.["Bedrooms"] ?? ""}
                  />
                </Form.Group>

                <Form.Group as={Col} controlId="formGridBathrooms">
                  <Form.Label>Bathrooms</Form.Label>
                  <Form.Control
                    type="text"
                    name="Bathrooms"
                    required
                    defaultValue={this.state.card?.["Bathrooms"] ?? ""}
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
                    defaultValue={this.state.card?.["Square Feet"] ?? ""}
                  />
                </Form.Group>

                <Form.Group as={Col} controlId="formGridLotSize">
                  <Form.Label>Lot Size</Form.Label>
                  <Form.Control
                    type="text"
                    name="LotSize"
                    required
                    defaultValue={this.state.card?.["Lot Size"] ?? ""}
                  />
                </Form.Group>

                <Form.Group as={Col} controlId="formGridGarageSize">
                  <Form.Label>Garage Size</Form.Label>
                  <Form.Control
                    type="text"
                    name="GarageSize"
                    required
                    defaultValue={this.state.card?.["Garage Size"] ?? ""}
                  />
                </Form.Group>
              </Row>

              <Form.Group controlId="formGridDescription">
                <Form.Label>Description</Form.Label>
                <Form.Control
                  as="textarea"
                  rows={3}
                  name="Description"
                  required
                  defaultValue={this.state.card?.["Description"] ?? ""}
                />
              </Form.Group>

              <div>List Photo (Only One Image Please)</div>
              <br />
              <Dropzone
                getUploadParams={() => ({
                  url: `/upload/image/${this.state.user}`,
                })}
                onChangeStatus={this.onListChange}
                accept="image/*"
                maxFiles={1}
              />

              <br />
              <div>Photo Array</div>
              <br />
              <Dropzone
                getUploadParams={() => ({
                  url: `/upload/image/${this.state.user}`,
                })}
                onChangeStatus={this.onArrayChange}
                accept="image/*"
              />

              <br />
              <br />
              <Button style={buttonStyle} variant="primary" type="submit">
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
    );
  }
}

export default withRouter(MyListing);