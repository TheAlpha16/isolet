from faker import Faker
import random
import json

fake = Faker()

chall_types = ['static', 'dynamic', 'on-demand']

categories = ["web","pwn","reversing","crypto","misc","forensics","stego","osint","hardware","networking","mobile","blockchain","cloud","ai","iot"]

def generate_fake_challenge():
    category = random.choice(categories)
    chall_type = random.choice(chall_types)
    files = [fake.file_name(extension='txt')
             for _ in range(random.randint(0, 2))]
    hints = [{"hint": fake.sentence(), "cost": random.randint(1, 500), "visible": random.choice([True, False])} for _ in range(random.randint(0, 2))]
    tags = [fake.word() for _ in range(random.randint(0, 3))]

    final = {
        "chall_name": fake.catch_phrase(),
        "category": category,
        "prompt": fake.paragraph(),
        "flag": fake.password(length=12, special_chars=True),
        "type": chall_type,
        "points": random.randint(50, 500),
        "author": fake.user_name(),
    }

    if len(hints) > 0:
        final["hints"] = hints
    if len(files) > 0:
        final["files"] = files
    if len(tags) > 0:
        final["tags"] = tags
    if chall_type != 'static':
        final["port"] = random.randint(1024, 65535)
        final["subd"] = fake.hostname()
        final["cpu"] = random.randint(1, 10)
        final["mem"] = random.randint(1, 16)
    
    return final

num_challenges = 10
fake_challenges = [generate_fake_challenge() for _ in range(num_challenges)]

fake_challenges_json = json.dumps(fake_challenges, indent=4)
open('../../challenges/challs.json', 'w').write(fake_challenges_json)
