import React, { memo } from 'react';
import { Routes, Route } from 'react-router';
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
  log('Main rendered', { loggedIn, user });
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
          element={<MyListing loggedIn={loggedIn} user={user} />}
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
          element={<MyListing loggedIn={loggedIn} user={user} />}
        />
      </Routes>
    </main>
  );
});

export default Main;