import { render, screen, cleanup } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { vi } from 'vitest';
import NavBar from '../NavBar';

// Mock @react-oauth/google
vi.mock('@react-oauth/google', () => ({
  GoogleLogin: ({ onSuccess, onError }) => (
    <button 
      data-testid="google-login-button"
      onClick={() => {
        const mockCredential = {
          credential: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.test'
        };
        onSuccess(mockCredential);
      }}
    >
      Sign in with Google
    </button>
  ),
}));

beforeEach(() => {
  // Mock console to suppress logs
  vi.spyOn(console, 'log').mockImplementation(() => {});
  vi.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('NavBar', () => {
  it('renders Google Login button when not logged in', () => {
    render(
      <MemoryRouter>
        <NavBar 
          loggedIn={false} 
          user={null} 
          onLoginSuccess={() => {}} 
          onLoginError={() => {}} 
          onSignOut={() => {}} 
          authLoading={false} 
        />
      </MemoryRouter>
    );
    expect(screen.getByTestId('google-login-button')).toBeInTheDocument();
    expect(screen.getByText('Sign in with Google')).toBeInTheDocument();
  });

  it('renders user dropdown when logged in', () => {
    render(
      <MemoryRouter>
        <NavBar 
          loggedIn={true} 
          user="test@example.com" 
          onLoginSuccess={() => {}} 
          onLoginError={() => {}} 
          onSignOut={() => {}} 
          authLoading={false} 
        />
      </MemoryRouter>
    );
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
    expect(screen.getByText('List Your Property')).toBeInTheDocument();
    // Don't test for dropdown items as they're not visible until clicked
  });

  it('shows user email in dropdown when logged in', () => {
    render(
      <MemoryRouter>
        <NavBar 
          loggedIn={true} 
          user="test@example.com" 
          onLoginSuccess={() => {}} 
          onLoginError={() => {}} 
          onSignOut={() => {}} 
          authLoading={false} 
        />
      </MemoryRouter>
    );
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('renders navigation links', () => {
    render(
      <MemoryRouter>
        <NavBar 
          loggedIn={false} 
          user={null} 
          onLoginSuccess={() => {}} 
          onLoginError={() => {}} 
          onSignOut={() => {}} 
          authLoading={false} 
        />
      </MemoryRouter>
    );
    expect(screen.getByText('Search Listings')).toBeInTheDocument();
    expect(screen.getByText('List Your Property')).toBeInTheDocument();
  });
});