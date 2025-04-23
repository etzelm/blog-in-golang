// src/components/__tests__/Listing.test.jsx

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, waitFor } from '@testing-library/react';
import { BrowserRouter, MemoryRouter, Route, Routes } from 'react-router'; // Use react-router-dom for Routes and Route
import Listing from '../Listing';
import '@testing-library/jest-dom';

// Mock fetch
let fetchMock;

beforeEach(() => {
  fetchMock = vi.fn(); // Use let for reassignment
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

  it('handles no MLS ID in URL when no listing prop is provided', async () => {
    render(
      <MemoryRouter initialEntries={['/realtor/listing']}>
        <Routes>
          <Route path="/realtor/listing" element={<Listing />} />
        </Routes>
      </MemoryRouter>
    );

    // Wait for a common element to appear before checking fetch
    await waitFor(() => expect(screen.getByText(/last updated/i)).toBeInTheDocument());
    // No MLS in URL, fetch should not be called
    expect(fetchMock).not.toHaveBeenCalled();
    // No listing data, should not render listing details
    expect(screen.queryByText(/Price:/)).not.toBeInTheDocument();
  });

  it('does not render carousel images when Photo Array is null', async () => {
    const { listings } = await import('../../../test-data');
    const mockListing = { ...listings[0], 'Photo Array': null }; // Photo Array is null

    fetchMock.mockImplementationOnce((url) => {
      if (url === `/listing/${mockListing.MLS}`) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([mockListing]),
        });
      }
      return Promise.resolve({ ok: false, json: () => Promise.resolve({}) });
    });

    const { container } = render(
      <MemoryRouter initialEntries={[`/realtor/listing?MLS=${mockListing.MLS}`]}>
        <Routes>
          <Route path="/realtor/listing" element={<Listing />} />
        </Routes>
      </MemoryRouter>
    );

    // Wait for a common element to appear before checking carousel
    await waitFor(() => expect(screen.getByText(/last updated/i)).toBeInTheDocument());
    await waitFor(() => {
      const carouselElement = container.querySelector('.carousel.slide');
      expect(carouselElement).toBeInTheDocument();
      const carouselImages = container.querySelectorAll('.carousel-item img');
      expect(carouselImages.length).toBe(0); // Should be 0 images
    });
  });

  it('does not render carousel images when Photo Array is an empty array', async () => {
    const { listings } = await import('../../../test-data');
    const mockListing = { ...listings[0], 'Photo Array': [] }; // Photo Array is empty

    fetchMock.mockImplementationOnce((url) => {
      if (url === `/listing/${mockListing.MLS}`) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([mockListing]),
        });
      }
      return Promise.resolve({ ok: false, json: () => Promise.resolve({}) });
    });

    const { container } = render(
      <MemoryRouter initialEntries={[`/realtor/listing?MLS=${mockListing.MLS}`]}>
        <Routes>
          <Route path="/realtor/listing" element={<Listing />} />
        </Routes>
      </MemoryRouter>
    );

    // Wait for a common element to appear before checking carousel
    await waitFor(() => expect(screen.getByText(/last updated/i)).toBeInTheDocument());
    await waitFor(() => {
      const carouselElement = container.querySelector('.carousel.slide');
      expect(carouselElement).toBeInTheDocument();
      const carouselImages = container.querySelectorAll('.carousel-item img');
      expect(carouselImages.length).toBe(0); // Should be 0 images
    });
  });
});