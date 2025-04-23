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

it('renders NavDropdown when logged in', () => {
  render(
    <MemoryRouter>
      <NavBar loggedIn={true} user="test@example.com" onSignIn={() => {}} onSignOut={() => {}} authLoading={false} />
    </MemoryRouter>
  );
  expect(screen.getByText('test@example.com')).toBeInTheDocument();
  expect(screen.getByText('List Your Property')).toBeInTheDocument();
});

it('renders disabled Sign In button when authLoading is true', () => {
  render(
    <MemoryRouter>
      <NavBar loggedIn={false} user={null} onSignIn={() => {}} onSignOut={() => {}} authLoading={true} />
    </MemoryRouter>
  );
  const signInButton = screen.getByText('Signing In...');
  expect(signInButton).toBeInTheDocument();
  expect(signInButton).toBeDisabled();
});