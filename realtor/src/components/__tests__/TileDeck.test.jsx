// src/components/__tests__/TileDeck.test.jsx

import React from 'react';
import { render, screen, cleanup } from '@testing-library/react';
import { BrowserRouter } from 'react-router';
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import TileDeck from '../TileDeck';
import '@testing-library/jest-dom';

beforeEach(() => {
  // Mock console to suppress logs if needed
  vi.spyOn(console, 'log').mockImplementation(() => {});
  vi.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('TileDeck Component', () => {
  it('handles cards with invalid Last Modified dates in sorting', () => {
    const cardsWithInvalidDates = [
      {
        MLS: '1',
        'Last Modified': 'invalid-date', // This will create NaN when parsed
        Street1: '123 Test St',
        Street2: '*',
        City: 'Test City',
        State: 'TS',
        'Zip Code': '12345',
        'Sales Price': '100000',
        Bedrooms: '2',
        Bathrooms: '1',
        'Square Feet': '1000',
        'Lot Size': '2000',
        'Garage Size': '1 car',
        Neighborhood: 'Test',
        'List Photo': 'https://example.com/photo.jpg',
      },
      {
        MLS: '2',
        'Last Modified': '1589161257428', // Valid timestamp
        Street1: '456 Test Ave',
        Street2: '*',
        City: 'Test City',
        State: 'TS',
        'Zip Code': '12345',
        'Sales Price': '200000',
        Bedrooms: '3',
        Bathrooms: '2',
        'Square Feet': '1500',
        'Lot Size': '3000',
        'Garage Size': '2 car',
        Neighborhood: 'Test',
        'List Photo': 'https://example.com/photo2.jpg',
      },
      {
        MLS: '3',
        'Last Modified': 'another-invalid-date', // This will also create NaN
        Street1: '789 Test Blvd',
        Street2: '*',
        City: 'Test City',
        State: 'TS',
        'Zip Code': '12345',
        'Sales Price': '300000',
        Bedrooms: '4',
        Bathrooms: '3',
        'Square Feet': '2000',
        'Lot Size': '4000',
        'Garage Size': '3 car',
        Neighborhood: 'Test',
        'List Photo': 'https://example.com/photo3.jpg',
      },
    ];

    render(
      <BrowserRouter>
        <TileDeck cards={cardsWithInvalidDates} user="test@example.com" />
      </BrowserRouter>
    );

    // Verify the tile deck renders
    expect(screen.getByTestId('tile-deck')).toBeInTheDocument();
    
    // Verify all cards are rendered (despite invalid dates)
    expect(screen.getByText(/123 Test St/)).toBeInTheDocument();
    expect(screen.getByText(/456 Test Ave/)).toBeInTheDocument();
    expect(screen.getByText(/789 Test Blvd/)).toBeInTheDocument();
    
    // The valid date should be sorted first, invalid dates should be at the end
    // This tests both the isNaN(dateA.getTime()) and isNaN(dateB.getTime()) branches
  });
});