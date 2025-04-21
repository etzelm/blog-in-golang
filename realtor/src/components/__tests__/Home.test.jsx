import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen, waitFor, cleanup } from '@testing-library/react';
import { BrowserRouter } from 'react-router';
import Home from '../Home';
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
  vi.spyOn(console, 'log').mockImplementation(() => {});
  vi.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  cleanup();
  vi.restoreAllMocks(); // Restore console mocks
});

describe('Home Component', () => {
  it('renders Home component without crashing', async () => {
    render(<BrowserRouter><Home /></BrowserRouter>);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));
    expect(fetchMock).toHaveBeenCalledWith('/listings');
    const tileDecks = screen.queryAllByTestId('tile-deck');
    expect(tileDecks.length).toBeGreaterThan(0);
    const tiles = screen.queryAllByTestId(/tile-\d+/);
    expect(tiles.length).toBe(0);
  });

  it('renders listings when API returns data', async () => {
    const { listings } = await import('../../../test-data');
    fetchMock.mockImplementation(() =>
      Promise.resolve({
        json: () => Promise.resolve(listings),
        ok: true,
      })
    );

    render(<BrowserRouter><Home /></BrowserRouter>);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));
    expect(fetchMock).toHaveBeenCalledWith('/listings');
    const tileDecks = screen.queryAllByTestId('tile-deck');
    expect(tileDecks.length).toBeGreaterThan(0);
    expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
    expect(screen.getByTestId('tile-1234567891')).toBeInTheDocument();
    expect(screen.getByTestId('tile-1234567892')).toBeInTheDocument();
    expect(screen.getByTestId('tile-1234567893')).toBeInTheDocument();
    expect(screen.getByTestId('tile-1234567894')).toBeInTheDocument();
    const tiles = screen.getAllByTestId(/tile-\d+/);
    expect(tiles.length).toBe(5);
  });

  it('filters out deleted listings', async () => {
    const { listings } = await import('../../../test-data');
    const modifiedListings = [
      ...listings,
      {
        ...listings[0],
        MLS: '9999999999',
        deleted: 'true'
      }
    ];
    fetchMock.mockImplementation(() =>
      Promise.resolve({
        json: () => Promise.resolve(modifiedListings),
        ok: true,
      })
    );

    render(<BrowserRouter><Home /></BrowserRouter>);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));
    const tileDecks = screen.queryAllByTestId('tile-deck');
    expect(tileDecks.length).toBeGreaterThan(0);
    const tiles = screen.getAllByTestId(/tile-\d+/);
    expect(tiles.length).toBe(5);
    expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
    expect(screen.queryByTestId('tile-9999999999')).not.toBeInTheDocument();
  });

  it('renders correct listing details in Tile components', async () => {
    const { listings } = await import('../../../test-data');
    const singleListing = [listings[0]];
    fetchMock.mockImplementation(() =>
      Promise.resolve({
        json: () => Promise.resolve(singleListing),
        ok: true,
      })
    );

    render(<BrowserRouter><Home /></BrowserRouter>);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));
    const tileDecks = screen.queryAllByTestId('tile-deck');
    expect(tileDecks.length).toBeGreaterThan(0);
    const tile = screen.getByTestId('tile-1234567890');
    expect(tile).toHaveTextContent('123 Real Avenue, Apt. 56 | Bend, OR 97701');
    expect(tile).toHaveTextContent('Price: $503,000');
    expect(tile).toHaveTextContent('Square Feet: 1200 sqft | Lot Size: 1600 sqft');
    expect(tile).toHaveTextContent('Beds: 3 | Baths: 1');
  });

  it('handles API errors gracefully', async () => {
    fetchMock.mockImplementation(() =>
      Promise.resolve({
        ok: false,
        status: 500,
        json: () => Promise.resolve({}),
      })
    );

    render(<BrowserRouter><Home /></BrowserRouter>);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));
    const tileDecks = screen.queryAllByTestId('tile-deck');
    expect(tileDecks.length).toBeGreaterThan(0);
    const tiles = screen.queryAllByTestId(/tile-\d+/);
    expect(tiles.length).toBe(0);
  });

  it('shows empty TileDeck before API response', async () => {
    const { listings } = await import('../../../test-data');
    fetchMock.mockImplementation(() =>
      new Promise((resolve) => {
        setTimeout(() => {
          resolve({
            json: () => Promise.resolve(listings),
            ok: true,
          });
        }, 100);
      })
    );

    render(<BrowserRouter><Home /></BrowserRouter>);
    const tileDecks = screen.queryAllByTestId('tile-deck');
    expect(tileDecks.length).toBeGreaterThan(0);
    const tiles = screen.queryAllByTestId(/tile-\d+/);
    expect(tiles.length).toBe(0);

    await waitFor(
      () => expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument(),
      { timeout: 2000 }
    );
    expect(fetchMock).toHaveBeenCalledTimes(1);
    const finalTiles = screen.getAllByTestId(/tile-\d+/);
    expect(finalTiles.length).toBe(5);
  });
});