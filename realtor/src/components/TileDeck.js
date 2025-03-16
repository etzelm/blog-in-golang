import React from "react";
import { Row, Col } from 'react-bootstrap';
import Tile from "./Tile";

export default class TileDeck extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
          loggedIn: this.props == null ? null : this.props.loggedIn,
          user: this.props == null ? null : this.props.user,
          cards: []
        };
    }

    render() {
      const sortedCards = [...this.state.cards].sort(
        (a, b) => new Date(b['Last Modified']) - new Date(a['Last Modified'])
      );
      return (
<div style={{ padding: '2vw 2vw' }}>  
  <Row className="g-4" xs={1} md={2} lg={3} xxl={4}>  
    {sortedCards.map((card) => (  
      <Col key={card.MLS} className="mb-2">  
        <Tile card={card} user={this.state.user} />  
      </Col>  
    ))}  
  </Row>  
</div>

);

}

}
