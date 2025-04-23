import { expect, test, vi } from 'vitest';

// Mock root element
document.body.innerHTML = '<div id="root"></div>';

// Pre-mock react-dom/client before import
vi.mock('react-dom/client', async () => {
  const actual = await vi.importActual('react-dom/client');
  return {
    ...actual,
    createRoot: vi.fn(() => ({
      render: vi.fn()
    }))
  };
});

test('renders index.jsx without crashing', async () => {
  await import('../../index.jsx');
  const root = document.getElementById('root');
  expect(root).not.toBeNull();

  const { createRoot } = await import('react-dom/client');
  expect(createRoot).toHaveBeenCalled();
});
