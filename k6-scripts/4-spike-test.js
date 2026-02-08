/**
 * K6 Spike Test - Go Fiber Service
 * 
 * Purpose: Test sudden traffic spike (flash sale, viral content)
 * VUs: 10 → 500 (instant!) → 10
 * Duration: 5 minutes
 * 
 * Run: k6 run 4-spike-test.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Gauge } from 'k6/metrics';

const survivedSpike = new Counter('survived_spike');
const failedDuringSpike = new Counter('failed_during_spike');
const currentVUs = new Gauge('current_vus');

export const options = {
  stages: [
    { duration: '1m', target: 10 },     // Normal baseline
    { duration: '10s', target: 500 },   // 🚀 SUDDEN SPIKE!
    { duration: '2m', target: 500 },    // Sustained spike
    { duration: '1m', target: 10 },     // Back to normal
    { duration: '1m', target: 0 },      // Cool down
  ],
  
  thresholds: {
    http_req_duration: ['p(95)<3000'],     // Lenient during spike
    http_req_failed: ['rate<0.20'],        // Accept 20% errors
    survived_spike: ['count>1000'],        // At least 1000 successful requests
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export function setup() {
  console.log('⚡ Starting Spike Test...');
  console.log(`Target: ${BASE_URL}`);
  console.log(`Spike: 10 → 500 users in 10 seconds!`);
  
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
  
  throw new Error('Setup login failed');
}

export default function(data) {
  currentVUs.add(__VU);
  
  const startTime = Date.now();
  
  // Simple endpoint - most likely to survive spike
  const res = http.get(
    `${BASE_URL}/api/v1/users?page=1&page_size=10`,
    {
      headers: {
        'Authorization': `Bearer ${data.accessToken}`,
      },
    }
  );
  
  const duration = Date.now() - startTime;
  
  const success = check(res, {
    'survived spike': (r) => r.status === 200,
    'response within timeout': () => duration < 5000,
  });
  
  if (success) {
    survivedSpike.add(1);
  } else {
    failedDuringSpike.add(1);
    
    // Log failures during spike for debugging
    if (__VU > 100) {  // Only during spike
      console.error(`❌ Spike failure: Status ${res.status}, Duration ${duration}ms, VU ${__VU}`);
    }
  }
  
  // Minimal think time - maximum pressure
  sleep(0.5);
}

export function teardown(data) {
  console.log('✅ Spike Test Complete!');
  console.log('');
  console.log('🔍 Key Observations:');
  console.log('- Did system survive the initial spike?');
  console.log('- How long to recover?');
  console.log('- Were there cascading failures?');
  console.log('- Did errors persist after spike ended?');
  console.log('');
  console.log('💡 Improvements to Consider:');
  console.log('- Auto-scaling configuration');
  console.log('- Connection pool sizing');
  console.log('- Rate limiting');
  console.log('- Circuit breakers');
}
