from locust import HttpUser, task

class StatusUser(HttpUser):
    @task
    def status(self):
        self.client.get("/status")
