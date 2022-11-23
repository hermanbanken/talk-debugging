from locust import HttpUser, between, task
import random

# usage: $ locust -f loadtest.py

equations = ['6 + 6 * 6', '1 / 0 + 100 / 10 * 42', '1 + 1 + 40 + 1 + 1', '1 * 1 * 1 * 1 * 1 * 1']

class CalculatorUser(HttpUser):
    wait_time = between(5, 15)
    
    def on_start(self):
        self.client.get("/")
    
    @task
    def index(self):
        self.client.post("/", {
            "equation": random.choice(equations),
        })
