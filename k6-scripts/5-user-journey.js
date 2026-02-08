/**
 * K6 User Journey Test - Go Fiber Service
 * 
 * Purpose: Simulate realistic user flow
 * Scenario: Login → Browse → Create → Update → Delete → Logout
 * VUs: 20 concurrent users
 * Duration: 10 minutes
 * 
 * Run: k6 run 5-user-journey.js
 */

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { Counter, Trend } from 'k6/metrics';

const journeyCompleted = new Counter('journey_completed');
const journeyFailed = new Counter('journey_failed');
const journeyDuration = new Trend('journey_duration');

export const options = {
  stages: [
    { duration: '2m', target: 20 },   // Ramp up
    { duration: '6m', target: 20 },   // Sustained
    { duration: '2m', target: 0 },    // Ramp down
  ],
  
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    http_req_failed: ['rate<0.05'],
    journey_completed: ['count>100'],
    'group_duration{group:::01_Authentication}': ['p(95)<1000'],
    'group_duration{group:::02_Browse}': ['p(95)<800'],
    'group_duration{group:::03_CRUD_Operations}': ['p(95)<3000'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export function setup() {
  console.log('👤 Starting User Journey Test...');
  console.log(`Simulating realistic user behavior`);
}

export default function() {
  const journeyStart = Date.now();
  let accessToken = '';
  let refreshToken = '';
  let createdUserId = '';
  
  try {
    // ========================================
    // PHASE 1: Authentication
    // ========================================
    group('01_Authentication', function() {
      // Login
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
      
      check(loginRes, {
        'login successful': (r) => r.status === 200,
      });
      
      if (loginRes.status !== 200) {
        journeyFailed.add(1);
        return;
      }
      
      const tokens = JSON.parse(loginRes.body).data;
      accessToken = tokens.access_token;
      refreshToken = tokens.refresh_token;
      
      sleep(1);
    });
    
    const authHeaders = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${accessToken}`,
    };
    
    // ========================================
    // PHASE 2: Browse & Search
    // ========================================
    group('02_Browse', function() {
      // List users - page 1
      const listRes1 = http.get(
        `${BASE_URL}/api/v1/users?page=1&page_size=10`,
        { headers: authHeaders }
      );
      
      check(listRes1, {
        'browse: list page 1 success': (r) => r.status === 200,
      });
      
      sleep(2); // User reading...
      
      // List users - page 2
      const listRes2 = http.get(
        `${BASE_URL}/api/v1/users?page=2&page_size=10`,
        { headers: authHeaders }
      );
      
      check(listRes2, {
        'browse: list page 2 success': (r) => r.status === 200,
      });
      
      sleep(1);
      
      // Search users
      const searchRes = http.get(
        `${BASE_URL}/api/v1/users?search=test&page=1&page_size=20`,
        { headers: authHeaders }
      );
      
      check(searchRes, {
        'browse: search success': (r) => r.status === 200,
      });
      
      sleep(2);
    });
    
    // ========================================
    // PHASE 3: CRUD Operations
    // ========================================
    group('03_CRUD_Operations', function() {
      // CREATE User
      const createPayload = JSON.stringify({
        username: `user_${randomString(8)}`,
        name: `Test User ${randomString(5)}`,
        email: `${randomString(10)}@test.com`,
        password: 'Testpassword!',
      });
      
      const createRes = http.post(
        `${BASE_URL}/api/v1/users`,
        createPayload,
        { headers: authHeaders }
      );
      
      const createSuccess = check(createRes, {
        'crud: create user success': (r) => r.status === 201 || r.status === 200,
      });
      
      if (!createSuccess) {
        console.warn('Create failed, skipping rest of CRUD');
        return;
      }
      
      createdUserId = JSON.parse(createRes.body).data.id;
      
      sleep(1);
      
      // READ User by ID
      const readRes = http.get(
        `${BASE_URL}/api/v1/users/${createdUserId}`,
        { headers: authHeaders }
      );
      
      check(readRes, {
        'crud: read user success': (r) => r.status === 200,
        'crud: correct user returned': (r) => {
          const data = JSON.parse(r.body).data;
          return data.id === createdUserId;
        },
      });
      
      sleep(2);
      
      // UPDATE User
      const updatePayload = JSON.stringify({
        name: `Updated Name ${randomString(5)}`,
      });
      
      const updateRes = http.put(
        `${BASE_URL}/api/v1/users/${createdUserId}`,
        updatePayload,
        { headers: authHeaders }
      );
      
      check(updateRes, {
        'crud: update user success': (r) => r.status === 200,
      });
      
      sleep(1);
      
      // DELETE User
      const deleteRes = http.del(
        `${BASE_URL}/api/v1/users/${createdUserId}`,
        null,
        { headers: authHeaders }
      );
      
      check(deleteRes, {
        'crud: delete user success': (r) => r.status === 200,
      });
      
      sleep(1);
    });
    
    // ========================================
    // PHASE 4: Token Refresh
    // ========================================
    group('04_Token_Refresh', function() {
      const refreshRes = http.post(
        `${BASE_URL}/api/v1/auth/refresh`,
        JSON.stringify({
          refresh_token: refreshToken,
        }),
        {
          headers: { 'Content-Type': 'application/json' },
        }
      );
      
      check(refreshRes, {
        'token refresh success': (r) => r.status === 200,
      });
      
      if (refreshRes.status === 200) {
        const newTokens = JSON.parse(refreshRes.body).data;
        accessToken = newTokens.access_token;
        refreshToken = newTokens.refresh_token;
      }
      
      sleep(1);
    });
    
    // ========================================
    // PHASE 5: Logout
    // ========================================
    group('05_Logout', function() {
      const logoutRes = http.post(
        `${BASE_URL}/api/v1/auth/logout`,
        JSON.stringify({
          refresh_token: refreshToken,
        }),
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`,
          },
        }
      );
      
      check(logoutRes, {
        'logout success': (r) => r.status === 200,
      });
    });
    
    // Journey completed successfully
    const journeyTime = Date.now() - journeyStart;
    journeyDuration.add(journeyTime);
    journeyCompleted.add(1);
    
    console.log(`✅ Journey completed in ${journeyTime}ms`);
    
  } catch (error) {
    journeyFailed.add(1);
    console.error(`❌ Journey failed: ${error}`);
  }
  
  // Think time before next iteration
  sleep(5);
}

export function teardown(data) {
  console.log('✅ User Journey Test Complete!');
  console.log('');
  console.log('📊 Journey Breakdown:');
  console.log('1. Authentication - Login flow');
  console.log('2. Browse - List & search users');
  console.log('3. CRUD - Create, Read, Update, Delete');
  console.log('4. Token Refresh - Session management');
  console.log('5. Logout - Clean session termination');
}
