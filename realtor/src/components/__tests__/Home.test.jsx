import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router';
import Home from '../Home';
import '@testing-library/jest-dom';

// Mock the fetch API
global.fetch = vi.fn(() =>
  Promise.resolve({
    json: () => Promise.resolve([]), // Mock empty listings for simplicity
    ok: true,
  })
);

// Mock TileDeck to avoid rendering its internals
vi.mock('../TileDeck', () => ({
  default: (props) => <div data-testid="tile-deck"></div>,
}));

describe('Home Component', () => {
  beforeEach(() => {
    fetch.mockClear();
  });

  it('renders Home component without crashing', async () => {
    render(
      <BrowserRouter>
        <Home />
      </BrowserRouter>
    );

    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
    expect(fetch).toHaveBeenCalledWith('/listings');
    expect(screen.getByTestId('tile-deck')).toBeInTheDocument();
  });
});