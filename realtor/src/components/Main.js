import React from "react";
import { Route } from 'react-router-dom'
import Home from "./Home";
import Listing from "./Listing";
import Search from "./Search";
import MyListing from "./MyListing";
import MyListings from "./MyListings";

export default class Main extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            loggedIn: this.props == null ? null : this.props.loggedIn,
            user: this.props == null ? null : this.props.user
        };
    }

    render() {
        return (
            <main>
                <Route 
                    exact path='/realtor' 
                    render={(props) => <Home {...props} />}
                />
                <Route 
                    exact path='/realtor/search' 
                    render={(props) => <Search {...props} />}
                />
                <Route 
                    exact path='/realtor/new' 
                    render={(props) => <MyListing {...props} />}
                />
                <Route 
                    exact path='/realtor/listing' 
                    render={(props) => <Listing {...props} />}
                />
                <Route 
                    exact path='/realtor/my-listings' 
                    render={(props) => <MyListings {...props} />}
                />
                <Route 
                    exact path='/realtor/my-listing' 
                    render={(props) => <MyListing {...props} />}
                />
            </main>
        );
    }

}