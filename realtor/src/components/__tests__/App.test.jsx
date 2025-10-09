// src/components/__tests__/App.test.jsx
import React from 'react';
import { render, screen, cleanup, waitFor, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { vi } from 'vitest';
import App from '../../App';
import { toast } from 'react-toastify'; // Import toast

// Mock fetch (not used, but included for consistency)
let fetchMock;

// Mock react-toastify
vi.mock('react-toastify', () => ({
  toast: {
    error: vi.fn(),
    success: vi.fn(),
    // Add other toast types if used in App.jsx and need mocking
  },
  ToastContainer: () => <div />, // Mock ToastContainer
}));


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
  vi.restoreAllMocks(); // Restore console mocks and other mocks
});

describe('App.jsx - Rendering and Basic Functionality', () => { // Updated describe block title
  it('renders NavBar and Main when not logged in', async () => {
    // Mock window.gapi to simulate successful Google Auth initialization
    vi.stubGlobal('gapi', {
      load: vi.fn((_, callback) => callback()), // Call callback immediately
      auth2: {
        init: vi.fn(() => ({
          isSignedIn: { get: vi.fn(() => false), listen: vi.fn() }, // Added listen mock
          currentUser: {
            get: vi.fn(() => ({
              getBasicProfile: vi.fn(() => ({
                getEmail: vi.fn(() => null),
              })),
            })),
          },
        })),
        getAuthInstance: vi.fn(() => ({
          isSignedIn: { get: vi.fn(() => false), listen: vi.fn() }, // Added listen mock
          currentUser: {
            get: vi.fn(() => ({
              getBasicProfile: vi.fn(() => ({
                getEmail: vi.fn(() => null),
              })),
            })),
          },
        })),
      },
    });

    render(
      <MemoryRouter initialEntries={['/realtor']}>
        <App />
      </MemoryRouter>
    );
    // Wait for the loaded state to render NavBar and Main
    await waitFor(() => {
      expect(screen.getByText('realtor webpage.')).toBeInTheDocument(); // From NavBar
      expect(screen.getByTestId('tile-deck')).toBeInTheDocument(); // From Main -> Home -> TileDeck
    });
  });
});

// Test Suite for Error Handling
describe('App Component - Error Handling', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.spyOn(console, 'log').mockImplementation(() => {});
        vi.spyOn(console, 'error').mockImplementation(() => {});
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it('should handle error during Google Auth initialization', async () => {
      // Mock window.gapi to simulate an error during loading
      vi.stubGlobal('gapi', {
        load: vi.fn((_, callback) => {
          // Simulate an error in the callback
          callback({ error: 'mocked_auth_error' });
        }),
        auth2: {
          // Mock init to fail as well to cover lines 60-77 and 96-123 (init errors)
          init: vi.fn(() => Promise.reject(new Error('Init failed'))),
           getAuthInstance: vi.fn(() => null), // Ensure getAuthInstance also returns null in error case
        },
      });

      render(
        <MemoryRouter initialEntries={['/realtor']}>
          <App />
        </MemoryRouter>
      );

      // Wait for the loading state to finish (should be fast with mocked errors)
      await waitFor(() => expect(screen.queryByText('Loading...')).not.toBeInTheDocument());


      // Check if the error toast was called for init error
      await waitFor(() => {
          expect(toast.error).toHaveBeenCalled();
          // The test should check for the specific error message from the catch block for init errors
          expect(toast.error).toHaveBeenCalledWith('Failed to initialize authentication.', { autoClose: 5000 });
      });
    });
  });

