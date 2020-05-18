import React from "react";
import CardColumns from 'react-bootstrap/CardColumns'
import Tile from "./Tile";

export default class TileDeck extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            cards: []
        };
    }

    render() {
        const columnStyle = {
            paddingLeft: '2vw', 
            paddingRight: '2vw'
        }

        return (
            <div>

                <br/>
                <CardColumns style={columnStyle} className="card-columns">
                    
                    {

                        this.props.cards
                                    .sort((a, b) => a['Last Modified'] < b['Last Modified'] ? 1 : -1)
                                    .map((card) => (<Tile key={card.MLS} card={card} user={this.props.user}/>))

                    }

                </CardColumns>

            </div>
        );
  }

}
