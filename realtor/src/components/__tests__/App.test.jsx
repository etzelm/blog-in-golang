import React from 'react';
import { render, screen, cleanup, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { vi } from 'vitest';
import App from '../../App';

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

describe('App.jsx', () => {
  it('renders NavBar and Main when not logged in', async () => {
    // Mock window.gapi to simulate successful Google Auth initialization
    vi.stubGlobal('gapi', {
      load: vi.fn((_, callback) => callback()), // Call callback immediately
      auth2: {
        init: vi.fn(() => ({
          isSignedIn: { get: vi.fn(() => false) },
          currentUser: {
            get: vi.fn(() => ({
              getBasicProfile: vi.fn(() => ({
                getEmail: vi.fn(() => null),
              })),
            })),
          },
        })),
        getAuthInstance: vi.fn(() => ({
          isSignedIn: { get: vi.fn(() => false) },
          currentUser: {
            get: vi.fn(() => ({
              getBasicProfile: vi.fn(() => ({
                getEmail: vi.fn(() => null),
              })),
            })),
          },
        })),
      },
    });

    render(
      <MemoryRouter initialEntries={['/realtor']}>
        <App />
      </MemoryRouter>
    );

    // Wait for the loaded state to render NavBar and Main
    await waitFor(() => {
      expect(screen.getByText('realtor webpage.')).toBeInTheDocument(); // From NavBar
      expect(screen.getByTestId('tile-deck')).toBeInTheDocument(); // From Main -> Home -> TileDeck
    });
  });
});