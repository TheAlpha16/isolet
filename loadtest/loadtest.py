from locust import HttpUser, task, between
import random
import hashlib
import logging
import os
import time

# Configure logging
logging.basicConfig(
    filename='ctf_load_test.log', level=logging.ERROR,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

class CTFUser(HttpUser):
    wait_time = between(1, 3)  # Simulate real user wait times

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.token = None
        self.email = None
        self.password = "Strongpasswd123."
        self.username = None
        self.team_id = None
        self.unique_identifier = hashlib.md5(os.urandom(16)).hexdigest()

    def on_start(self):
        """Simulate user setup: register, login, create team."""
        self.email = f"{self.unique_identifier}@test.com"
        self.username = self.unique_identifier

        if not self.register():
            logging.error(f"User {self.username} failed to register")
            return

        if not self.login():
            logging.error(f"User {self.username} failed to login")
            return

        # time.sleep(random.uniform(0.5, 1.5))  # Simulate delay after login

        if not self.create_team():
            logging.error(f"User {self.username} failed to create a team")

    def _request(self, method, endpoint, data=None, expected_status={200}, retries=30):
        """Helper method for API requests with retries."""
        headers = {"Cookie": f"token={self.token};"} if self.token else {}
        for _ in range(retries):
            with self.client.request(method, endpoint, headers=headers, data=data, catch_response=True) as response:
                if response.status_code in expected_status:
                    response.success()
                    return response.json() if response.text else {}
                response.failure(f"{endpoint} failed: {response.text}")
                logging.error(f"{method} {endpoint} failed: {response.text}")
                time.sleep(0.8)  # Backoff before retry
        return None

    def register(self):
        """Register a new user."""

        while True:
            if bool(self._request("POST", "/auth/register", {
                "email": self.email,
                "password": self.password,
                "username": self.username,
                "confirm": self.password
            })) is True:
                logging.info(f"User {self.username} registered successfully")
                return True
            time.sleep(0.8)

    def login(self):
        """Login and obtain token."""
        while True:
            response = self._request("POST", "/auth/login", {
                "email": self.email,
                "password": self.password
            })
            if response and 'token' in self.client.cookies:
                self.token = self.client.cookies['token']
                logging.info(f"User {self.username} logged in successfully")
                return True
            time.sleep(0.8)

    def create_team(self):
        """Create a team."""
        while True:
            response = self._request("POST", "/onboard/team/create", {
                "teamname": self.unique_identifier,
                "password": "notasecurepasswd"
            }, expected_status={200, 403})
            if response:
                self.team_id = response.get('team_id')
                self.login()
                logging.info(f"User {self.username} created team {self.unique_identifier} successfully")
                return True
            time.sleep(0.8)

    @task(20)
    def get_metadata(self):
        """Fetch user metadata."""
        self._request("GET", "/auth/metadata")

    @task(20)
    def get_identify(self):
        """Identify user."""
        self._request("GET", "/api/identify")

    @task(5)
    def self_team(self):
        """Fetch team details."""
        self._request("GET", "/api/profile/team/self")

    @task(30)
    def get_challenges(self):
        """Fetch available challenges."""
        self._request("GET", "/api/challs")

    @task(20)
    def submit_flag(self):
        """Submit a flag with random correctness."""
        chall_id = str(random.randint(1, 60))
        flag = f"isolet-dev{{random_flag_{chall_id}}}" if random.random(
        ) < 0.05 else "isolet-dev{fake}"
        self._request("POST", "/api/submit",
                      {"chall_id": chall_id, "flag": flag}, expected_status={200, 403, 400})

    @task(2)
    def check_status(self):
        """Check challenge status."""
        self._request("GET", "/api/status")

    @task(10)
    def view_scoreboard(self):
        """View scoreboard."""
        self._request("GET", "/api/scoreboard?page=1")

    @task(10)
    def view_topscore(self):
        """View top scoreboard."""
        self._request("GET", "/api/scoreboard/top")

    def on_stop(self):
        """Optional cleanup on user stop."""
        pass
