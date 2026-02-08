/**
 * K6 Load Test - Go Fiber Service
 * 
 * Purpose: Test normal production load
 * VUs: 10 → 50 → 10
 * Duration: 10 minutes
 * 
 * Run: k6 run 2-load-test.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

// Custom metrics
const loginDuration = new Trend('custom_login_duration');
const loginSuccessRate = new Rate('custom_login_success');
const authErrors = new Counter('custom_auth_errors');

export const options = {
  stages: [
    { duration: '2m', target: 10 },   // Ramp up to 10 users
    { duration: '5m', target: 50 },   // Ramp up to 50 users
    { duration: '2m', target: 50 },   // Stay at 50 users
    { duration: '1m', target: 0 },    // Ramp down
  ],
  
  thresholds: {
    http_req_duration: ['p(95)<1000', 'p(99)<2000'],
    http_req_failed: ['rate<0.05'],
    'http_req_duration{scenario:login}': ['p(95)<800'],
    'http_req_duration{scenario:users}': ['p(95)<500'],
    custom_login_success: ['rate>0.95'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

// Test users pool
const TEST_USERS = [
  { email: 'user1@example.com', password: 'password' },
  { email: 'user2@example.com', password: 'password' },
  { email: 'user3@example.com', password: 'password' },
  { email: 'user4@example.com', password: 'password' },
  { email: 'user5@example.com', password: 'password' },
];

export function setup() {
  console.log('📊 Starting Load Test...');
  console.log(`Target: ${BASE_URL}`);
  console.log(`Stages: 10→50→10 users over 10 minutes`);
}

export default function() {
  // Select random user
  const user = TEST_USERS[Math.floor(Math.random() * TEST_USERS.length)];
  
  // Scenario 1: Login
  const loginStart = Date.now();
  
  const loginRes = http.post(
    `${BASE_URL}/api/v1/auth/login`,
    JSON.stringify(user),
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { scenario: 'login' },
    }
  );
  
  const loginDur = Date.now() - loginStart;
  loginDuration.add(loginDur);
  
  const loginSuccess = check(loginRes, {
    'login: status is 200': (r) => r.status === 200,
    'login: response time < 1s': () => loginDur < 1000,
  });
  
  loginSuccessRate.add(loginSuccess);
  
  if (!loginSuccess) {
    authErrors.add(1);
    sleep(2);
    return;
  }
  
  const tokens = JSON.parse(loginRes.body).data;
  const accessToken = tokens.access_token;
  
  sleep(1);
  
  // Scenario 2: Get Users List
  const usersRes = http.get(
    `${BASE_URL}/api/v1/users?page=1&page_size=20`,
    {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
      tags: { scenario: 'users' },
    }
  );
  
  check(usersRes, {
    'users list: status is 200': (r) => r.status === 200,
    'users list: has pagination': (r) => {
      const body = JSON.parse(r.body);
      return body.meta !== undefined;
    },
  });
  
  sleep(2);
  
  // Scenario 3: Search Users
  const searchRes = http.get(
    `${BASE_URL}/api/v1/users?search=test&page=1&page_size=10`,
    {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
      tags: { scenario: 'users' },
    }
  );
  
  check(searchRes, {
    'search users: status is 200': (r) => r.status === 200,
  });
  
  sleep(1);
  
  // Scenario 4: Refresh Token
  const refreshRes = http.post(
    `${BASE_URL}/api/v1/auth/refresh`,
    JSON.stringify({
      refresh_token: tokens.refresh_token,
    }),
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { scenario: 'login' },
    }
  );
  
  check(refreshRes, {
    'refresh: status is 200': (r) => r.status === 200,
    'refresh: new token received': (r) => {
      const body = JSON.parse(r.body);
      return body.data?.access_token !== undefined;
    },
  });
  
  sleep(3);
}

export function teardown(data) {
  console.log('✅ Load Test Complete!');
  console.log('📈 Check the summary for detailed metrics');
}
