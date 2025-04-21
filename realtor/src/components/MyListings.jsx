import React, { useEffect } from "react";
import TileDeck from "./TileDeck";

const MyListings = (props) => {
  const [cards, setCards] = React.useState([]);

  useEffect(() => {
    const fetchMyListings = async () => {
      if (props.loggedIn && props.user) {
        try {
          const response = await fetch('/listings');
          const data = await response.json();

          const myListings = data.filter(card => card.User === props.user);
          setCards(myListings);
        } catch (error) {
          console.error('Error fetching listings:', error);
        }
      }
    };

    fetchMyListings();
  }, [props.loggedIn, props.user]);

  return (
    <div>
      {!props.loggedIn && (
        <div style={{ textAlign: "center", marginTop: "50px" }}>
          <h3>Please sign in above to see your listed properties.</h3>
          <h3>Thank you.</h3>
        </div>
      )}

      {props.loggedIn && cards.length === 0 && (
        <div style={{ textAlign: "center", marginTop: "50px" }}>
          <h3>You have no properties listed with us.</h3>
        </div>
      )}

      {props.loggedIn && cards.length > 0 && (
        <div style={{ backgroundColor: 'LightGray', margin: "0px", marginTop: "50px", padding: "0px", height: "240vh" }}>
          <TileDeck cards={cards} user={props.user} />
        </div>
      )}
    </div>
  );
};

export default MyListings;
