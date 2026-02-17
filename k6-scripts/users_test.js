import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 100 },
    { duration: '30s', target: 300 },
    { duration: '30s', target: 500 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'], 
    http_req_failed: ['rate==0'], 
  },
};

const URL = 'http://127.0.0.1:3000/api/v1/users';
const JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzAxMDU4NjMsImlkIjoiYWM1OTU1NTktM2NhNi00YTQyLWFlM2UtYmQ5MzU4MDkzNjQzIiwibmFtZSI6IkFkbWluaXN0cmF0b3IiLCJ1c2VybmFtZSI6ImFkbWluaXN0cmF0b3IifQ.fTdFIwafD9uwAlY27-FbqsfinYeoIypOrN54-RJcIFE";

export default function () {
  const params = {
    headers: {
      Authorization: `Bearer ${JWT_TOKEN}`,
      'Content-Type': 'application/json',
    },
  };

  const res = http.get(URL, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0);
}
