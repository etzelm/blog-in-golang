# sample realtor react app

This project is built with [Vite](https://vitejs.dev/) and [React](https://reactjs.org/). Previously bootstrapped with Create React App, it has been migrated to Vite for faster development and build times.

## Code Coverage (v8 Report)

| File           | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s   |
| :------------- | :------ | :------- | :------ | :------ | :------------------ |
| **All files** | **91.24** | **90.22** | **92.1** | **91.24** |                     |
| `src`          | 80.18   | 81.81    | 80      | 80.18   |                     |
| `App.jsx`      | 77.89   | 81.81    | 80      | 77.89   | 67-87,97-98         |
| `index.jsx`    | 100     | 100      | 100     | 100     |                     |
| `src/components` | 92.21   | 90.65    | 93.93   | 92.21   |                     |
| `Home.jsx`     | 100     | 100      | 100     | 100     |                     |
| `Listing.jsx`  | 94.89   | 84.21    | 100     | 94.89   | 32-34,42-43,45-46   |
| `Main.jsx`     | 100     | 100      | 100     | 100     |                     |
| `MyListing.jsx`| 85.21   | 92.15    | 80      | 85.21   | ...,270-292,301-310 |
| `MyListings.jsx`| 94.59  | 91.66    | 100     | 94.59   | 17-18               |
| `NavBar.jsx`   | 100     | 75       | 100     | 100     | 20                  |
| `Search.jsx`   | 98.39   | 85.29    | 100     | 98.39   | 83-85               |
| `Tile.jsx`     | 100     | 96       | 100     | 100     | 9                   |
| `TileDeck.jsx` | 92.3    | 83.33    | 100     | 92.3    | 14-15               |
| `test-data`    | 100     | 100      | 100     | 100     |                     |
| `index.js`     | 100     | 100      | 100     | 100     |                     |

## Available Scripts

In the project directory, you can run:

### `yarn start` or `yarn dev`

Runs the app in the development mode using Vite.
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits with fast Hot Module Replacement (HMR).

### `yarn test`

Launches the Vitest test runner.
See the [Vitest documentation](https://vitest.dev/) for more information.

### `yarn test:watch`

Runs tests in watch mode with Vitest.

### `yarn build`

Builds the app for production to the `build` folder using Vite.
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.
Your app is ready to be deployed!

### `yarn preview`

Locally preview the production build using Vite's preview server.

## Learn More

You can learn more in the [Vite documentation](https://vitejs.dev/guide/).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

Vite supports dynamic imports for code splitting. See: <https://vitejs.dev/guide/features.html#dynamic-import>

### Analyzing the Bundle Size

Use `yarn build` and check the build output for bundle analysis. You can also use tools like `rollup-plugin-visualizer`.

### Environment Variables

Vite uses `VITE_` prefixed environment variables. See: <https://vitejs.dev/guide/env-and-mode.html>

### Advanced Configuration

Customize Vite configuration in `vite.config.js`. See: <https://vitejs.dev/config/>

### Deployment

See the Vite deployment guide: <https://vitejs.dev/guide/static-deploy.html>
