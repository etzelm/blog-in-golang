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
        this.setState({ cards: data })

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
                <TileDeck cards={this.state.cards}/>
            </div>
        );
    }

}
