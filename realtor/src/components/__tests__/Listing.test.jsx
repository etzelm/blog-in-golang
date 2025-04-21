import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { render, screen, cleanup } from '@testing-library/react';
import { BrowserRouter } from 'react-router';
import Listing from '../Listing';
import '@testing-library/jest-dom';

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

describe('Listing Component', () => {
  it('renders Listing component with listing prop without crashing', async () => {
    const { listings } = await import('../../../test-data');
    const listing = listings[0]; // listing1.json, MLS: 1234567890
    const { container } = render(<BrowserRouter><Listing listing={listing} /></BrowserRouter>);
    expect(fetchMock).not.toHaveBeenCalled();
    const card = container.querySelector('.card');
    expect(card).toBeInTheDocument();
    expect(screen.getByText(/last updated/i)).toBeInTheDocument();
  });
});