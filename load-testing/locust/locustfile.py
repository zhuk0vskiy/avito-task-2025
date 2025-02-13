from locust import HttpUser, task, between
import json

class MyUser(HttpUser):
    wait_time = between(1, 2)

    @task
    def authenticate(self):
        payload = {
            "username": "test",
            "password": "pass"
        }
        headers = {'Content-Type': 'application/json'}
        response = self.client.post("/api/auth", json=payload, headers=headers)
        if response.status_code == 200:
            token = response.json()["token"]
            self.token = token
            self.client.cookies.set("jwtToken1", token)

    # @task
    # def get_info(self):
    #     if hasattr(self, 'token'):  # Проверяем, что токен был получен
    #         headers = {
    #             'Authorization': f'Bearer {self.token}',
    #             'Content-Type': 'application/json'
    #         }
    #         self.client.get("/api/info", headers=headers)

    # @task
    # def send_coins(self):
    #     payload = {
    #         "toUser": "testuser2",
    #         "amount": "1"
    #     }
    #     headers = {
    #         'Authorization': f'Bearer {self.client.cookies.get("jwtToken1")}',
    #         'Content-Type': 'application/json'
    #     }
    #     self.client.post("/api/send", json=payload, headers=headers)