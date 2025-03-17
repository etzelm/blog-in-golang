import React from 'react';
import { Row, Col } from 'react-bootstrap';
import Tile from './Tile';

const TileDeck = ({ cards, user }) => {
  const sortedCards = [...cards].sort(
    (a, b) => new Date(b['Last Modified']) - new Date(a['Last Modified'])
  );

  return (
<div style={{ padding: '2vw 2vw' }}>  
  <Row className="g-4" xs={1} md={2} lg={3} xxl={4}>  
    {sortedCards.map((card) => (  
      <Col key={card.MLS} className="mb-4">  
        <Tile card={card} user={user} />  
      </Col>  
    ))}  
  </Row>  
</div>

);
}; 
export default TileDeck;