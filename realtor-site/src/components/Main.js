import React from "react";
import { Route } from 'react-router-dom'
import Home from "./Home";
import Listing from "./Listing";
import NewListing from "./NewListing";

export default class Main extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            loggedIn: false,
            user: null
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
                    exact path='/realtor/new' 
                    render={(props) => <NewListing {...props} />}
                />
                <Route 
                    exact path='/realtor/listing' 
                    render={(props) => <Listing {...props} />}
                />
            </main>
        );
    }

}