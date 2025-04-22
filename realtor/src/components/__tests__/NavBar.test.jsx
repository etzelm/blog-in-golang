import React from 'react';
import { render, screen, cleanup } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import NavBar from '../NavBar';

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

describe('NavBar', () => {
  it('renders sign in button when not logged in', () => {
    render(
      <MemoryRouter>
        <NavBar loggedIn={false} user={null} onSignIn={() => {}} onSignOut={() => {}} authLoading={false} />
      </MemoryRouter>
    );
    expect(screen.getByText('Sign In')).toBeInTheDocument();
  });
});