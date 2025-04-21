import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { render, screen, cleanup } from '@testing-library/react';
import { BrowserRouter } from 'react-router';
import MyListing from '../MyListing';
import '@testing-library/jest-dom';

// Mock fetch
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

describe('MyListing Component', () => {
  it('renders sign-in message without user', () => {
    render(<BrowserRouter><MyListing /></BrowserRouter>);
    expect(fetchMock).not.toHaveBeenCalled();
    expect(screen.getByText(/please sign in above to list your property/i)).toBeInTheDocument();
    expect(screen.queryAllByTestId('tile-deck')).toHaveLength(0);
    expect(screen.queryAllByTestId(/tile-\d+/)).toHaveLength(0);
  });
});