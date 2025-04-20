import React, { useRef, useEffect, useState, useCallback, useMemo } from 'react';
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

// Error Boundary Component
class MyListingErrorBoundary extends React.Component {
  state = { error: null };

  static getDerivedStateFromError(error) {
    return { error: error.message };
  }

  componentDidCatch(error, errorInfo) {
    console.error('ErrorBoundary caught error:', {
      error: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
      timestamp: new Date().toISOString(),
    });
  }

  render() {
    if (this.state.error) {
      return (
        <div>
          <h3>Error rendering MyListing: {this.state.error}</h3>
          <p>Check the console for details.</p>
        </div>
      );
    }
    return this.props.children;
  }
}

// Structured logging function
const log = (message, data = {}) => {
  console.log(JSON.stringify({
    message,
    timestamp: new Date().toISOString(),
    ...data,
  }, null, 2));
};

const MyListing = ({ loggedIn, user }) => {
  const location = useLocation();
  const params = useParams();
  const instanceId = useMemo(() => Math.random().toString(36).substr(2, 9), []);
  const formRef = useRef(null);
  const listDropzoneRef = useRef(null);
  const arrayDropzoneRef = useRef(null);
  const isMountedRef = useRef(true);

  const [state, setState] = useState(() => {
    const search = location?.search || '';
    const urlParams = new URLSearchParams(search);
    const isCreateMode = !urlParams.get('id');
    log('Initializing state', { isCreateMode, loggedIn, user, instanceId });
    return {
      loggedIn: !!loggedIn,
      user: user || null,
      card: null,
      loaded: isCreateMode,
    };
  });

  // Guarded state setter
  const safeSetState = useCallback((updater) => {
    if (isMountedRef.current) {
      setState((prev) => {
        const newState = typeof updater === 'function' ? updater(prev) : updater;
        log('State updated', { instanceId, newState });
        return newState;
      });
    } else {
      log('Blocked state update after unmount', { instanceId });
    }
  }, [instanceId]);

  // Fetch listing data
  useEffect(() => {
    const search = location?.search || '';
    const urlParams = new URLSearchParams(search);
    const listingId = urlParams.get('id');

    if (!listingId) {
      log('Create mode, no fetch needed', { instanceId });
      return;
    }

    const abortController = new AbortController();

    const fetchListing = async () => {
      log('fetchListing started', { instanceId, listingId });
      try {
        const response = await fetch(`/listing/${listingId}`, { signal: abortController.signal });
        log('fetchListing response', { instanceId, listingId, status: response.status, ok: response.ok });
        if (!response.ok) throw new Error(`Failed to fetch listing: ${response.status}`);
        const data = await response.json();
        log('fetchListing data', {
          instanceId,
          listingId,
          dataLength: data.length,
          firstItem: data[0] ? { ...data[0], 'Photo Array': data[0]['Photo Array']?.length || 0 } : null,
        });
        if (data.length > 0 && isMountedRef.current) {
          safeSetState((prev) => ({ ...prev, card: data[0], loaded: true }));
        }
      } catch (error) {
        if (error.name === 'AbortError') {
          log('fetchListing aborted', { instanceId, listingId });
          return;
        }
        log('fetchListing error', { instanceId, listingId, error: error.message });
      }
    };

    fetchListing();

    return () => {
      log('Fetch effect cleanup', { instanceId });
      abortController.abort();
    };
  }, [location, instanceId, safeSetState]);

  // Log mount/unmount and prop changes
  useEffect(() => {
    log('Component mounted', {
      instanceId,
      isClient: typeof window !== 'undefined',
      loggedIn,
      user,
      location: location.pathname + location.search,
      params,
    });

    return () => {
      isMountedRef.current = false;
      log('Component unmounted', {
        instanceId,
        loggedIn,
        user,
        cardExists: !!state.card,
        loaded: state.loaded,
      });
    };
  }, [loggedIn, user, location, params, instanceId, state.card, state.loaded]);

  // Clean up Dropzone on unmount
  useEffect(() => {
    return () => {
      if (listDropzoneRef.current) {
        log('Cleaning up list Dropzone', { instanceId });
        listDropzoneRef.current.removeAllFiles();
      }
      if (arrayDropzoneRef.current) {
        log('Cleaning up array Dropzone', { instanceId });
        arrayDropzoneRef.current.removeAllFiles();
      }
    };
  }, [instanceId]);

  const onSubmit = async (event) => {
    event.preventDefault();
    log('Form submission started', { instanceId });
    const { card } = state;
    const elements = formRef.current?.elements;
    if (!elements) {
      log('Form elements missing', { instanceId });
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
        log('Missing required field', { instanceId, field });
        NotificationManager.warning('Missing required field', `Please fill in ${field}`, 3000);
        return;
      }
    }

    try {
      log('Submitting listing', { instanceId, mls: newUuid, user: state.user });
      const rawResponse = await fetch('/listings/add/HowMuchDoesSecurityCost', {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(json),
      });
      const responseData = await rawResponse.json();
      log('Submission response', { instanceId, responseData, status: rawResponse.status, ok: rawResponse.ok });
      if (rawResponse.ok) {
        NotificationManager.success('Success', 'Listing submitted', 3000);
        log('Submission successful', { instanceId, mls: newUuid });
      } else {
        throw new Error(`Failed to submit listing: ${rawResponse.status}`);
      }
    } catch (error) {
      log('Submission error', { instanceId, error: error.message });
      NotificationManager.warning('Submission failed', 'Please try again', 3000);
    }
  };

  const onListChange = useCallback(({ meta }, status) => {
    log('onListChange triggered', { instanceId, metaName: meta.name, status });
    if (!isMountedRef.current) {
      log('Blocked onListChange after unmount', { instanceId });
      return;
    }
    const sml = 'https://realtor-site-images.s3-us-west-1.amazonaws.com/media/';
    const path = `${sml}${state.user}/${meta.name}`;
    safeSetState((prev) => {
      const newCard = prev.card ? { ...prev.card } : {};
      if (status === 'done') {
        newCard['List Photo'] = path;
      } else if (status === 'removed' && newCard['List Photo'] === path) {
        newCard['List Photo'] = '';
      }
      return { ...prev, card: newCard };
    });
  }, [state.user, instanceId, safeSetState]);

  const onArrayChange = useCallback(({ meta }, status) => {
    log('onArrayChange triggered', { instanceId, metaName: meta.name, status });
    if (!isMountedRef.current) {
      log('Blocked onArrayChange after unmount', { instanceId });
      return;
    }
    const sml = 'https://realtor-site-images.s3-us-west-1.amazonaws.com/media/';
    const path = `${sml}${state.user}/${meta.name}`;
    safeSetState((prev) => {
      const newCard = prev.card ? { ...prev.card } : { 'Photo Array': [] };
      const photoArr = Array.isArray(newCard['Photo Array']) ? [...newCard['Photo Array']] : [];
      if (status === 'done') {
        photoArr.push(path);
      } else if (status === 'removed') {
        const index = photoArr.indexOf(path);
        if (index !== -1) {
          photoArr.splice(index, 1);
        }
      }
      newCard['Photo Array'] = photoArr;
      return { ...prev, card: newCard };
    });
  }, [state.user, instanceId, safeSetState]);

  const onRemove = useCallback((photo) => {
    log('onRemove triggered', { instanceId, photo });
    if (!isMountedRef.current) {
      log('Blocked onRemove after unmount', { instanceId });
      return;
    }
    safeSetState((prev) => {
      const newCard = prev.card ? { ...prev.card } : { 'Photo Array': [] };
      const photoArr = Array.isArray(newCard['Photo Array']) ? [...newCard['Photo Array']] : [];
      const index = photoArr.indexOf(photo);
      if (index !== -1) {
        photoArr.splice(index, 1);
      }
      newCard['Photo Array'] = photoArr;
      return { ...prev, card: newCard };
    });
  }, [instanceId, safeSetState]);

  log('Rendering MyListing', {
    instanceId,
    isClient: typeof window !== 'undefined',
    loggedIn: state.loggedIn,
    loaded: state.loaded,
    cardExists: !!state.card,
    photoArrayLength: state.card?.['Photo Array']?.length || 0,
    location: location.pathname + location.search,
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
  log('Carousel data', { instanceId, photoCount: photos.length, photos });

  if (!state.loaded) {
    log('Rendering loading state', { instanceId });
    return <div>Loading...</div>;
  }

  if (!state.loggedIn) {
    log('Rendering not logged in state', { instanceId });
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

  log('Rendering main content', { instanceId });
  return (
    <MyListingErrorBoundary>
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
    </MyListingErrorBoundary>
  );
};

export default MyListing;