import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';


const slowResponseRate = new Rate('slow_responses');

let authToken = '';

export const options = {
  scenarios: {
    // ramping_load: {
    //   executor: 'ramping-arrival-rate',
    //   startRate: 0,
    //   timeUnit: '1s',
    //   stages: [
    //     { duration: '1m', target: 1000 }, // увеличение до 1000 RPS за 2 минуты
    //     // { duration: '8m', target: 1000 }, // поддержание 1000 RPS
    //   ],
    //   preAllocatedVUs: 200,
    //   maxVUs: 5000,
    // },
    // constant_load: {
    //   executor: 'constant-arrival-rate',
    //   rate: 1000,                // Number of iterations per timeUnit
    //   timeUnit: '1s',         // Generate 50 iterations per minute
    //   duration: '5m',        // Test duration
    //   preAllocatedVUs: 10,    // Initial pool of VUs
    //   maxVUs: 2000,           // Maximum number of VUs to handle the rate
    // },
    // shared_iterations_example: {
    //   executor: 'per-vu-iterations',
    //   vus: 10000,
    //   iterations: 1000
    // },
    ramping_rate_example: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 10,
      maxVUs: 3000,
      stages: [
        { duration: '1m', target: 500 },
        { duration: '4m', target: 500 },
        // { duration: '1m', target: 0 },
      ],
    }
  },
};

const predefinedUsernames = [
  "user1_test",
  "user2_test",
  "user3_test",
  "user4_test",
  "user5_test",
  "user6_test",
  "user7_test",
  "user8_test",
  "user9_test",
  "user10_test"
];

// Функция для получения случайного username из массива
function getRandomUsername() {
  const index = Math.floor(Math.random() * predefinedUsernames.length);
  return predefinedUsernames[index];
}

function generateRandomString(length = 10) {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < length; i++) {
      result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
}

export function getAuthToken() {
  const payload = {
    // "username": generateRandomString(),
    "username": getRandomUsername(),
    "password": "sdfgsgdfgsg",
  };

  const params = {
    headers: {
      'accept': 'application/json',
      'Content-Type': 'application/json',
    },
  };

  const response = http.post('http://172.20.1.1:8081/api/auth', JSON.stringify(payload), params);


  check(response, {
    'Status is 200': (r) => r.status === 200,
    'Response contains token': (r) => r.body.includes('token'),
    'Response time is less than 50ms': (r) => r.timings.duration < 50
  });

  if (response.timings.duration >= 50) {
    slowResponseRate.add(1);
    console.log(`Slow response detected: ${response.timings.duration}ms`);
  }

  // Сохраняем токен
  const responseBody = JSON.parse(response.body);
  authToken = responseBody.token;

  return authToken;
}


export function getInfoRequest(authToken) {
  const params = {
    headers: {
      'accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authToken}`
    },
  };

  const response = http.get('http://172.20.1.1:8081/api/info', params);
  
  check(response, {
    'Authenticated request successful': (r) => r.status === 200
  });
}

export function sendCoinRequest(authToken) {
  const payload = {
      "toUser": getRandomUsername(), // используем ту же функцию для получения случайного username
      "amount": Math.floor(Math.random() * 10) + 1 // случайное число от 1 до 10
  };

  const params = {
      headers: {
          'accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${authToken}`
      },
  };

  const response = http.post('http://172.20.1.1:8081/api/sendCoin', JSON.stringify(payload), params);

  check(response, {
      'Send coin request successful': (r) => r.status === 200,
      'Response time is less than 50ms': (r) => r.timings.duration < 50
  });

  if (response.timings.duration >= 50) {
      slowResponseRate.add(1);
      console.log(`Slow response detected: ${response.timings.duration}ms`);
  }
}

export function setup() {
  console.log('Setup: Creating test users...');
  
  const createdUsers = [];
  
  // Создаем тестовых пользователей
  for (const username of predefinedUsernames) {
      const payload = {
          "username": username,
          "password": "sdfgsgdfgsg",
      };

      const params = {
          headers: {
              'accept': 'application/json',
              'Content-Type': 'application/json',
          },
      };

      // Регистрируем пользователя
      const response = http.post('http://172.20.1.1:8081/api/register', JSON.stringify(payload), params);
      
      // Проверяем результат
      const success = check(response, {
          'User registration successful': (r) => r.status === 200// 409 если пользователь уже существует
      });

      if (success) {
          createdUsers.push(username);
      }

      // Небольшая пауза между регистрациями
      sleep(0.1);
  }

  console.log(`Setup completed. Created/Verified ${createdUsers.length} users`);
  
  // Возвращаем результаты setup
  return {
      setupCompleted: true,
      usersCreated: createdUsers.length,
      timestamp: new Date().toISOString()
  };
}

// Обновляем основную функцию
export default function (data) {
  // Проверяем, что setup прошел успешно
  if (!data.setupCompleted) {
      console.error('Setup was not completed successfully');
      return;
  }

  // Используем authToken вместо token
  authToken = getAuthToken();
  // getInfoRequest(authToken);
  sendCoinRequest(authToken);
  
  sleep(Math.random() * 4 + 1);
}