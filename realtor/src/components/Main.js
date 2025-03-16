import React from "react";
import { Routes, Route } from 'react-router-dom';
import Home from "./Home";
import Listing from "./Listing";
import Search from "./Search";
import MyListing from "./MyListing";
import MyListings from "./MyListings";

export default class Main extends React.Component {
  render() {
    return (
      <main>
        <Routes>
          <Route 
            path='/realtor' 
            element={
              <Home 
                loggedIn={this.props.loggedIn} 
                user={this.props.user}
              />
            }
          />
          <Route 
            path='/realtor/search' 
            element={
              <Search 
                loggedIn={this.props.loggedIn} 
                user={this.props.user}
              />
            }
          />
          <Route 
            path='/realtor/new' 
            element={
              <MyListing 
                loggedIn={this.props.loggedIn} 
                user={this.props.user}
              />
            }
          />
          <Route 
            path='/realtor/listing' 
            element={
              <Listing 
                loggedIn={this.props.loggedIn} 
                user={this.props.user}
              />
            }
          />
          <Route 
            path='/realtor/my-listings' 
            element={
              <MyListings 
                loggedIn={this.props.loggedIn} 
                user={this.props.user}
              />
            }
          />
          <Route 
            path='/realtor/my-listing' 
            element={
              <MyListing 
                loggedIn={this.props.loggedIn} 
                user={this.props.user}
              />
            }
          />
        </Routes>
      </main>
    );
  }
}