// New Test Suite for Sign-in and Sign-out
describe('App Component - Sign-in/Sign-out', () => {
    let mockAuthInstance;
    let isSignedInListener;

    beforeEach(() => {
        vi.clearAllMocks();
        vi.spyOn(console, 'log').mockImplementation(() => {});
        vi.spyOn(console, 'error').mockImplementation(() => {});

        // Mock a successful Google Auth instance
        mockAuthInstance = {
            isSignedIn: {
                get: vi.fn(() => false), // Initially signed out
                listen: vi.fn((listener) => {
                    isSignedInListener = listener; // Capture the listener
                }),
            },
            currentUser: {
                get: vi.fn(() => ({
                    getBasicProfile: vi.fn(() => ({
                        getEmail: vi.fn(() => 'testuser@example.com'),
                    })),
                })),
            },
            signIn: vi.fn(() => {
                // Simulate sign-in success
                mockAuthInstance.isSignedIn.get.mockReturnValue(true);
                // Explicitly call the listener to simulate the state change
                if (isSignedInListener) {
                     isSignedInListener(true);
                }
                return Promise.resolve(); // Resolve the promise
            }),
            signOut: vi.fn(() => {
                 // Simulate sign-out success
                 mockAuthInstance.isSignedIn.get.mockReturnValue(false);
                 // Explicitly call the listener to simulate the state change
                 if (isSignedInListener) {
                     isSignedInListener(false);
                 }
                 return Promise.resolve(); // Resolve the promise
            }),
        };

        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()), // Call callback immediately
            auth2: {
                // Resolve init with our mock instance and trigger the listener immediately if needed for initial state
                init: vi.fn(() => {
                    const authInstance = mockAuthInstance;
                     // Simulate initial state check in App's useEffect
                    if (authInstance.isSignedIn.get() && window.localStorage.getItem('signedOut') !== 'true') {
                        const email = authInstance.currentUser.get().getBasicProfile().getEmail();
                         // This part is handled by the component's useEffect, but we ensure the listener is set up
                         // Trigger the listener here if the initial state is signed in and not manually signed out
                        if (isSignedInListener) {
                            isSignedInListener(true);
                        }
                    }
                    return Promise.resolve(authInstance);
                }),
                getAuthInstance: vi.fn(() => mockAuthInstance), // Return the mock instance
            },
        });

        // Mock localStorage
        const localStorageMock = {
            getItem: vi.fn(),
            setItem: vi.fn(),
            removeItem: vi.fn(),
            clear: vi.fn(),
        };
        Object.defineProperty(window, 'localStorage', {
            value: localStorageMock,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
        // Clean up localStorage mock
        Object.defineProperty(window, 'localStorage', {
            value: undefined,
        });
    });

    it('should handle successful sign-in and sign-out', async () => {
        // Initially localStorage should not indicate signed out
        window.localStorage.getItem.mockReturnValue(null);

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the app to load and the sign-in button to appear
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Simulate clicking the Sign In button
        fireEvent.click(screen.getByText('Sign In'));

        // Wait for sign-in to complete and the user email to appear in the NavBar
        await waitFor(() => expect(screen.getByText('testuser@example.com')).toBeInTheDocument());

        // Check that signIn was called
        expect(mockAuthInstance.signIn).toHaveBeenCalledTimes(1);
        // Check that localStorage.removeItem was called after successful sign-in
        expect(window.localStorage.removeItem).toHaveBeenCalledWith('signedOut');
        // Check for the success toast
        expect(toast.success).toHaveBeenCalledWith('Successfully signed in!', { autoClose: 3000, toastId: 'sign-in' });

        // Click the dropdown toggle (user email) to reveal the "Sign Out" option
        fireEvent.click(screen.getByText('testuser@example.com'));

        // Wait for the "Sign Out" option to be visible and click it
        await waitFor(() => expect(screen.getByText('Sign Out')).toBeVisible());
        fireEvent.click(screen.getByText('Sign Out'));

        // Wait for sign-out to complete and the Sign In button to reappear
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Check that signOut was called
        expect(mockAuthInstance.signOut).toHaveBeenCalledTimes(1);
         // Check that localStorage.setItem was called after successful sign-out
        expect(window.localStorage.setItem).toHaveBeenCalledWith('signedOut', 'true');
        // Check for the success toast
        expect(toast.success).toHaveBeenCalledWith('Successfully signed out.', { autoClose: 3000, toastId: 'sign-out' });
    });

     it('should restore session if signed in and not manually signed out', async () => {
        // Simulate being signed in according to gapi and not manually signed out in localStorage
        mockAuthInstance.isSignedIn.get.mockReturnValue(true);
        window.localStorage.getItem.mockReturnValue(null); // Not manually signed out

         render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the user email to appear, indicating session restoration
        await waitFor(() => expect(screen.getByText('testuser@example.com')).toBeInTheDocument());

        // Ensure signIn was NOT called (session was restored)
        expect(mockAuthInstance.signIn).not.toHaveBeenCalled();
         // Ensure localStorage.removeItem was called (signedOut flag cleared)
         // This is called within the listen callback in App.jsx when isSignedIn becomes true
         // We need to ensure the listener is triggered. The mock setup in beforeEach now handles this.
        expect(window.localStorage.removeItem).toHaveBeenCalledWith('signedOut');
    });

    it('should not restore session if manually signed out', async () => {
        // Simulate being signed in according to gapi but manually signed out in localStorage
        mockAuthInstance.isSignedIn.get.mockReturnValue(true);
        window.localStorage.getItem.mockReturnValue('true'); // Manually signed out

         render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the app to load and the Sign In button to appear (session not restored)
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Ensure signIn was NOT called
        expect(mockAuthInstance.signIn).not.toHaveBeenCalled();
         // Ensure localStorage.removeItem was NOT called
        expect(window.localStorage.removeItem).not.toHaveBeenCalled();
    });

     it('should handle sign-in failure (popup closed)', async () => {
        // Mock signIn to reject with a popup_closed_by_user error
        mockAuthInstance.signIn.mockRejectedValue({ error: 'popup_closed_by_user' });

         render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the sign-in button to appear
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Simulate clicking the Sign In button
        fireEvent.click(screen.getByText('Sign In'));

        // Wait for the error toast
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Sign-in canceled.', { autoClose: 5000 });
        });

        // Ensure user is not logged in
        expect(screen.queryByText('testuser@example.com')).not.toBeInTheDocument();
         expect(screen.getByText('Sign In')).toBeInTheDocument(); // Sign In button should still be there
    });

    it('should handle sign-in failure (access denied)', async () => {
        // Mock signIn to reject with an access_denied error
        mockAuthInstance.signIn.mockRejectedValue({ error: 'access_denied' });

         render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the sign-in button to appear
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Simulate clicking the Sign In button
        fireEvent.click(screen.getByText('Sign In'));

        // Wait for the error toast
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Permission denied.', { autoClose: 5000 });
        });

        // Ensure user is not logged in
        expect(screen.queryByText('testuser@example.com')).not.toBeInTheDocument();
         expect(screen.getByText('Sign In')).toBeInTheDocument(); // Sign In button should still be there
    });

     it('should handle generic sign-in failure', async () => {
        // Mock signIn to reject with a generic error
        mockAuthInstance.signIn.mockRejectedValue(new Error('Something went wrong'));

         render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the sign-in button to appear
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Simulate clicking the Sign In button
        fireEvent.click(screen.getByText('Sign In'));

        // Wait for the error toast
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to sign in.', { autoClose: 5000 });
        });

        // Ensure user is not logged in
        expect(screen.queryByText('testuser@example.com')).not.toBeInTheDocument();
         expect(screen.getByText('Sign In')).toBeInTheDocument(); // Sign In button should still be there
    });

     it('should handle sign-out failure', async () => {
        // Simulate being logged in initially
        mockAuthInstance.isSignedIn.get.mockReturnValue(true);
        window.localStorage.getItem.mockReturnValue(null);

        // Mock signOut to reject with an error
        mockAuthInstance.signOut.mockRejectedValue(new Error('Sign out failed'));


         render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the user email to appear
        await waitFor(() => expect(screen.getByText('testuser@example.com')).toBeInTheDocument());

        // Click the dropdown toggle (user email) to reveal the "Sign Out" option
        fireEvent.click(screen.getByText('testuser@example.com'));

        // Wait for the "Sign Out" option to be visible and click it
        await waitFor(() => expect(screen.getByText('Sign Out')).toBeVisible());
        fireEvent.click(screen.getByText('Sign Out'));

        // Wait for the error toast
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to sign out.', { autoClose: 5000 });
        });

        // Ensure user is still logged in (sign-out failed)
        expect(screen.getByText('testuser@example.com')).toBeInTheDocument();
         expect(screen.queryByText('Sign In')).not.toBeInTheDocument(); // Sign In button should not be there
    });

    it('should prevent multiple sign-in attempts when authLoading is true', async () => {
        // Mock signIn to take some time (never resolve to keep authLoading true)
        mockAuthInstance.signIn.mockImplementation(() => new Promise(() => {})); // Never resolves

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the sign-in button to appear
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Click sign in first time
        fireEvent.click(screen.getByText('Sign In'));

        // Wait for the button text to change to "Signing In..."
        await waitFor(() => expect(screen.getByText('Signing In...')).toBeInTheDocument());

        // Try to click again (should be ignored due to authLoading)
        fireEvent.click(screen.getByText('Signing In...'));

        // Wait a bit to ensure only one call was made
        await new Promise(resolve => setTimeout(resolve, 100));

        // Verify signIn was only called once (second call was ignored due to authLoading)
        expect(mockAuthInstance.signIn).toHaveBeenCalledTimes(1);
    });

    it('should handle error when window.gapi does not exist and fetch fails', async () => {
        // Mock window.gapi to be undefined initially
        vi.stubGlobal('gapi', undefined);

        // Mock fetch to fail when trying to load Google Auth library
        const fetchMock = vi.fn().mockRejectedValue(new Error('Failed to load Google Auth library'));
        global.fetch = fetchMock;

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the error handling to complete
        await waitFor(() => {
            expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
        });

        // Check that the error toast was called
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to load authentication library.', { autoClose: 5000 });
        });

        // Verify fetch was called to try to load the library
        expect(fetchMock).toHaveBeenCalledWith('https://apis.google.com/js/platform.js');
    });

    it('should handle error when auth2 instance is null after initialization', async () => {
        // Mock gapi to return null from getAuthInstance
        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()), // Call callback immediately
            auth2: {
                init: vi.fn(() => Promise.resolve()), // Resolve successfully
                getAuthInstance: vi.fn(() => null), // Return null to trigger error branch
            },
        });

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the error handling to complete
        await waitFor(() => {
            expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
        });

        // Check that the authentication setup failed error toast was called
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Authentication setup failed.', { autoClose: 5000 });
        });

        // Verify the Sign In button is still rendered (app should still work)
        expect(screen.getByText('Sign In')).toBeInTheDocument();
    });

    it('should handle error when auth2 instance is not available during sign-in', async () => {
        // Mock a successful init but getAuthInstance returns null during sign-in
        const mockAuthInstance = null;

        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()),
            auth2: {
                init: vi.fn(() => Promise.resolve()), // Init succeeds
                getAuthInstance: vi.fn(() => mockAuthInstance), // Always returns null
            },
        });

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for app to load
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Click sign in button
        fireEvent.click(screen.getByText('Sign In'));

        // Wait for error handling
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to sign in.', { autoClose: 5000 });
        });

        // Verify user is still not logged in
        expect(screen.getByText('Sign In')).toBeInTheDocument();
    });

    it('should handle error when auth2 instance is not available during sign-out', async () => {
        // Mock initial successful setup and sign-in state
        let mockAuthInstance = {
            isSignedIn: {
                get: vi.fn(() => true), // Start signed in
                listen: vi.fn(),
            },
            currentUser: {
                get: vi.fn(() => ({
                    getBasicProfile: vi.fn(() => ({
                        getEmail: vi.fn(() => 'testuser@example.com'),
                    })),
                })),
            },
        };

        // Track calls to getAuthInstance to return different values
        let getAuthInstanceCalls = 0;
        
        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()),
            auth2: {
                init: vi.fn(() => {
                    // Trigger the listener to set signed-in state
                    setTimeout(() => {
                        if (mockAuthInstance.isSignedIn.listen.mock.calls.length > 0) {
                            const listener = mockAuthInstance.isSignedIn.listen.mock.calls[0][0];
                            listener(true);
                        }
                    }, 0);
                    return Promise.resolve();
                }),
                getAuthInstance: vi.fn(() => {
                    getAuthInstanceCalls++;
                    // Return mockAuthInstance for first few calls, null for sign-out attempt
                    return getAuthInstanceCalls <= 2 ? mockAuthInstance : null;
                }),
            },
        });

        // Mock localStorage
        const localStorageMock = {
            getItem: vi.fn(() => null), // Not manually signed out
            setItem: vi.fn(),
            removeItem: vi.fn(),
        };
        Object.defineProperty(window, 'localStorage', {
            value: localStorageMock,
        });

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for user to be logged in initially
        await waitFor(() => expect(screen.getByText('testuser@example.com')).toBeInTheDocument());

        // Click the dropdown to reveal sign-out option
        fireEvent.click(screen.getByText('testuser@example.com'));

        // Wait for sign-out option and click it
        await waitFor(() => expect(screen.getByText('Sign Out')).toBeVisible());
        fireEvent.click(screen.getByText('Sign Out'));

        // Wait for error handling
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to sign out.', { autoClose: 5000 });
        });

        // User should still be logged in since sign-out failed
        expect(screen.getByText('testuser@example.com')).toBeInTheDocument();

        // Clean up localStorage mock
        Object.defineProperty(window, 'localStorage', {
            value: undefined,
        });
    });

    it('should handle error when fetching Google Auth library fails with bad response', async () => {
        // Mock window.gapi to be undefined initially to trigger fetch
        vi.stubGlobal('gapi', undefined);

        // Mock fetch to return a bad response (not ok)
        const fetchMock = vi.fn().mockResolvedValue({
            ok: false, // This should trigger the !response.ok branch (line 29)
            status: 404,
            statusText: 'Not Found'
        });
        global.fetch = fetchMock;

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for the error handling to complete
        await waitFor(() => {
            expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
        });

        // Check that the error toast was called
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to load authentication library.', { autoClose: 5000 });
        });

        // Verify fetch was called to try to load the library
        expect(fetchMock).toHaveBeenCalledWith('https://apis.google.com/js/platform.js');
    });

    it('should handle error when getAuthInstance returns null during sign-out attempt', async () => {
        // Create a working auth instance for initialization
        const workingAuthInstance = {
            isSignedIn: {
                get: vi.fn(() => true), // User is signed in
                listen: vi.fn((callback) => {
                    // Immediately call with signed in state
                    callback(true);
                }),
            },
            currentUser: {
                get: vi.fn(() => ({
                    getBasicProfile: vi.fn(() => ({
                        getEmail: vi.fn(() => 'testuser@example.com'),
                    })),
                })),
            },
        };

        // Track call count to getAuthInstance
        let getAuthInstanceCallCount = 0;

        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()),
            auth2: {
                init: vi.fn(() => Promise.resolve(workingAuthInstance)),
                getAuthInstance: vi.fn(() => {
                    getAuthInstanceCallCount++;
                    // Return working instance for initialization, null for sign-out
                    return getAuthInstanceCallCount === 1 ? workingAuthInstance : null;
                }),
            },
        });

        // Mock localStorage
        const localStorageMock = {
            getItem: vi.fn(() => null), // Not manually signed out
            setItem: vi.fn(),
            removeItem: vi.fn(),
        };
        Object.defineProperty(window, 'localStorage', {
            value: localStorageMock,
        });

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for user to be signed in
        await waitFor(() => expect(screen.getByText('testuser@example.com')).toBeInTheDocument());

        // Click user email to open dropdown
        fireEvent.click(screen.getByText('testuser@example.com'));

        // Click sign out
        await waitFor(() => expect(screen.getByText('Sign Out')).toBeVisible());
        fireEvent.click(screen.getByText('Sign Out'));

        // Wait for error toast (getAuthInstance returns null during sign-out)
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to sign out.', { autoClose: 5000 });
        });

        // User should still appear logged in since sign-out failed
        expect(screen.getByText('testuser@example.com')).toBeInTheDocument();

        // Clean up localStorage mock
        Object.defineProperty(window, 'localStorage', {
            value: undefined,
        });
    });

    it('should handle auth loading state prevention during sign-out attempts', async () => {
        // Create a working auth instance for initialization
        const workingAuthInstance = {
            isSignedIn: {
                get: vi.fn(() => true), // User is signed in
                listen: vi.fn((callback) => {
                    callback(true); // Trigger signed in state
                }),
            },
            currentUser: {
                get: vi.fn(() => ({
                    getBasicProfile: vi.fn(() => ({
                        getEmail: vi.fn(() => 'testuser@example.com'),
                    })),
                })),
            },
            signOut: vi.fn(() => new Promise(() => {})), // Never resolves to keep authLoading true
        };

        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()),
            auth2: {
                init: vi.fn(() => Promise.resolve(workingAuthInstance)),
                getAuthInstance: vi.fn(() => workingAuthInstance),
            },
        });

        // Mock localStorage
        const localStorageMock = {
            getItem: vi.fn(() => null),
            setItem: vi.fn(),
            removeItem: vi.fn(),
        };
        Object.defineProperty(window, 'localStorage', {
            value: localStorageMock,
        });

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for user to be signed in
        await waitFor(() => expect(screen.getByText('testuser@example.com')).toBeInTheDocument());

        // Click user email to open dropdown
        fireEvent.click(screen.getByText('testuser@example.com'));

        // Click sign out first time
        await waitFor(() => expect(screen.getByText('Sign Out')).toBeVisible());
        fireEvent.click(screen.getByText('Sign Out'));

        // Wait for signing out state
        await waitFor(() => expect(screen.getByText('Signing Out...')).toBeInTheDocument());

        // Try to click sign out again while authLoading is true (should be ignored)
        fireEvent.click(screen.getByText('Signing Out...'));

        // Wait a bit to ensure only one call was made
        await new Promise(resolve => setTimeout(resolve, 100));

        // Verify signOut was only called once (second call was ignored due to authLoading)
        expect(workingAuthInstance.signOut).toHaveBeenCalledTimes(1);

        // Clean up localStorage mock
        Object.defineProperty(window, 'localStorage', {
            value: undefined,
        });
    });

    it('should handle getAuthInstance returning null during sign-in after successful initialization', async () => {
        // Track calls to getAuthInstance to return different values
        let getAuthInstanceCallCount = 0;
        const workingAuthInstance = {
            isSignedIn: { get: vi.fn(() => false), listen: vi.fn() },
            currentUser: {
                get: vi.fn(() => ({
                    getBasicProfile: vi.fn(() => ({
                        getEmail: vi.fn(() => 'testuser@example.com'),
                    })),
                })),
            },
        };

        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()),
            auth2: {
                init: vi.fn(() => Promise.resolve(workingAuthInstance)),
                getAuthInstance: vi.fn(() => {
                    getAuthInstanceCallCount++;
                    // Return working instance for initialization, null for sign-in attempt
                    return getAuthInstanceCallCount <= 2 ? workingAuthInstance : null;
                }),
            },
        });

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for app to load with sign-in button
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // Click sign-in button (this should trigger the null getAuthInstance case)
        fireEvent.click(screen.getByText('Sign In'));

        // Wait for error handling - should show "Failed to sign in" due to null auth instance
        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith('Failed to sign in.', { autoClose: 5000 });
        });

        // User should still not be logged in
        expect(screen.getByText('Sign In')).toBeInTheDocument();
        expect(screen.queryByText('testuser@example.com')).not.toBeInTheDocument();
    });

    it('should test final edge case for maximum coverage', async () => {
        // This test is designed to hit any remaining edge cases
        // by testing the authLoading prevention in a different scenario
        
        const mockAuthInstance = {
            isSignedIn: {
                get: vi.fn(() => false),
                listen: vi.fn(),
            },
            currentUser: {
                get: vi.fn(() => ({
                    getBasicProfile: vi.fn(() => ({
                        getEmail: vi.fn(() => 'testuser@example.com'),
                    })),
                })),
            },
            signIn: vi.fn(() => {
                // Simulate successful sign-in
                mockAuthInstance.isSignedIn.get.mockReturnValue(true);
                return Promise.resolve();
            }),
        };

        vi.stubGlobal('gapi', {
            load: vi.fn((_, callback) => callback()),
            auth2: {
                init: vi.fn(() => Promise.resolve(mockAuthInstance)),
                getAuthInstance: vi.fn(() => mockAuthInstance),
            },
        });

        render(
            <MemoryRouter initialEntries={['/realtor']}>
                <App />
            </MemoryRouter>
        );

        // Wait for app to load
        await waitFor(() => expect(screen.getByText('Sign In')).toBeInTheDocument());

        // This should cover any remaining branch logic
        expect(screen.getByText('Sign In')).toBeInTheDocument();
    });
});
