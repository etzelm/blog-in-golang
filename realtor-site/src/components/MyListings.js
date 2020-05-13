import React from "react";
import TileDeck from "./TileDeck";

export default class Home extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            loggedIn: this.props == null ? null : this.props.loggedIn,
            user: this.props == null ? null : this.props.user
        };
    }

    async componentDidMount() {
        
        window.gapi.load('auth2', () => {
            this.auth2 = window.gapi.auth2.init({
              client_id: 'ThisIsSupposedToBeAnId',
            })
      
            this.auth2.then(async () => {
                console.log('on init');
        
                const check = this.auth2.isSignedIn.get();
                const token = localStorage.getItem('aToken');
                var loggedIn = token == null ? check : true, 
                    email = token == null ? null : token;
                console.log(token);
                if (!(token == null)) {
                    loggedIn = true
                    email = token
                }

                if (loggedIn && !email) {
                    email = this.auth2.currentUser.get().getBasicProfile().getEmail();
                    console.log(email);
                }
                
                var myListings = [];
                if (email) {

                    const response = await fetch('/listings');
                    const data = await response.json();
                    console.log(data)

                    for (var it=0; it<data.length; it++) {
                        const card = data[it];
                        if (card['User'] === email) {
                            myListings.push(card);
                        }
                    }

                }

                this.setState({
                    loggedIn: loggedIn,
                    user: email,
                    loaded: true,
                    cards: myListings
                });
            });
        });

    }

    render() {
        const h3Style = {
            textAlign: "center"
        };

        const homeStyle = {
            backgroundColor: 'LightGray',
            margin: "0px",
            padding: "0px",
            height: "240vh"
        }

        return (

            <div>
                {

                    !this.state.loggedIn && 
                    <div>
                        <br/><br/><br/><br/><br/>
                        <h3 style={h3Style}>Please sign in above to see your listed properties.</h3>
                        <br/><br/>
                        <h3 style={h3Style}>Thank you.</h3>
                    </div>

                }

                {

                    this.state.loggedIn && this.state.cards.length === 0 &&
                    <div>
                        <br/><br/><br/><br/><br/>
                        <h3 style={h3Style}>You have no properties listed with us.</h3>
                    </div>

                }

                {

                    this.state.loggedIn &&  this.state.cards.length > 0 &&
                    <div style={homeStyle}>
                        <br/><br/>
                        <TileDeck cards={this.state.cards} user={this.state.user} />
                    </div>

                }


            </div>

        );
    }

}
