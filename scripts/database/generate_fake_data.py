import random

categories = [
    "Miscellaneous", "Reversing", "Pwn", "Linux", "Web", "Crypto", "Forensics", "OSINT", "Steganography"
]
challenge_types = ["static"] * 30 + ["dynamic"] * 18 + ["on-demand"] * 12
deployment_options = {
    "static": "http",
    "dynamic": ["http", "nc"],
    "on-demand": "http"
}
authors = ["TheAlpha", "Naughtyb0y", "e4stw1nd", "CyberKid", "Hackerman"]
tags_pool = ["pwn", "heap", "python", "reversing", "crypto", "web", "linux", "misc", "df", "forensics"]

category_sql = "-- Create categories\nINSERT INTO categories (category_id, category_name) VALUES\n"
category_sql += ",\n".join([f"    ({i+1}, '{categories[i]}')" for i in range(len(categories))]) + ";\n"

challenge_sql = "-- Add challenge data\nINSERT INTO challenges (\n    chall_id, chall_name, category_id, \n    prompt, flag, type, \n    points, files, requirements, \n    author, tags, links,\n    deployment, port, subd\n) VALUES \n"

hints_sql = "-- Add hints\nINSERT INTO hints (hid, chall_id, hint, cost) VALUES\n"

hint_counter = 1
challenge_entries = []
hint_entries = []

for i in range(1, 61):  
    chall_name = f"challenge_{i}"
    category_id = random.randint(1, 9)
    chall_type = random.choice(challenge_types)
    
    deployment = (
        random.choice(deployment_options[chall_type])
        if isinstance(deployment_options[chall_type], list)
        else deployment_options[chall_type]
    )
    port = 443 if chall_type == "on-demand" else (80 if deployment == "http" else random.randint(30000, 31000))
    subd = f"subd_{i}" if chall_type == "on-demand" else ""

    flag = f"isolet-dev{{random_flag_{i}}}"
    points = random.choice([100, 200, 300, 400, 500])
    files = "{}" if chall_type == "on-demand" else f"{{file_{i}.bin}}"
    requirements = "{}"
    author = random.choice(authors)
    tags = "{" + ",".join(random.sample(tags_pool, random.randint(1, 3))) + "}"
    links = "{}"

    challenge_entries.append(
        f"({i}, '{chall_name}', {category_id}, 'Solve {chall_name}!', '{flag}', '{chall_type}', "
        f"{points}, '{files}', '{requirements}', '{author}', '{tags}', '{links}', '{deployment}', {port}, '{subd}')"
    )

    num_hints = random.choices([0, 1, 2, 3], weights=[0.3, 0.5, 0.15, 0.05])[0]
    for _ in range(num_hints):
        hint_entries.append(
            f"({hint_counter}, {i}, 'Hint for {chall_name}', {random.choice([10, 20, 30])})"
        )
        hint_counter += 1

challenge_sql += ",\n".join(challenge_entries) + ";\n"
hints_sql += ",\n".join(hint_entries) + ";\n" if hint_entries else "-- No hints generated\n"

update_sql = "UPDATE challenges SET visible = true;\nUPDATE hints SET visible = true;"

final_sql = category_sql + "\n" + challenge_sql + "\n" + hints_sql + "\n" + update_sql + "\n\n" + """SELECT setval('categories_category_id_seq', 4);
SELECT setval('challenges_chall_id_seq', 60);
SELECT setval('hints_hid_seq', 69);""" + "\n\n"

final_sql += """INSERT INTO config (key, value) VALUES
('EVENT_START', '1728813375'),
('EVENT_END', '1799986175'),
('POST_EVENT', 'false'),
('TEAM_LEN', '4'),
('CTF_NAME', 'isolet-dev'),
('CONCURRENT_INSTANCES', '1'),
('INSTANCE_HOSTNAME', 'ctf.isolet-dev.in'),
('INSTANCE_TIME', '30'),
('MAX_INSTANCE_TIME', '60'),
('EMAIL_ID', 'menoone@isolet-dev.in'),
('EMAIL_AUTH', 'jdrglctcdtdlbi'),
('PUBLIC_URL', 'isolet-dev.in'),
('SMTP_HOST', 'smtp.gmail.com'),
('SMTP_PORT', '587'),
('SYNC_CONFIG_SECONDS', '30'),
('USER_REGISTRATION', 'true'),
('API_RATE_LIMIT', 'false'),
('EMAIL_VERIFICATION', 'false');"""

with open("fake_data.sql", "w") as f:
    f.write(final_sql)

print("Fake data SQL script generated: fake_data.sql")
