import React from 'react';
import { render, screen, cleanup } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import Main from '../Main';

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

describe('Main', () => {
  it('renders Home component at /realtor route', () => {
    render(
      <MemoryRouter initialEntries={['/realtor']}>
        <Main loggedIn={false} user={null} />
      </MemoryRouter>
    );
    // Home renders a div with homeStyle, so we check for its presence
    expect(screen.getByTestId('tile-deck')).toBeInTheDocument();
  });
});