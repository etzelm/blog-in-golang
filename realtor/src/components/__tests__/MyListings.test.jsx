import React from 'react';
import { render, screen, cleanup } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import MyListings from '../MyListings';

// Mock fetch (not used, but included for consistency)
let fetchMock;
beforeEach(() => {
  fetchMock = vi.fn(() =>
    Promise.resolve({
      json: () => Promise.resolve([]),
      ok: true,
    })
  );
  global.fetch = fetchMock;
  fetchMock.mockClear();
  // Mock console to suppress logs
  vi.spyOn(console, 'log').mockImplementation(() => {});
  vi.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  cleanup();
  vi.restoreAllMocks(); // Restore console mocks
});

describe('MyListings', () => {
  it('renders sign in message when not logged in', () => {
    render(
      <MemoryRouter>
        <MyListings loggedIn={false} user={null} />
      </MemoryRouter>
    );
    expect(screen.getByText('Please sign in above to see your listed properties.')).toBeInTheDocument();
  });
});