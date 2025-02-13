import http from 'k6/http';
import { sleep, check } from 'k6';

export default function () {
  const payload = {
    "username": "test5456456",
    "password": "test"
  };

  const params = {
    headers: {
      'accept': 'application/json',
      'Content-Type': 'application/json',
    },
  };

  // Добавляем задержку перед запросом
  sleep(0.1);

  const response = http.post('http://localhost:8080/api/auth', JSON.stringify(payload), params);

  // Проверяем успешность запроса
  // check(response, {
  //   'Status is 200': (r) => r.status === 200,
  //   'Response contains token': (r) => r.body.includes('token')
  // });

  // Добавляем задержку после запроса
  sleep(0.1);
}