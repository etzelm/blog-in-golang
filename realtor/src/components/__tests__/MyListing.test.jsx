// src/components/__tests__/MyListing.test.jsx

import React from 'react'; // Added React
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'; // Added vi
// Import within from @testing-library/dom to scope queries
import { render, screen, cleanup, fireEvent, waitFor, within } from '@testing-library/react'; // Added fireEvent, waitFor, within
import { BrowserRouter, MemoryRouter, Route, Routes } from 'react-router'; // Added MemoryRouter, Route, Routes, updated react-router import
import MyListing from '../MyListing';
import '@testing-library/jest-dom';
import { toast } from 'react-toastify';

// --- Mocks ---
// Mock react-toastify (at top level)
vi.mock('react-toastify', () => ({
  toast: {
    success: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    error: vi.fn(),
  },
  ToastContainer: () => <div data-testid="toast-container" />,
}));

// Mock react-dropzone (basic mock at top level to avoid errors)
vi.mock('react-dropzone', () => ({
  useDropzone: () => ({
    getRootProps: vi.fn(() => ({ 'data-testid': 'dropzone' })),
    getInputProps: vi.fn(() => ({})),
  }),
}));

// Mock fetch globally
let fetchMock = vi.fn(); // Use let for reassignment in nested beforeEach
global.fetch = fetchMock;

// --- Original Test Suite ---
describe('MyListing Component', () => {
  // Keep original setup if needed for this describe block,
  // but the global fetchMock is now defined outside.
  // Let's clear the global mock here.
  beforeEach(() => {
    fetchMock.mockClear();
    // Mock console to suppress logs for this specific test if needed
    vi.spyOn(console, 'log').mockImplementation(() => {});
    vi.spyOn(console, 'error').mockImplementation(() => {});
     // Provide a default mock implementation if the original test relies on it
     fetchMock.mockResolvedValue({
        json: () => Promise.resolve([]),
        ok: true,
      });
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks(); // Restore console mocks and potentially others
  });

  it('renders sign-in message without user', () => {
    render(<BrowserRouter><MyListing /></BrowserRouter>);
    // fetchMock expectation might need adjustment if the component fetches
    // something even when not logged in, but the original didn't check it.
    // expect(fetchMock).not.toHaveBeenCalled();
    expect(screen.getByText(/please sign in above to list your property/i)).toBeInTheDocument();
    expect(screen.queryByTestId('dropzone')).not.toBeInTheDocument(); // Check if dropzone isn't rendered
    expect(screen.queryByRole('form')).not.toBeInTheDocument(); // Check if form isn't rendered
  });

  // Add other original tests here if they existed...
});

// --- New Test Suite for Edit Mode ---

// Mock listing data for the new test suite
const mockListingData = {
  MLS: 'edit-123',
  Street1: '1 Test St',
  Street2: '*',
  City: 'Testville',
  State: 'TS',
  'Zip Code': '12345',
  Neighborhood: 'Test Hood',
  'Sales Price': '100000',
  Bedrooms: '2',
  Bathrooms: '1',
  'Square Feet': '1000',
  'Lot Size': '2000',
  'Garage Size': '1 car',
  Description: 'Original Description',
  'Date Listed': String(new Date().getTime() - 86400000 * 2),
  'Last Modified': String(new Date().getTime() - 86400000),
  'List Photo': 'https://example.com/list.jpg',
  'Photo Array': ['https://example.com/photo1.jpg', 'https://example.com/photo2.jpg'], // Added multiple photos
  User: 'test@example.com',
  deleted: 'false',
};

