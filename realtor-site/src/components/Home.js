import React from "react";
import TileDeck from "./TileDeck";

export default class Home extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            cards: []
        };
    }

    async componentDidMount() {

        const response = await fetch('/listings');
        const data = await response.json();
        console.log(data)

        var listings = [];
        for (var it=0; it<data.length; it++) {
            const card = data[it];
            if (card['deleted'] === "false") {
                listings.push(card);
            }
        }

        this.setState({ cards: listings })

    }

    render() {
        const homeStyle = {
            backgroundColor: 'LightGray',
            margin: "0px",
            padding: "0px",
            height: "240vh"
        }

        return (
            <div style={homeStyle}>
                <br/><br/>
                <TileDeck cards={this.state.cards}/>
            </div>
        );
    }

}
