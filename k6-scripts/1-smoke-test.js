/**
 * K6 Smoke Test - Go Fiber Service
 * 
 * Purpose: Verify basic functionality with minimal load
 * VUs: 2
 * Duration: 2 minutes
 * 
 * Run: k6 run 1-smoke-test.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 2,
  duration: '2m',
  
  thresholds: {
    http_req_duration: ['p(95)<500'],    // 95% of requests < 500ms
    http_req_failed: ['rate<0.01'],      // Error rate < 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export function setup() {
  console.log('🔥 Starting Smoke Test...');
  console.log(`Target: ${BASE_URL}`);
  
  // Login to get token for authenticated tests
  const loginRes = http.post(
    `${BASE_URL}/api/v1/auth/login`,
    JSON.stringify({
      email: 'admin@example.com',
      password: 'password',
    }),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );
  
  if (loginRes.status !== 200) {
    console.error('❌ Login failed! Check your test user credentials');
    console.error(`Response: ${loginRes.body}`);
  }
  
  const data = JSON.parse(loginRes.body);
  return { 
    accessToken: data.data?.access_token || '',
    refreshToken: data.data?.refresh_token || '',
  };
}

export default function(data) {
  // Test 1: Login
  const loginRes = http.post(
    `${BASE_URL}/api/v1/auth/login`,
    JSON.stringify({
      email: 'admin@example.com',
      password: 'password',
    }),
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: 'Login' },
    }
  );
  
  check(loginRes, {
    'login: status is 200': (r) => r.status === 200,
    'login: has access token': (r) => {
      const body = JSON.parse(r.body);
      return body.data?.access_token !== undefined;
    },
  });
  
  sleep(1);
  
  // Test 2: Get Users List (authenticated)
  if (data.accessToken) {
    const usersRes = http.get(
      `${BASE_URL}/api/v1/users?page=1&page_size=10`,
      {
        headers: {
          'Authorization': `Bearer ${data.accessToken}`,
        },
        tags: { name: 'GetUsers' },
      }
    );
    
    check(usersRes, {
      'get users: status is 200': (r) => r.status === 200,
      'get users: has data array': (r) => {
        const body = JSON.parse(r.body);
        return Array.isArray(body.data);
      },
    });
  }
  
  sleep(1);
}

export function teardown(data) {
  console.log('✅ Smoke Test Complete!');
}
