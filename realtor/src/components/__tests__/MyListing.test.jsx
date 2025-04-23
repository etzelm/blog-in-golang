// src/components/__tests__/MyListing.test.jsx

import React from 'react'; // Added React
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'; // Added vi
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/react'; // Added fireEvent, waitFor
import { BrowserRouter, MemoryRouter, Route, Routes } from 'react-router'; // Added MemoryRouter, Route, Routes, updated react-router import
import MyListing from '../MyListing';
import '@testing-library/jest-dom';
import { NotificationManager } from 'react-notifications'; // Added NotificationManager

// --- Mocks ---
// Mock react-notifications (at top level)
vi.mock('react-notifications', () => ({
  NotificationManager: {
    success: vi.fn(),
    warning: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
  },
  NotificationContainer: () => <div data-testid="notification-container" />, // Mock Container rendering
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
  'Photo Array': ['https://example.com/photo1.jpg'],
  User: 'test@example.com',
  deleted: 'false',
};

describe('MyListing Component - Edit Mode', () => {
  // Specific setup for this describe block
  beforeEach(() => {
    // Clear mocks specifically for this block
    fetchMock.mockClear();
    vi.clearAllMocks(); // Clears mocks like NotificationManager too

    // Mock fetch for getting the listing data for THIS block
    fetchMock.mockImplementation((url) => {
      if (url.startsWith('/listing/')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([mockListingData]), // Return the mock data in an array
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
       expect(NotificationManager.success).toHaveBeenCalledWith('Success', 'Listing submitted', 3000);
    });

     // 6. Check that dropzone is rendered (since user is logged in)
     expect(screen.getAllByTestId('dropzone')).toHaveLength(2);
  });
});