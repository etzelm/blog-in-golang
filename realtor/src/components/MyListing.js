import React, { useRef, useEffect, useState, useCallback } from 'react';
import { useLocation, useParams } from 'react-router';
import Card from 'react-bootstrap/Card';
import Carousel from 'react-bootstrap/Carousel';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import 'react-dropzone-uploader/dist/styles.css';
import Dropzone from 'react-dropzone-uploader';
import { v4 as uuid } from 'uuid';
import { NotificationContainer, NotificationManager } from 'react-notifications';

const MyListing = ({ loggedIn, user }) => {
  const location = useLocation();
  const params = useParams();
  const instanceId = Math.random().toString(36).substr(2, 9);
  const formRef = useRef(null);
  const listDropzoneRef = useRef(null);
  const arrayDropzoneRef = useRef(null);

  const [state, setState] = useState(() => {
    const search = location?.search || '';
    const urlParams = new URLSearchParams(search);
    const isCreateMode = !urlParams.get('id');
    return {
      loggedIn: !!loggedIn,
      user: user || null,
      card: null,
      loaded: isCreateMode,
    };
  });

  // Memoized fetch function
  const fetchListing = useCallback(async (listingId, signal) => {
    try {
      console.log(`fetchListing: Fetching listing [${instanceId}]`, { listingId });
      const response = await fetch(`/listing/${listingId}`, { signal });
      console.log(`fetchListing: Fetch response [${instanceId}]`, {
        listingId,
        status: response.status,
        ok: response.ok,
      });
      if (!response.ok) throw new Error(`Failed to fetch listing: ${response.status}`);
      const data = await response.json();
      console.log(`fetchListing: Fetch data [${instanceId}]`, {
        listingId,
        dataLength: data.length,
        firstItem: data[0] ? { ...data[0], 'Photo Array': data[0]['Photo Array']?.length || 0 } : null,
      });
      if (data.length > 0) {
        return data[0];
      }
      return null;
    } catch (error) {
      if (error.name === 'AbortError') return null;
      console.error(`fetchListing: Error fetching listing [${instanceId}]`, { error: error.message });
      return null;
    }
  }, [instanceId]);

  // Effect for fetching listing data
  useEffect(() => {
    const search = location?.search || '';
    const urlParams = new URLSearchParams(search);
    const listingId = urlParams.get('id');

    if (!listingId) {
      console.log(`useEffect: Create mode, no fetch needed [${instanceId}]`);
      setState((prev) => ({ ...prev, loaded: true }));
      return;
    }

    let isMounted = true;
    const abortController = new AbortController();

    fetchListing(listingId, abortController.signal).then((listing) => {
      if (isMounted && listing) {
        setState((prev) => ({ ...prev, card: listing, loaded: true }));
      }
    });

    return () => {
      isMounted = false;
      abortController.abort();
    };
  }, [location, fetchListing, instanceId]);

  // Effect for logging mount/unmount
  useEffect(() => {
    console.log(`useEffect: Mounting MyListing [${instanceId}]`, {
      isClient: typeof window !== 'undefined',
      loggedIn,
      user,
      location: location ? location.pathname + location.search : 'undefined',
    });

    return () => {
      console.log(`useEffect cleanup: Unmounting MyListing [${instanceId}]`, {
        loggedIn,
        user,
        cardExists: !!state.card,
        loaded: state.loaded,
      });
    };
  }, [loggedIn, user, location, instanceId, state.card, state.loaded]);

  const onSubmit = async (event) => {
    event.preventDefault();
    const { card } = state;
    const elements = formRef.current?.elements;
    if (!elements) {
      console.error(`onSubmit: Form elements missing [${instanceId}]`);
      NotificationManager.warning('Form error', 'Please try again', 3000);
      return;
    }

    const sanitizeInput = (value) => (typeof value === 'string' ? value.trim() : '');
    const time = new Date().getTime();
    const firstTime = card?.['Date Listed'] || time;
    const newUuid = card?.['MLS'] || uuid();
    const status = card?.['deleted'] || 'false';

    const json = {
      Bathrooms: sanitizeInput(elements.Bathrooms.value),
      Bedrooms: sanitizeInput(elements.Bedrooms.value),
      City: sanitizeInput(elements.City.value),
      'Date Listed': `${firstTime}`,
      deleted: status,
      Description: sanitizeInput(elements.Description.value),
      'Garage Size': sanitizeInput(elements.GarageSize.value),
      'Last Modified': `${time}`,
      'List Photo': card?.['List Photo'] || '',
      'Lot Size': sanitizeInput(elements.LotSize.value),
      MLS: newUuid,
      Neighborhood: sanitizeInput(elements.Neighborhood.value),
      'Photo Array': Array.isArray(card?.['Photo Array']) ? card['Photo Array'] : [],
      'Sales Price': sanitizeInput(elements.Price.value),
      'Square Feet': sanitizeInput(elements.SquareFeet.value),
      State: sanitizeInput(elements.State.value),
      Street1: sanitizeInput(elements.Address.value),
      Street2: sanitizeInput(elements.Address2.value) || '*',
      User: state.user || '',
      'Zip Code': sanitizeInput(elements.ZipCode.value),
    };

    const requiredFields = [
      'Bathrooms',
      'Bedrooms',
      'City',
      'Description',
      'Garage Size',
      'Lot Size',
      'Neighborhood',
      'Sales Price',
      'Square Feet',
      'State',
      'Street1',
      'Zip Code',
    ];
    for (const field of requiredFields) {
      if (!json[field]) {
        console.error(`onSubmit: Missing required field [${instanceId}]`, { field });
        NotificationManager.warning('Missing required field', `Please fill in ${field}`, 3000);
        return;
      }
    }

    try {
      console.log(`onSubmit: Submitting listing [${instanceId}]`, { mls: newUuid, user: state.user });
      const rawResponse = await fetch('/listings/add/HowMuchDoesSecurityCost', {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(json),
      });
      const responseData = await rawResponse.json();
      console.log(`onSubmit: Response data [${instanceId}]`, { responseData, status: rawResponse.status, ok: rawResponse.ok });
      if (rawResponse.ok) {
        NotificationManager.success('Success', 'Listing submitted', 3000);
        console.log(`onSubmit: Submission successful [${instanceId}]`, { mls: newUuid });
      } else {
        throw new Error(`Failed to submit listing: ${rawResponse.status}`);
      }
    } catch (error) {
      console.error(`onSubmit: Submission error [${instanceId}]`, { error: error.message });
      NotificationManager.warning('Submission failed', 'Please try again', 3000);
    }
  };

  const onListChange = useCallback(({ meta }, status) => {
    const sml = 'https://files.mitchelletzel.com/media/';
    const path = `${sml}${state.user}/${meta.name}`;
    setState((prev) => {
      const newCard = { ...prev.card } || {};
      if (status === 'done') {
        newCard['List Photo'] = path;
      } else if (status === 'removed' && newCard['List Photo'] === path) {
        newCard['List Photo'] = '';
      }
      console.log(`onListChange: Updating card [${instanceId}]`, {
        status,
        metaName: meta.name,
        newListPhoto: newCard['List Photo'],
      });
      return { ...prev, card: newCard };
    });
  }, [state.user, instanceId]);

  const onArrayChange = useCallback(({ meta }, status) => {
    const sml = 'https://files.mitchelletzel.com/media/';
    const path = `${sml}${state.user}/${meta.name}`;
    setState((prev) => {
      const newCard = { ...prev.card } || { 'Photo Array': [] };
      let photoArr = [...(newCard['Photo Array'] || [])];
      if (status === 'done') {
        photoArr.push(path);
      } else if (status === 'removed') {
        photoArr = photoArr.filter((p) => p !== path);
      }
      newCard['Photo Array'] = photoArr;
      console.log(`onArrayChange: Updating photo array [${instanceId}]`, {
        status,
        metaName: meta.name,
        newPhotoArray: newCard['Photo Array'],
      });
      return { ...prev, card: newCard };
    });
  }, [state.user, instanceId]);

  const onRemove = useCallback((photo) => {
    setState((prev) => {
      const newCard = { ...prev.card } || { 'Photo Array': [] };
      newCard['Photo Array'] = (newCard['Photo Array'] || []).filter((p) => p !== photo);
      console.log(`onRemove: Removing photo [${instanceId}]`, {
        removedPhoto: photo,
        newPhotoArray: newCard['Photo Array'],
      });
      return { ...prev, card: newCard };
    });
  }, [instanceId]);

  console.log(`render: Rendering MyListing [${instanceId}]`, {
    isClient: typeof window !== 'undefined',
    loggedIn: state.loggedIn,
    loaded: state.loaded,
    cardExists: !!state.card,
    photoArrayLength: state.card?.['Photo Array']?.length || 0,
    location: location ? location.pathname + location.search : 'undefined',
    params,
  });

  const h3Style = { textAlign: 'center' };
  const listingStyle = {
    backgroundColor: 'Gray',
    margin: '0px',
    padding: '0px',
    height: '400vh',
  };
  const cardStyle = {
    width: '90vw',
    paddingTop: '3vw',
    paddingLeft: '2vw',
    paddingRight: '2vw',
    paddingBottom: '3vw',
    margin: 'auto',
    backgroundColor: 'LightGray',
  };
  const card2Style = {
    width: '82vw',
    paddingTop: '3vw',
    paddingLeft: '2vw',
    paddingRight: '2vw',
    paddingBottom: '3vw',
    margin: 'auto',
    backgroundColor: 'White',
  };
  const carouselStyle = {
    width: '70vw',
    height: '25vw',
    margin: 'auto',
    paddingBottom: '7vw',
    borderStyle: 'solid',
    borderWidth: '8px',
  };
  const itemStyle = {
    backgroundSize: 'auto',
    objectFit: 'cover',
    width: '100%',
    height: '24vw',
    overflow: 'hidden',
    alignItems: 'center',
  };
  const buttonStyle = { margin: '0', position: 'absolute', left: '50%', transform: 'translateX(-50%)' };

  const photos = Array.isArray(state.card?.['Photo Array']) ? state.card['Photo Array'] : [];
  console.log(`render: Carousel data [${instanceId}]`, {
    photoCount: photos.length,
    photos,
  });

  if (!state.loaded) {
    return <div>Loading...</div>;
  }

  if (!state.loggedIn) {
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

  const isClient = typeof window !== 'undefined';
  let carouselContent = <div>Carousel placeholder</div>;
  if (isClient) {
    carouselContent = (
      <Carousel style={carouselStyle}>
        {photos.map((photo) => (
          <Carousel.Item style={itemStyle} key={photo}>
            <img className="d-block w-100" src={photo} alt="Property" />
            <Carousel.Caption>
              <Button variant="primary" onClick={() => onRemove(photo)}>
                Remove
              </Button>
            </Carousel.Caption>
          </Carousel.Item>
        ))}
      </Carousel>
    );
  }

  return (
    <div style={listingStyle}>
      <br />
      <br />
      <br />
      <Card style={cardStyle}>
        <h3 style={h3Style}>{!state.card ? 'List your property with us.' : 'Edit your listing'}</h3>
        <br />
        <br />
        <p style={{ whiteSpace: 'pre-wrap' }}>{carouselContent}</p>
        <br />
        <Card style={card2Style}>
          <Form ref={formRef} onSubmit={onSubmit}>
            <Form.Group controlId="formGridAddress1">
              <Form.Label>Address</Form.Label>
              <Form.Control type="text" name="Address" required defaultValue={state.card?.Street1 || ''} />
            </Form.Group>

            <Form.Group controlId="formGridAddress2">
              <Form.Label>Address 2</Form.Label>
              <Form.Control
                type="text"
                name="Address2"
                defaultValue={state.card?.Street2 === '*' ? '' : state.card?.Street2 || ''}
              />
            </Form.Group>

            <Row>
              <Form.Group as={Col} controlId="formGridCity">
                <Form.Label>City</Form.Label>
                <Form.Control type="text" name="City" required defaultValue={state.card?.City || ''} />
              </Form.Group>

              <Form.Group as={Col} controlId="formGridState">
                <Form.Label>State</Form.Label>
                <Form.Control type="text" name="State" required defaultValue={state.card?.State || ''} />
              </Form.Group>

              <Form.Group as={Col} controlId="formGridZipCode">
                <Form.Label>Zip Code</Form.Label>
                <Form.Control type="text" name="ZipCode" required defaultValue={state.card?.['Zip Code'] || ''} />
              </Form.Group>
            </Row>

            <Row>
              <Form.Group as={Col} controlId="formGridPrice">
                <Form.Label>Sales Price</Form.Label>
                <Form.Control type="text" name="Price" required defaultValue={state.card?.['Sales Price'] || ''} />
              </Form.Group>

              <Form.Group as={Col} controlId="formGridNeighborhood">
                <Form.Label>Neighborhood</Form.Label>
                <Form.Control
                  type="text"
                  name="Neighborhood"
                  required
                  defaultValue={state.card?.Neighborhood || ''}
                />
              </Form.Group>
            </Row>

            <Row>
              <Form.Group as={Col} controlId="formGridBedrooms">
                <Form.Label>Bedrooms</Form.Label>
                <Form.Control type="text" name="Bedrooms" required defaultValue={state.card?.Bedrooms || ''} />
              </Form.Group>

              <Form.Group as={Col} controlId="formGridBathrooms">
                <Form.Label>Bathrooms</Form.Label>
                <Form.Control type="text" name="Bathrooms" required defaultValue={state.card?.Bathrooms || ''} />
              </Form.Group>
            </Row>

            <Row>
              <Form.Group as={Col} controlId="formGridSquareFeet">
                <Form.Label>Square Feet</Form.Label>
                <Form.Control
                  type="text"
                  name="SquareFeet"
                  required
                  defaultValue={state.card?.['Square Feet'] || ''}
                />
              </Form.Group>

              <Form.Group as={Col} controlId="formGridLotSize">
                <Form.Label>Lot Size</Form.Label>
                <Form.Control type="text" name="LotSize" required defaultValue={state.card?.['Lot Size'] || ''} />
              </Form.Group>

              <Form.Group as={Col} controlId="formGridGarageSize">
                <Form.Label>Garage Size</Form.Label>
                <Form.Control
                  type="text"
                  name="GarageSize"
                  required
                  defaultValue={state.card?.['Garage Size'] || ''}
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
                defaultValue={state.card?.Description || ''}
              />
            </Form.Group>

            {isClient && (
              <>
                <div>List Photo (Only One Image Please)</div>
                <br />
                <Dropzone
                  ref={listDropzoneRef}
                  getUploadParams={() => ({
                    url: `/upload/image/${state.user || ''}`,
                  })}
                  onChangeStatus={onListChange}
                  accept="image/*"
                  maxFiles={1}
                />

                <br />
                <div>Photo Array</div>
                <br />
                <Dropzone
                  ref={arrayDropzoneRef}
                  getUploadParams={() => ({
                    url: `/upload/image/${state.user || ''}`,
                  })}
                  onChangeStatus={onArrayChange}
                  accept="image/*"
                />
              </>
            )}

            <br />
            <br />
            <Button style={buttonStyle} variant="primary" type="submit">
              Submit
            </Button>
            <br />
            <br />
            {isClient && <NotificationContainer />}
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
};

export default MyListing;