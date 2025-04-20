import React from 'react';
import { Row, Col } from 'react-bootstrap';
import Tile from './Tile';

const TileDeck = ({ cards, user }) => {
  // Add basic input validation
  if (!Array.isArray(cards)) {
    console.warn('cards prop should be an array');
    return null;
  }

  const sortedCards = [...cards].sort((a, b) => {
    // Add date validation
    const dateA = new Date(a['Last Modified']);
    const dateB = new Date(b['Last Modified']);
    
    // Handle invalid dates
    if (isNaN(dateA.getTime())) return 1;
    if (isNaN(dateB.getTime())) return -1;
    
    return dateB - dateA;
  });

  return (
    <div data-testid="tile-deck" style={{ padding: '2vw 2vw' }}>
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