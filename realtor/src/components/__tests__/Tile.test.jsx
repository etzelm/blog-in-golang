// src/components/__tests__/Tile.test.jsx

import React from 'react';
import { render, screen, fireEvent, waitFor, cleanup } from '@testing-library/react';
import { BrowserRouter } from 'react-router'; // Use react-router-dom
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import Tile from '../Tile'; // Adjust path as needed
import '@testing-library/jest-dom';

// --- Mock Data ---
const mockListing = {
  MLS: 'tile-test-123',
  Street1: '123 Test St',
  Street2: '*',
  City: 'Testville',
  State: 'TS',
  'Zip Code': '12345',
  Bedrooms: '3',
  Bathrooms: '2',
  'Square Feet': '1500',
  'Lot Size': '3000',
  'Sales Price': '250,000',
  'List Photo': 'https://placehold.co/300x200/eee/ccc?text=Placeholder',
  'Last Modified': String(new Date().getTime() - 86400000 * 2), // 2 days ago
  deleted: 'false', // Start with listing not deleted
};

const mockUser = 'testuser@example.com';

// --- Mock fetch ---
let fetchMock;

beforeEach(() => {
  // Setup fetch mock before each test
  fetchMock = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ status: 'success' }), // Mock successful API response
    })
  );
  global.fetch = fetchMock;
  fetchMock.mockClear();

  // Mock console to suppress expected logs during tests
  vi.spyOn(console, 'log').mockImplementation(() => {});
  vi.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  // Clean up DOM and restore mocks after each test
  cleanup();
  vi.restoreAllMocks();
});

// --- Test Suite ---
describe('Tile Component', () => {
  // --- New Test Added Below ---

  it('renders edit/remove buttons for logged-in user and handles state change on click', async () => {
    // Render the Tile with a user prop
    render(
      <BrowserRouter>
        <Tile card={mockListing} user={mockUser} />
      </BrowserRouter>
    );

    // Verify Edit Listing button is present
    expect(screen.getByText('Edit Listing')).toBeInTheDocument();

    // Verify Remove Listing button is present (since deleted is 'false')
    const removeButton = screen.getByText('Remove Listing');
    expect(removeButton).toBeInTheDocument();
    expect(screen.queryByText('Publish Listing')).not.toBeInTheDocument();

    // Simulate clicking the Remove Listing button
    fireEvent.click(removeButton);

    // Wait for the fetch call to complete and the component to re-render
    await waitFor(() => {
      // Verify fetch was called correctly to update the listing status
      expect(fetchMock).toHaveBeenCalledTimes(1);
      expect(fetchMock).toHaveBeenCalledWith(
        '/listings/add/HowMuchDoesSecurityCost', // Endpoint from Tile.jsx
        expect.objectContaining({
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
          },
          // Check that the body contains the updated 'deleted' status
          body: expect.stringContaining(`"MLS":"${mockListing.MLS}"`),
          body: expect.stringContaining(`"deleted":"true"`), // Should now be true
        })
      );
    });

    // Wait for the button text to change after state update
    await waitFor(() => {
      // Verify the button now says "Publish Listing"
      expect(screen.getByText('Publish Listing')).toBeInTheDocument();
      // Verify the "Remove Listing" button is gone
      expect(screen.queryByText('Remove Listing')).not.toBeInTheDocument();
    });

    // Optional: Simulate clicking "Publish Listing" to toggle back
    const publishButton = screen.getByText('Publish Listing');
    fireEvent.click(publishButton);

    // Wait for the second fetch call and re-render
    await waitFor(() => {
      // Verify fetch was called again
      expect(fetchMock).toHaveBeenCalledTimes(2);
       // Check that the body contains the updated 'deleted' status (back to false)
       expect(fetchMock).toHaveBeenCalledWith(
        '/listings/add/HowMuchDoesSecurityCost',
        expect.objectContaining({
          body: expect.stringContaining(`"deleted":"false"`),
        })
      );
    });

     // Wait for the button text to change back
     await waitFor(() => {
      // Verify the button now says "Remove Listing" again
      expect(screen.getByText('Remove Listing')).toBeInTheDocument();
      expect(screen.queryByText('Publish Listing')).not.toBeInTheDocument();
    });
  });
});
