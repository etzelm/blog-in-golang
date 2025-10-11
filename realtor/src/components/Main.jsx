import React, { memo } from 'react';
import { Routes, Route, useLocation } from 'react-router';
import Home from './Home';
import Listing from './Listing';
import Search from './Search';
import MyListing from './MyListing';
import MyListings from './MyListings';

const log = (message, data = {}) => {
  console.log(JSON.stringify({
    message,
    timestamp: new Date().toISOString(),
    ...data,
  }, null, 2));
};

const Main = memo(({ loggedIn, user }) => {
  const location = useLocation();
  log('Main rendered', { loggedIn, user });
  
  // Create a stable key for MyListing to prevent unnecessary re-mounts
  const getMyListingKey = () => {
    const urlParams = new URLSearchParams(location.search);
    const mlsId = urlParams.get('MLS') || 'new';
    return `${user || 'anonymous'}-${mlsId}`;
  };
  
  return (
    <main>
      <Routes>
        <Route 
          path='/realtor' 
          element={<Home loggedIn={loggedIn} user={user} />}
        />
        <Route 
          path='/realtor/search' 
          element={<Search loggedIn={loggedIn} user={user} />}
        />
        <Route 
          path='/realtor/new' 
          element={<MyListing key={getMyListingKey()} loggedIn={loggedIn} user={user} />}
        />
        <Route 
          path='/realtor/listing' 
          element={<Listing loggedIn={loggedIn} user={user} />}
        />
        <Route 
          path='/realtor/my-listings' 
          element={<MyListings loggedIn={loggedIn} user={user} />}
        />
        <Route 
          path='/realtor/my-listing' 
          element={<MyListing key={getMyListingKey()} loggedIn={loggedIn} user={user} />}
        />
      </Routes>
    </main>
  );
});

export default Main;