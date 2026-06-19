/**
 * Test setup file
 * 
 * This file runs before all tests.
 * Use it to set up test environment, global mocks, etc.
 */

// Check if API server is running for E2E tests
const API_URL = process.env.API_URL || 'http://localhost:3001';

if (process.env.E2E_TEST) {
  console.log(`E2E tests will run against: ${API_URL}`);
  console.log('Make sure the API server is running with test database.');
}
