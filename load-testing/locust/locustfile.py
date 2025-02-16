from locust import HttpUser, task, between
import random
import json

class CoinUser(HttpUser):
    # Ждем между запросами от 0.1 до 0.3 секунды
    wait_time = between(0.1, 0.3)
    
    def on_start(self):
        """Выполняется один раз при старте для каждого пользователя"""
        # Авторизация
        payload = {
            "username": self.get_random_username(),
            "password": "sdfgsgdfgsg"
        }
        
        response = self.client.post(
            "/api/auth",
            json=payload,
            headers={'Content-Type': 'application/json'}
        )
        
        if response.status_code == 200:
            self.token = response.json()["token"]
        else:
            self.token = None

    def get_random_username(self):
        """Получение случайного имени пользователя из предопределенного списка"""
        usernames = [
            "user1_test",
            "user2_test",
            "user3_test",
            "user4_test",
            "user5_test"
        ]
        return random.choice(usernames)

    @task(1)
    def send_coin(self):
        """Отправка монет"""
        if not self.token:
            return

        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.token}'
        }

        payload = {
            "toUsername": self.get_random_username(),
            "amount": random.randint(1, 10)
        }

        self.client.post(
            "/api/sendCoin",
            json=payload,
            headers=headers
        )

# Запуск:
# locust -f locustfile.py --host=http://172.20.1.1:8081