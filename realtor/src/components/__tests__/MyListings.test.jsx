import React from 'react';
import { render, screen, cleanup, waitFor } from '@testing-library/react';
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

it('renders user listings when logged in with fetched data', async () => {
  const userEmail = 'test@example.com';
  const mockListings = [
    {
      MLS: '1234567890',
      Street1: '123 Real St',
      City: 'Bend',
      State: 'OR',
      ZipCode: '97701',
      User: userEmail,
      deleted: 'false',
    },
  ];
  fetchMock.mockImplementation(() =>
    Promise.resolve({
      json: () => Promise.resolve(mockListings),
      ok: true,
    })
  );

  render(
    <MemoryRouter>
      <MyListings loggedIn={true} user={userEmail} />
    </MemoryRouter>
  );

  await waitFor(() => {
    expect(fetchMock).toHaveBeenCalledWith('/listings');
    expect(screen.getByTestId('tile-deck')).toBeInTheDocument();
    expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
  });
});