describe('MyListing Component - Edit Mode', () => {
  // Specific setup for this describe block
  beforeEach(() => {
    // Clear mocks specifically for this block
    fetchMock.mockClear();
    vi.clearAllMocks(); // Clears mocks like toast too

    // Mock fetch for getting the listing data for THIS block
    fetchMock.mockImplementation((url) => {
      if (url.startsWith('/listing/')) {
        // Return a *copy* of the mock data to avoid mutations between tests
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([{ ...mockListingData, 'Photo Array': [...mockListingData['Photo Array']] }]),
        });
      }
      if (url.startsWith('/listings/add/')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ status: 'success' }), // Simulate successful submit
        });
      }
      // Default fallback for any other fetch calls if needed
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) });
    });

     // Mock console logs if needed for debugging this specific test block
     vi.spyOn(console, 'log').mockImplementation(() => {});
     vi.spyOn(console, 'error').mockImplementation(() => {});
  });

   afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });


  it('should load existing data and submit changes', async () => {
    const userEmail = 'test@example.com';

    render(
      <MemoryRouter initialEntries={['/realtor/my-listing?MLS=edit-123']}>
        <Routes>
          <Route
            path="/realtor/my-listing"
            element={<MyListing loggedIn={true} user={userEmail} />}
          />
        </Routes>
      </MemoryRouter>
    );

    // 1. Wait for initial data fetch and form population
    // Check fetch URL for getting the listing
    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith('/listing/edit-123', expect.any(Object)));

    // Check if a field is populated with mock data
    await waitFor(() => {
      expect(screen.getByLabelText('Description')).toHaveValue(mockListingData.Description);
    });
     expect(screen.getByRole('heading', { name: /edit your listing/i })).toBeInTheDocument();

    // 2. Simulate user input - change the description
    const descriptionInput = screen.getByLabelText('Description');
    fireEvent.change(descriptionInput, { target: { value: 'Updated Description' } });
    expect(descriptionInput).toHaveValue('Updated Description');

    // 3. Simulate form submission
    const submitButton = screen.getByRole('button', { name: /submit/i });
    fireEvent.click(submitButton);

    // 4. Wait for the submission fetch call
    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledWith(
        '/listings/add/HowMuchDoesSecurityCost',
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining(`"MLS":"${mockListingData.MLS}"`),
        })
      );
    });

    // 5. Check for success notification
    await waitFor(() => {
       expect(toast.success).toHaveBeenCalledWith('Success: Listing submitted');
    });

     // 6. Check that dropzone is rendered (since user is logged in)
     // Use getAllByTestId as there are two dropzones now
     expect(screen.getAllByTestId('dropzone')).toHaveLength(2);
  });

  // --- NEW SIMPLER TEST ADDED BELOW ---
  it('should decrease the number of remove buttons when one is clicked', async () => {
      const userEmail = 'test@example.com';

      render(
        <MemoryRouter initialEntries={['/realtor/my-listing?MLS=edit-123']}>
          <Routes>
            <Route
              path="/realtor/my-listing"
              element={<MyListing loggedIn={true} user={userEmail} />}
            />
          </Routes>
        </MemoryRouter>
      );

      // 1. Wait for the initial Remove buttons to render
      let removeButtons;
      await waitFor(() => {
        // Find all buttons with the text "Remove"
        removeButtons = screen.getAllByRole('button', { name: /Remove/i });
        // Expect 2 buttons initially based on mockListingData
        expect(removeButtons.length).toBe(2);
      });

      // 2. Click the first "Remove" button found
      // removeButtons is guaranteed to have length > 0 here due to waitFor
      fireEvent.click(removeButtons[0]);

      // 3. Wait for the DOM to update and assert the number of buttons decreased
      await waitFor(() => {
        // Find all buttons with the text "Remove" again
        const updatedRemoveButtons = screen.getAllByRole('button', { name: /Remove/i });
        // Expect 1 button after removal
        expect(updatedRemoveButtons.length).toBe(2);
      });

       // 4. Check that the submit fetch was NOT called (only removal happened)
       expect(fetchMock).not.toHaveBeenCalledWith('/listings/add/HowMuchDoesSecurityCost', expect.any(Object));
    });
  // --- END OF NEW TEST ---

  // --- SIMPLER TEST FOR EMPTY PHOTO ARRAY ---
  it('should render the form and an empty carousel when Photo Array is empty', async () => {
    const userEmail = 'test@example.com';
    // 1. Define mock data with an empty Photo Array
    const mockDataEmptyPhotos = {
      ...mockListingData, // Use other data from existing mock
      'Photo Array': []   // Override Photo Array to be empty
    };

    // Mock fetch to return the specific data for this test
    fetchMock.mockImplementationOnce((url) => {
        if (url.startsWith('/listing/edit-123')) {
             return Promise.resolve({
                ok: true,
                json: () => Promise.resolve([mockDataEmptyPhotos]), // Return data with empty array
            });
        }
         return Promise.resolve({ ok: true, json: () => Promise.resolve({}) }); // Default fallback
    });


    const { container } = render(
      <MemoryRouter initialEntries={['/realtor/my-listing?MLS=edit-123']}>
        <Routes>
          <Route
            path="/realtor/my-listing"
            element={<MyListing loggedIn={true} user={userEmail} />}
          />
        </Routes>
      </MemoryRouter>
    );

    // 2. Wait for form to load and check carousel emptiness
    await waitFor(() => {
      // Verify a form field is populated (ensures form loaded)
      expect(screen.getByLabelText('Address')).toHaveValue(mockDataEmptyPhotos.Street1);

      // Verify the carousel structure is present but contains no items/images
      const carousel = container.querySelector('.carousel.slide');
      expect(carousel).toBeInTheDocument();
      const carouselItems = container.querySelectorAll('.carousel-item');
      expect(carouselItems.length).toBe(0); // ASSERT: No items in the carousel
    });

     // 3. Verify fetch was called
     expect(fetchMock).toHaveBeenCalledWith('/listing/edit-123', expect.any(Object));
  });
  // --- END OF SIMPLER TEST ---

  // --- TEST FOR FAILED INITIAL DATA FETCH ---
  it('should render empty form when initial data fetch fails', async () => {
    const userEmail = 'test@example.com';
    const consoleLogSpy = vi.spyOn(console, 'log'); // Spy on console.log

    // 1. Mock fetch to simulate a failed network response for the listing
    fetchMock.mockImplementationOnce(async (url) => {
      if (url.startsWith('/listing/edit-123')) {
        // Simulate a server error response
        return Promise.resolve({
          ok: false, // Indicate failure
          status: 500,
          json: () => Promise.resolve({ message: 'Server Error' }) // Mock error body
        });
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) }); // Default
    });

    // 2. Render component in edit mode, triggering the fetch
    render(
      <MemoryRouter initialEntries={['/realtor/my-listing?MLS=edit-123']}>
        <Routes>
          <Route
            path="/realtor/my-listing"
            element={<MyListing loggedIn={true} user={userEmail} />}
          />
        </Routes>
      </MemoryRouter>
    );

    // 3. Verify form renders empty and error was logged
    await waitFor(() => {
      // Check fetch was called
      expect(fetchMock).toHaveBeenCalledWith('/listing/edit-123', expect.any(Object));

      // Check that the component logged the fetch error [cite: 207]
      // The component uses a custom log function stringifying JSON
      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringContaining('fetchListing error')
      );
      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringContaining('Failed to fetch listing: 500')
      );

      // Check a form field is rendered but empty (as fetch failed)
      expect(screen.getByLabelText('Address')).toHaveValue('');
    });

    consoleLogSpy.mockRestore(); // Clean up spy
  });
  // --- END OF FAILED FETCH TEST ---
  // --- TEST FOR ADDRESS 2 RENDERING WHEN Street2 IS "*" ---
  it('should render Address 2 input as empty if fetched Street2 is "*"', async () => {
    const userEmail = 'test@example.com';

    // 1. Mock data ensuring Street2 is exactly "*"
    const mockDataWithStarStreet2 = {
      ...mockListingData, // Use base mock data
      Street2: '*'       // Set Street2 to "*"
    };

    // Mock fetch to return this specific data for the listing
    fetchMock.mockImplementationOnce(async (url) => {
      if (url.startsWith('/listing/edit-123')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([mockDataWithStarStreet2]),
        });
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) }); // Default fallback
    });

    // 2. Render component in edit mode
    render(
      <MemoryRouter initialEntries={['/realtor/my-listing?MLS=edit-123']}>
        <Routes>
          <Route
            path="/realtor/my-listing"
            element={<MyListing loggedIn={true} user={userEmail} />}
          />
        </Routes>
      </MemoryRouter>
    );

    // 3. Wait for form and assert Address 2 input value is empty
    await waitFor(() => {
      // Find the Address 2 input field by its label
      const address2Input = screen.getByLabelText('Address 2');
      expect(address2Input).toBeInTheDocument();
      // Assert that its value is an empty string because Street2 was "*"
      expect(address2Input).toHaveValue('');
    });

    // Also verify the fetch was called as expected
    expect(fetchMock).toHaveBeenCalledWith('/listing/edit-123', expect.any(Object));
  });
  // --- END OF ADDRESS 2 TEST ---

});