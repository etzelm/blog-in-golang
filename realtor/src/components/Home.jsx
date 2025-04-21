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
    try {
      const response = await fetch('/listings');
      if (!response.ok) {
        throw new Error(`Failed to fetch listings: ${response.status}`);
      }
      const data = await response.json();
      console.log('Raw data:', data);

      // Filter non-deleted listings and deduplicate by MLS
      const listings = data.filter(card => card && card.deleted === "false");
      const uniqueListings = Array.from(
        new Map(listings.map(item => [item.MLS, item])).values()
      );
      console.log('Filtered listings:', uniqueListings.map(l => ({ MLS: l.MLS, deleted: l.deleted })));
      this.setState({ cards: uniqueListings });
    } catch (error) {
      console.error('Error fetching listings:', error);
      this.setState({ cards: [] });
    }
  }

  render() {
    console.log('Rendering with cards:', this.state.cards.map(c => c.MLS));
    const homeStyle = {
      backgroundColor: 'LightGray',
      margin: "0px",
      padding: "0px",
      height: "240vh"
    };

    return (
      <div style={homeStyle}>
        <br />
        <br />
        <TileDeck cards={this.state.cards} />
      </div>
    );
  }
}