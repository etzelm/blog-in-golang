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
  vi.spyOn(console, 'error').mockImplementation(() => {}); // Spy on console.error
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

  it('should return all listings when Submit is clicked with no criteria', async () => {
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

    // 1. Wait for initial data to load and tiles to render
    let initialTiles;
    await waitFor(() => {
      initialTiles = screen.getAllByTestId(/tile-\d+/);
      expect(initialTiles.length).toBe(listings.filter(l => l.deleted === 'false').length);
    });

    // 2. Simulate clicking the Submit button without changing any inputs
    const submitButton = screen.getByRole('button', { name: /submit/i });
    fireEvent.click(submitButton);

    // 3. Wait for filtering (which should do nothing) and assert all original tiles are still present
    await waitFor(() => {
      const currentTiles = screen.getAllByTestId(/tile-\d+/);
      // Expect the same number of tiles as initially loaded
      expect(currentTiles.length).toBe(initialTiles.length);
      // Check if a specific tile is still there
      expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
       // Ensure the "No results" message is NOT shown
      expect(screen.queryByText('No listings match your criteria.')).not.toBeInTheDocument();
    });

    // Check console log for filter values (all should be null or empty)
    expect(console.log).toHaveBeenCalledWith('Filter values:', {
      City: null, State: null, 'Zip Code': null, Bedrooms: null, Bathrooms: null, MLS: null, 'Square Feet': null
    });
  });

  it('should handle error during initial listing fetch', async () => {
    const errorMessage = 'Failed to fetch';
    // Mock fetch to simulate an error
    fetchMock.mockRejectedValue(new Error(errorMessage));

    render(
      <MemoryRouter>
        <Search loggedIn={false} user={null} />
      </MemoryRouter>
    );

    // 1. Wait for the fetch attempt and the error handling
    await waitFor(() => {
      // Check if console.error was called with the expected message
      expect(console.error).toHaveBeenCalledWith("Error fetching listings:", expect.any(Error));
      // Optionally, check the specific error message if needed:
      // expect(console.error).toHaveBeenCalledWith("Error fetching listings:", new Error(errorMessage));
    });

    // 2. Assert that no listing tiles are rendered
    const tiles = screen.queryAllByTestId(/tile-\d+/);
    expect(tiles.length).toBe(0);

    // 3. Assert that the search form is still rendered
    expect(screen.getByLabelText('City')).toBeInTheDocument();
  });

  // --- NEW TEST ADDED BELOW ---
  it('should correctly filter out listings with missing data for filtered fields', async () => {
    // Define mock listings with missing data
    const mockListingsWithMissingData = [
      // Valid listing matching filter
      { MLS: 'valid-1', Street1: '1 Valid St', City: 'Bend', State: 'OR', Bedrooms: '3', deleted: 'false', 'List Photo': 'https://placehold.co/300x200/eee/aaa?text=Valid' },
      // Listing with missing City
      { MLS: 'missing-city', Street1: '2 Missing City St', City: null, State: 'OR', Bedrooms: '3', deleted: 'false', 'List Photo': 'https://placehold.co/300x200/eee/aaa?text=Missing+City' },
      // Listing with missing Bedrooms
      { MLS: 'missing-beds', Street1: '3 Missing Beds St', City: 'Bend', State: 'OR', Bedrooms: null, deleted: 'false', 'List Photo': 'https://placehold.co/300x200/eee/aaa?text=Missing+Beds' },
      // Listing that doesn't match filter criteria
      { MLS: 'other-data', Street1: '4 Other St', City: 'Redmond', State: 'OR', Bedrooms: '2', deleted: 'false', 'List Photo': 'https://placehold.co/300x200/eee/aaa?text=Other' },
    ];

    // Mock fetch to return these specific listings
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockListingsWithMissingData),
    });

    render(
      <MemoryRouter>
        <Search loggedIn={false} user={null} />
      </MemoryRouter>
    );

    // 1. Wait for initial data to load
    await waitFor(() => {
      expect(screen.getByTestId('tile-valid-1')).toBeInTheDocument();
      expect(screen.getByTestId('tile-missing-city')).toBeInTheDocument();
      expect(screen.getByTestId('tile-missing-beds')).toBeInTheDocument();
      expect(screen.getByTestId('tile-other-data')).toBeInTheDocument();
      expect(screen.getAllByTestId(/tile-\w+/).length).toBe(5); // All 4 initially
    });

    // 2. Find inputs and submit button
    const cityInput = screen.getByLabelText('City');
    const bedroomsInput = screen.getByLabelText('Bedrooms');
    const submitButton = screen.getByRole('button', { name: /submit/i });

    // 3. Simulate entering filter criteria ('Bend' and '3' bedrooms)
    fireEvent.change(cityInput, { target: { value: 'Bend' } });
    fireEvent.change(bedroomsInput, { target: { value: '3' } });

    // 4. Simulate clicking the Submit button
    fireEvent.click(submitButton);

    // 5. Wait for filtering and assert results
    await waitFor(() => {
      // Only the valid listing should remain
      expect(screen.getByTestId('tile-valid-1')).toBeInTheDocument();
      // Listings with missing data for filtered fields should be gone
      expect(screen.queryByTestId('tile-missing-city')).not.toBeInTheDocument();
      expect(screen.queryByTestId('tile-missing-beds')).not.toBeInTheDocument();
      // Listing that didn't match criteria should be gone
      expect(screen.queryByTestId('tile-other-data')).not.toBeInTheDocument();
      // Only 1 tile should remain
      expect(screen.getAllByTestId(/tile-\w+/).length).toBe(2);
    });

     // Check console log for filter values
     expect(console.log).toHaveBeenCalledWith('Filter values:', expect.objectContaining({ City: 'bend', Bedrooms: 3 }));
     // Check console log for filtered cards result
     expect(console.log).toHaveBeenCalledWith('Filtered cards:', expect.arrayContaining([expect.objectContaining({ MLS: 'valid-1' })]));
  });
  // --- END OF NEW TEST ---

  it('should show "No listings match your criteria" when no results are found', async () => {
    const { listings } = await import('../../../test-data');
    // Mock fetch to return listings
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(listings),
    });

    render(
      <MemoryRouter>
        <Search loggedIn={false} user={null} />
      </MemoryRouter>
    );

    // Wait for initial data to load
    await waitFor(() => {
      expect(screen.getByTestId('tile-1234567890')).toBeInTheDocument();
    });

    // Find inputs and submit button
    const cityInput = screen.getByLabelText('City');
    const submitButton = screen.getByRole('button', { name: /submit/i });

    // Enter criteria that won't match any listings
    fireEvent.change(cityInput, { target: { value: 'NonExistentCity' } });

    // Submit the form
    fireEvent.click(submitButton);

    // Wait for filtering and assert no results message appears
    await waitFor(() => {
      expect(screen.getByText('No listings match your criteria.')).toBeInTheDocument();
      // Verify no tiles are present
      expect(screen.queryAllByTestId(/tile-\d+/).length).toBe(0);
    });

    // Verify state.noResults was set to true by checking the UI
    expect(screen.getByText('No listings match your criteria.')).toBeVisible();
  });

  it('should log filter mismatches for debugging', async () => {
    const mockListings = [
      {
        MLS: 'test-1',
        Street1: '1 Test St',
        City: 'Portland', // This will NOT match our filter
        State: 'OR',
        Bedrooms: '2', // This will NOT match our filter
        deleted: 'false',
        'List Photo': 'https://example.com/photo.jpg'
      }
    ];

    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockListings),
    });

    render(
      <MemoryRouter>
        <Search loggedIn={false} user={null} />
      </MemoryRouter>
    );

    // Wait for initial data to load
    await waitFor(() => {
      expect(screen.getByTestId('tile-test-1')).toBeInTheDocument();
    });

    // Find inputs and submit button
    const cityInput = screen.getByLabelText('City');
    const bedroomsInput = screen.getByLabelText('Bedrooms');
    const submitButton = screen.getByRole('button', { name: /submit/i });

    // Enter criteria that will cause mismatches (for debugging logs)
    fireEvent.change(cityInput, { target: { value: 'Seattle' } }); // Won't match Portland
    fireEvent.change(bedroomsInput, { target: { value: '3' } }); // Won't match 2

    // Submit the form
    fireEvent.click(submitButton);

    // Wait for filtering and verify debug logs were called
    await waitFor(() => {
      // Check that the mismatch logging occurred (covers lines 83-85)
      // The City filter fails first, so that's what we should see in the logs
      expect(console.log).toHaveBeenCalledWith(
        "Card 0 failed filter: City (card: portland, filter: seattle)"
      );
    });

    // Verify no results message appears
    expect(screen.getByText('No listings match your criteria.')).toBeInTheDocument();
  });
});
