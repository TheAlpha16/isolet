# ISOLET

## Description
Isolet is a framework to deploy linux wargames like [Bandit](https://overthewire.org/wargames/bandit/). It uses pre-configured templates to provide isolated instance using [kubernetes](https://kubernetes.io) pods for each user. 

## Configuration


## API routes
| route | methods | parameters | response | sample |
|:---:|:---:|:---:|:---:|:---:|
| /api/challs | GET | NONE | challenges | [{"chall_id":0, "level":1, "name":"demo", "prompt":"solve it", "tags":["ssh", "cat"]}] |
| /api/launch | POST | chall_id, userid, level | status | {"status": "success", "message": "3b369c0b1fd5419b2f81da89cf5480d2 32747"} |
| /api/stop | POST | userid, level | status | {"status": "failure", "message": "User does not exist"} |
| /api/submit | POST | userid, level, flag | status | {"status": "failure", "message": "Flag copy detected. Incident reported!"} |