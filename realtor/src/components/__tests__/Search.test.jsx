import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen, waitFor, cleanup, fireEvent } from '@testing-library/react';
import { BrowserRouter, MemoryRouter } from 'react-router';
import Home from '../Home';
import Search from '../Search';
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

describe('Search.jsx', () => {
  it('renders search form when not logged in', async () => {
    render(
      <MemoryRouter>
        <Search loggedIn={false} user={null} />
      </MemoryRouter>
    );

    // Check for form labels to confirm form rendering
    expect(screen.getByLabelText('City')).toBeInTheDocument();
    expect(screen.getByLabelText('State')).toBeInTheDocument();
    expect(screen.getByLabelText('Zip Code')).toBeInTheDocument();
  });

  it('should clear inputs and reset results when Reset button is clicked', async () => {
    const { listings } = await import('../../../test-data');
    // Mock fetch to return initial data
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(listings),
    });

    render(
      <MemoryRouter>
        <Search loggedIn={false} user={null} />
      </MemoryRouter>
    );

    // Wait for initial data to load and tiles to render
    await waitFor(() => {
      expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
    });

    // Find the City input and the Reset button
    const cityInput = screen.getByLabelText('City');
    const resetButton = screen.getByRole('button', { name: /reset/i });

    // Simulate user typing in the City input
    fireEvent.change(cityInput, { target: { value: 'TestCity' } });
    expect(cityInput).toHaveValue('TestCity');

    // Simulate clicking the Reset button
    fireEvent.click(resetButton);

    // Assert that the City input is cleared
    expect(cityInput).toHaveValue('');

    // Assert that the original tiles are still (or again) present
    // (This assumes reset restores the initially fetched list)
    await waitFor(() => {
      expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
      // Optionally check for other original listings if needed
      expect(screen.getAllByTestId(/tile-\d+/).length).toBe(listings.filter(l => l.deleted === 'false').length);
    });

     // Check that console log for reset was called (covers line 221)
     // Note: This requires the spyOn for console.log in beforeEach
     expect(console.log).toHaveBeenCalledWith('Form reset, restored original cards');
  });

  it('should filter listings when Submit button is clicked with criteria', async () => {
    const { listings } = await import('../../../test-data');
    // Mock fetch to return initial data (includes listings in Bend and potentially others)
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(listings),
    });

    render(
      <MemoryRouter>
        <Search loggedIn={false} user={null} />
      </MemoryRouter>
    );

    // Wait for initial data to load and tiles to render
    await waitFor(() => {
      expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument(); // Assuming this is in Bend
    });

    // Find the City input and the Submit button
    const cityInput = screen.getByLabelText('City');
    const submitButton = screen.getByRole('button', { name: /submit/i });

    // Simulate user typing 'Bend' in the City input
    fireEvent.change(cityInput, { target: { value: 'Bend' } });
    expect(cityInput).toHaveValue('Bend');

    // Simulate clicking the Submit button
    fireEvent.click(submitButton);

    // Wait for the filtering logic to apply and component to re-render
    await waitFor(() => {
      // Assert that only listings in Bend are now visible
      // Check for a known Bend listing
      expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
      // Check that listings NOT in Bend (if any in test data) are GONE
      // (Add specific checks here if your test data includes non-Bend listings)

      // Count the tiles to ensure it matches the number of Bend listings
      const expectedBendListings = listings.filter(
        l => l.City.toLowerCase() === 'bend' && l.deleted === 'false'
      ).length;
      expect(screen.getAllByTestId(/tile-\d+/).length).toBe(expectedBendListings);

      // Ensure the "No results" message is NOT shown
      expect(screen.queryByText('No listings match your criteria.')).not.toBeInTheDocument();
    });

    // Check console log for filter values (covers line 200)
    expect(console.log).toHaveBeenCalledWith('Filter values:', expect.objectContaining({ City: 'bend' }));
     // Check console log for filtered cards result (covers line 209)
     expect(console.log).toHaveBeenCalledWith('Filtered cards:', expect.any(Array));
  });
});