/*
 * For a detailed explanation regarding each configuration property and type check, visit:
 * https://jestjs.io/docs/configuration
 */

export default {
  clearMocks: true,
  collectCoverage: true,
  coverageDirectory: 'coverage',
  moduleNameMapper: {
    '^.+\\.(jpg|png|svg|s?css)$': '<rootDir>/src/views/__test__/mock/GenericLoaderMock.jest.ts',
  },
  testMatch: [
    // "**/__tests__/Pern.test.js"
    '**/?(*.)+(spec|test).[tj]s?(x)',
  ],
};
