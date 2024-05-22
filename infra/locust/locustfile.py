from locust import HttpUser, TaskSet, task, between
import random
import uuid

class UserBehavior(TaskSet):

    @task
    def create_wallet_and_order(self):
        # Step 1: Create a new wallet
        response = self.client.post("/wallets")
        if response.status_code == 200:
            user_id = response.json()["user_id"]

            # Step 2: Deposit an amount in the wallet
            deposit_amount = random.randint(50, 150)  # random amount to deposit
            deposit_response = self.client.patch("/wallets/deposit", json={
                "userId": user_id,
                "amount": deposit_amount
            })

            if deposit_response.status_code == 200:
                # Step 3: Create an order
                self.client.post("/orders", json={
                    "userId": user_id,
                    "type": 1,
                    "quantity": 1,
                    "price": random.randint(1, 10)  # random price for the order
                })

class WebsiteUser(HttpUser):
    tasks = [UserBehavior]
    wait_time = between(1, 2)
