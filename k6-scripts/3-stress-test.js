/**
 * K6 Stress Test - Go Fiber Service
 * 
 * Purpose: Find system breaking point
 * VUs: 0 → 10 → 50 → 100 → 200 → 300
 * Duration: 20 minutes
 * 
 * Run: k6 run 3-stress-test.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate } from 'k6/metrics';

const errors = new Counter('errors');
const timeouts = new Counter('timeouts');
const slowRequests = new Counter('slow_requests');
const errorRate = new Rate('error_rate');

export const options = {
  stages: [
    { duration: '2m', target: 10 },    // Warm up
    { duration: '3m', target: 50 },    // Normal load
    { duration: '3m', target: 100 },   // Push harder
    { duration: '3m', target: 200 },   // Heavy load
    { duration: '3m', target: 300 },   // Breaking point?
    { duration: '3m', target: 300 },   // Sustained stress
    { duration: '3m', target: 0 },     // Recovery
  ],
  
  thresholds: {
    http_req_duration: ['p(95)<5000'],   // More lenient
    http_req_failed: ['rate<0.15'],      // Accept 15% errors
    error_rate: ['rate<0.20'],           // Max 20% error rate
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export function setup() {
  console.log('💪 Starting Stress Test...');
  console.log(`Target: ${BASE_URL}`);
  console.log(`⚠️  This will push your system to its limits!`);
  
  // Get initial token
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
  
  if (loginRes.status === 200) {
    const data = JSON.parse(loginRes.body);
    return { accessToken: data.data.access_token };
  }
  
  console.error('❌ Setup failed - cannot proceed');
  throw new Error('Setup login failed');
}

export default function(data) {
  const authHeaders = {
    'Authorization': `Bearer ${data.accessToken}`,
  };
  
  // Simulate various user behaviors
  const scenarios = [
    // Scenario 1: List users (most common)
    () => http.get(`${BASE_URL}/api/v1/users?page=1&page_size=10`, { headers: authHeaders }),
    
    // Scenario 2: Search users
    () => http.get(`${BASE_URL}/api/v1/users?search=test`, { headers: authHeaders }),
    
    // Scenario 3: Get user by ID (random)
    () => {
      const randomId = '00000000-0000-0000-0000-000000000001'; // Replace with actual ID
      return http.get(`${BASE_URL}/api/v1/users/${randomId}`, { headers: authHeaders });
    },
    
    // Scenario 4: Pagination
    () => {
      const page = Math.floor(Math.random() * 10) + 1;
      return http.get(`${BASE_URL}/api/v1/users?page=${page}&page_size=20`, { headers: authHeaders });
    },
  ];
  
  // Pick random scenario
  const scenario = scenarios[Math.floor(Math.random() * scenarios.length)];
  const startTime = Date.now();
  const res = scenario();
  const duration = Date.now() - startTime;
  
  // Track metrics
  const success = check(res, {
    'status is 2xx or 3xx': (r) => r.status >= 200 && r.status < 400,
    'not timeout': (r) => r.status !== 0,
    'response < 10s': () => duration < 10000,
  });
  
  if (!success) {
    errors.add(1);
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  if (res.status === 0) {
    timeouts.add(1);
    console.error('⏱️  Request timeout!');
  }
  
  if (duration > 5000) {
    slowRequests.add(1);
    console.warn(`⚠️  Slow request: ${duration}ms`);
  }
  
  // Aggressive think time (more requests per second)
  sleep(Math.random() * 2);
}

export function teardown(data) {
  console.log('✅ Stress Test Complete!');
  console.log('');
  console.log('📊 Analysis Tips:');
  console.log('- Check when error rate spiked');
  console.log('- Identify breaking point (VU count)');
  console.log('- Monitor system resources (CPU, Memory)');
  console.log('- Check database connection pool');
}
