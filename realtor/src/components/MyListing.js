import React from "react";
import { useLocation, useNavigate, useParams } from "react-router";
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
    this.listDropzoneRef = React.createRef();
    this.arrayDropzoneRef = React.createRef();
    this.isMounted = true;
    this.instanceId = Math.random().toString(36).substr(2, 9);

    this.state = {
      loggedIn: props.loggedIn ?? false,
      user: props.user ?? null,
      card: null,
      loaded: false,
    };

    console.log(`constructor: Initial props [${this.instanceId}]`, {
      loggedIn: this.state.loggedIn,
      user: this.state.user,
      location: props.location ? props.location.pathname + props.location.search : "undefined",
      params: props.params,
    });
  }

  async componentDidMount() {
    console.log(`componentDidMount: Mounting MyListing [${this.instanceId}]`, {
      isClient: typeof window !== "undefined",
      loggedIn: this.state.loggedIn,
      user: this.state.user,
      location: this.props.location ? this.props.location.pathname + this.props.location.search : "undefined",
    });
    try {
      if (!this.props.location) {
        console.warn(`componentDidMount: location prop missing [${this.instanceId}]`, { props: this.props });
        if (this.isMounted) this.setState({ loaded: true });
        return;
      }

      const search = this.props.location.search;
      const params = new URLSearchParams(search);
      const listingId = params.get("id");
      console.log(`componentDidMount: Fetching listing [${this.instanceId}]`, { listingId });
      if (listingId) {
        const response = await fetch(`/listing/${listingId}`);
        console.log(`componentDidMount: Fetch response [${this.instanceId}]`, {
          listingId,
          status: response.status,
          ok: response.ok,
        });
        if (!response.ok) throw new Error(`Failed to fetch listing: ${response.status}`);
        const data = await response.json();
        console.log(`componentDidMount: Fetch data [${this.instanceId}]`, {
          listingId,
          dataLength: data.length,
          firstItem: data[0] ? { ...data[0], "Photo Array": data[0]["Photo Array"]?.length || 0 } : null,
        });
        if (data.length > 0 && this.isMounted) {
          this.setState({ card: data[0] });
        }
      }
    } catch (error) {
      console.error(`componentDidMount: Error fetching listing [${this.instanceId}]`, { error: error.message });
    } finally {
      if (this.isMounted) {
        console.log(`componentDidMount: Setting loaded [${this.instanceId}]`, { loaded: true });
        this.setState({ loaded: true });
      }
    }
  }

  componentWillUnmount() {
    console.log(`componentWillUnmount: Unmounting MyListing [${this.instanceId}]`, {
      loggedIn: this.state.loggedIn,
      user: this.state.user,
      cardExists: !!this.state.card,
      loaded: this.state.loaded,
    });
    this.isMounted = false;
    if (this.listDropzoneRef.current) {
      this.listDropzoneRef.current.disable();
      console.log(`componentWillUnmount: Disabled listDropzone [${this.instanceId}]`);
    }
    if (this.arrayDropzoneRef.current) {
      this.arrayDropzoneRef.current.disable();
      console.log(`componentWillUnmount: Disabled arrayDropzone [${this.instanceId}]`);
    }
  }

  componentDidUpdate(prevProps, prevState) {
    if (
      prevState.loggedIn !== this.state.loggedIn ||
      prevState.user !== this.state.user ||
      prevState.card !== this.state.card ||
      prevState.loaded !== this.state.loaded
    ) {
      console.log(`componentDidUpdate: State changed [${this.instanceId}]`, {
        prevState: {
          loggedIn: prevState.loggedIn,
          user: prevState.user,
          cardExists: !!prevState.card,
          loaded: prevState.loaded,
        },
        currentState: {
          loggedIn: this.state.loggedIn,
          user: this.state.user,
          cardExists: !!this.state.card,
          loaded: this.state.loaded,
        },
      });
    }
  }

  shouldComponentUpdate(nextProps, nextState) {
    const shouldUpdate =
      this.state.loggedIn !== nextState.loggedIn ||
      this.state.user !== nextState.user ||
      this.state.card !== nextState.card ||
      this.state.loaded !== nextState.loaded ||
      (this.props.location?.pathname !== nextProps.location?.pathname ||
        this.props.location?.search !== nextProps.location?.search) ||
      JSON.stringify(this.props.params) !== JSON.stringify(nextProps.params);
    console.log(`shouldComponentUpdate [${this.instanceId}]`, {
      shouldUpdate,
      stateChanged: {
        loggedIn: this.state.loggedIn !== nextState.loggedIn,
        user: this.state.user !== nextState.user,
        card: this.state.card !== nextState.card,
        loaded: this.state.loaded !== nextState.loaded,
      },
      propsChanged: {
        locationPathname: this.props.location?.pathname !== nextProps.location?.pathname,
        locationSearch: this.props.location?.search !== nextProps.location?.search,
        params: JSON.stringify(this.props.params) !== JSON.stringify(nextProps.params),
      },
    });
    return shouldUpdate;
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
      console.log(`onSubmit: Submitting listing [${this.instanceId}]`, { mls: newUuid, user: this.state.user });
      const rawResponse = await fetch("/listings/add/HowMuchDoesSecurityCost", {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(json),
      });

      console.log(`onSubmit: Submission response [${this.instanceId}]`, { status: rawResponse.status, ok: rawResponse.ok });
      if (rawResponse.ok) {
        NotificationManager.success("Success", "Success", 3000);
        console.log(`onSubmit: Submission successful, navigating [${this.instanceId}]`, { mls: newUuid });
      } else {
        throw new Error(`Failed to submit listing: ${rawResponse.status}`);
      }
    } catch (error) {
      console.error(`onSubmit: Submission error [${this.instanceId}]`, { error: error.message });
      NotificationManager.warning("Failure", "Failure", 3000);
    }
  }

  onListChange({ meta }, status) {
    if (!this.isMounted) {
      console.warn(`onListChange: Called after unmount [${this.instanceId}]`, { metaName: meta.name, status });
      return;
    }
    const sml = "https://files.mitchelletzel.com/media/";
    const path = `${sml}${this.state.user}/${meta.name}`;
    let newCard = { ...this.state.card } || {};

    if (status === "done") {
      newCard["List Photo"] = path;
    } else if (status === "removed" && newCard["List Photo"] === path) {
      newCard["List Photo"] = "";
    }
    console.log(`onListChange: Updating card [${this.instanceId}]`, {
      status,
      metaName: meta.name,
      newListPhoto: newCard["List Photo"],
    });
    this.setState({ card: newCard });
  }

  onArrayChange({ meta }, status) {
    if (!this.isMounted) {
      console.warn(`onArrayChange: Called after unmount [${this.instanceId}]`, { metaName: meta.name, status });
      return;
    }
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
    console.log(`onArrayChange: Updating photo array [${this.instanceId}]`, {
      status,
      metaName: meta.name,
      newPhotoArray: newCard["Photo Array"],
    });
    this.setState({ card: newCard });
  }

  onRemove(photo) {
    if (!this.isMounted) {
      console.warn(`onRemove: Called after unmount [${this.instanceId}]`, { removedPhoto: photo });
      return;
    }
    let newCard = { ...this.state.card } || { "Photo Array": [] };
    newCard["Photo Array"] = (newCard["Photo Array"] || []).filter((p) => p !== photo);
    console.log(`onRemove: Removing photo [${this.instanceId}]`, {
      removedPhoto: photo,
      newPhotoArray: newCard["Photo Array"],
    });
    this.setState({ card: newCard });
  }

  render() {
    console.log(`render: Rendering MyListing [${this.instanceId}]`, {
      isClient: typeof window !== "undefined",
      loggedIn: this.state.loggedIn,
      loaded: this.state.loaded,
      cardExists: !!this.state.card,
      photoArrayLength: this.state.card?.["Photo Array"]?.length || 0,
      location: this.props.location ? this.props.location.pathname + this.props.location.search : "undefined",
      params: this.props.params,
    });

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
    console.log(`render: Carousel data [${this.instanceId}]`, {
      photoCount: photos.length,
      photos: photos,
    });

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

    let carouselContent;
    try {
      carouselContent = (
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
      );
    } catch (error) {
      console.error(`render: Carousel rendering error [${this.instanceId}]`, { error: error.message });
      carouselContent = <div>Error rendering carousel</div>;
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
          <p style={{ whiteSpace: "pre-wrap" }}>{carouselContent}</p>
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
                ref={this.listDropzoneRef}
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
                ref={this.arrayDropzoneRef}
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