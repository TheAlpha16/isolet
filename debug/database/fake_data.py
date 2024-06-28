from faker import Faker
import random
import json

fake = Faker()

chall_types = ['static', 'dynamic', 'on-demand']

categories = [
    {"category_id": 1, "category_name": "web"},
    {"category_id": 2, "category_name": "pwn"},
    {"category_id": 3, "category_name": "reversing"},
    {"category_id": 4, "category_name": "crypto"},
    {"category_id": 5, "category_name": "misc"},
]

def generate_fake_challenge():
    category = random.choice(categories)
    chall_type = random.choice(chall_types)
    files = [fake.file_name(extension='txt')
             for _ in range(random.randint(1, 3))]
    hints = [fake.sentence() for _ in range(random.randint(1, 3))]
    tags = [fake.word() for _ in range(random.randint(1, 5))]

    return {
        "chall_name": fake.catch_phrase(),
        "category_id": category["category_id"],
        "prompt": fake.paragraph(),
        "flag": fake.password(length=12, special_chars=True),
        "type": chall_type,
        "points": random.randint(50, 500),
        "files": files,
        "hints": hints,
        "author": fake.user_name(),
        "tags": tags,
        "port": random.randint(1024, 65535) if chall_type != 'static' else None,
        "subd": fake.hostname(),
        "cpu": random.randint(1, 10),
        "mem": random.randint(1, 16)
    }


num_challenges = 10
fake_challenges = [generate_fake_challenge() for _ in range(num_challenges)]

fake_challenges_json = json.dumps(fake_challenges, indent=4)
open('../../challenges/challs.json', 'w').write(fake_challenges_json)
